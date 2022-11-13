package listener

import (
	"bytes"
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
)

type CommentListener struct {
	OptionService      service.OptionService
	PostService        service.PostService
	PostAssembler      assembler.PostAssembler
	JournalService     service.JournalService
	SheetService       service.SheetService
	ThemeService       service.ThemeService
	EmailService       service.EmailService
	UserService        service.UserService
	BaseCommentService service.BaseCommentService
	Template           *template.Template
}

func NewCommentListener(
	optionService service.OptionService,
	postService service.PostService,
	journalService service.JournalService,
	sheetService service.SheetService,
	bus event.Bus,
	postAssembler assembler.PostAssembler,
	themeService service.ThemeService,
	emailService service.EmailService,
	userService service.UserService,
	template *template.Template,
	baseCommentService service.BaseCommentService,
) {
	c := &CommentListener{
		OptionService:      optionService,
		PostService:        postService,
		PostAssembler:      postAssembler,
		JournalService:     journalService,
		SheetService:       sheetService,
		ThemeService:       themeService,
		EmailService:       emailService,
		UserService:        userService,
		Template:           template,
		BaseCommentService: baseCommentService,
	}
	bus.Subscribe(event.CommentNewEventName, c.HandleCommentNew)
	bus.Subscribe(event.CommentReplyEventName, c.HandleCommentReply)
}

func (c *CommentListener) HandleCommentNew(ctx context.Context, ce event.Event) error {
	newCommentNotice, err := c.OptionService.GetOrByDefaultWithErr(ctx, property.CommentNewNotice, property.CommentNewNotice.DefaultValue)
	if err != nil {
		return err
	}
	if !newCommentNotice.(bool) {
		return nil
	}
	commentEvent, ok := ce.(*event.CommentNewEvent)
	if !ok {
		return nil
	}
	enabledAbsolutePath, err := c.OptionService.GetOrByDefaultWithErr(ctx, property.GlobalAbsolutePathEnabled, false)
	if err != nil {
		return err
	}
	blogBaseURL, err := c.OptionService.GetBlogBaseURL(ctx)
	if err != nil {
		return err
	}
	comment := commentEvent.Comment
	data := make(map[string]interface{})
	var subject string

	if comment.Type == consts.CommentTypePost || comment.Type == consts.CommentTypeSheet {
		post, err := c.PostService.GetByPostID(ctx, commentEvent.Comment.PostID)
		if err != nil {
			return nil
		}
		postDTO, err := c.PostAssembler.ConvertToMinimalDTO(ctx, post)
		if err != nil {
			return nil
		}
		data["pageFullPath"] = util.IfElse(enabledAbsolutePath.(bool), postDTO.FullPath, blogBaseURL+postDTO.FullPath).(string)
		data["pageTitle"] = postDTO.Title
		data["author"] = comment.Author
		data["content"] = comment.Content
		data["email"] = comment.Email
		data["status"] = comment.Status
		data["createTime"] = comment.CreateTime
		data["authorUrl"] = comment.AuthorURL
		if comment.Type == consts.CommentTypePost {
			subject = "Your blog post 《" + postDTO.Title + "》 has a new comment"
		} else {
			subject = "Your blog page 《" + postDTO.Title + "》 has a new comment"
		}
	} else if comment.Type == consts.CommentTypeJournal {
		journalPrefix, err := c.OptionService.GetJournalPrefix(ctx)
		if err != nil {
			return err
		}
		journals, err := c.JournalService.GetByJournalIDs(ctx, []int32{comment.PostID})
		if err != nil {
			return err
		}
		if len(journals) == 0 {
			return nil
		}
		journal := journals[comment.PostID]

		data["pageFullPath"] = blogBaseURL + "/" + journalPrefix
		data["pageTitle"] = journal.CreateTime.Format("2006-01-02 03:04")
		data["author"] = comment.Author
		data["content"] = comment.Content
		data["email"] = comment.Email
		data["status"] = comment.Status
		data["createTime"] = comment.CreateTime
		data["authorUrl"] = comment.AuthorURL
		subject = "Your blog journal has a new comment"
	}
	template := "common/mail_template/mail_notice"
	if exist, err := c.ThemeService.TemplateExist(ctx, "mail_template/mail_notice.tmpl"); err != nil && exist {
		t, err := c.ThemeService.Render(ctx, "mail_template/mail_notice")
		if err == nil {
			template = t
		}
	}
	content := bytes.Buffer{}
	err = c.Template.ExecuteTemplate(&content, template, data)
	if err != nil {
		return err
	}
	users, err := c.UserService.GetAllUser(ctx)
	if err != nil {
		return err
	}
	return c.EmailService.SendTemplateEmail(ctx, users[0].Email, subject, content.String())
}

func (c *CommentListener) HandleCommentReply(ctx context.Context, ce event.Event) error {
	commentReplyNotice, err := c.OptionService.GetOrByDefaultWithErr(ctx, property.CommentReplyNotice, property.CommentNewNotice.DefaultValue)
	if err != nil {
		return err
	}
	if !commentReplyNotice.(bool) {
		return nil
	}
	commentEvent, ok := ce.(*event.CommentReplyEvent)
	if !ok {
		return nil
	}

	comment := commentEvent.Comment
	if comment.Status != consts.CommentStatusPublished {
		return nil
	}
	parentComment, err := c.BaseCommentService.GetByID(ctx, comment.ParentID)
	if err != nil {
		return err
	}
	if !parentComment.AllowNotification || parentComment.Status != consts.CommentStatusPublished {
		return nil
	}

	blogTitle, err := c.OptionService.GetOrByDefaultWithErr(ctx, property.BlogTitle, "")
	if err != nil {
		return err
	}
	enabledAbsolutePath, err := c.OptionService.GetOrByDefaultWithErr(ctx, property.GlobalAbsolutePathEnabled, false)
	if err != nil {
		return err
	}
	blogBaseURL, err := c.OptionService.GetBlogBaseURL(ctx)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	var subject string
	if comment.Type == consts.CommentTypePost || comment.Type == consts.CommentTypeSheet {
		post, err := c.PostService.GetByPostID(ctx, commentEvent.Comment.PostID)
		if err != nil {
			return nil
		}
		postDTO, err := c.PostAssembler.ConvertToMinimalDTO(ctx, post)
		if err != nil {
			return nil
		}
		data["pageFullPath"] = util.IfElse(enabledAbsolutePath.(bool), postDTO.FullPath, blogBaseURL+postDTO.FullPath).(string)
		data["pageTitle"] = postDTO.Title
		data["baseAuthor"] = parentComment.Author
		data["baseContent"] = parentComment.Content
		data["replyAuthor"] = comment.Author
		data["replyContent"] = comment.Content
		data["baseAuthorEmail"] = parentComment.Email
		data["replyAuthorEmail"] = comment.Email
		data["status"] = comment.Status
		data["createTime"] = comment.CreateTime
		data["authorUrl"] = comment.AuthorURL
		if comment.Type == consts.CommentTypePost {
			subject = "You have a new reply in the 《" + post.Title + "》 article you comment on " + blogTitle.(string)
		} else {
			subject = "You have a new reply in the 《" + post.Title + "》 page you comment on " + blogTitle.(string)
		}
	} else if comment.Type == consts.CommentTypeJournal {
		blogBaseURL, err := c.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return err
		}
		journalPrefix, err := c.OptionService.GetJournalPrefix(ctx)
		if err != nil {
			return err
		}
		journals, err := c.JournalService.GetByJournalIDs(ctx, []int32{comment.PostID})
		if err != nil {
			return err
		}
		if len(journals) == 0 {
			return nil
		}
		journal := journals[comment.PostID]
		data["pageFullPath"] = blogBaseURL + "/" + journalPrefix
		data["pageTitle"] = journal.CreateTime.Format("2006-01-02 03:04")
		data["baseAuthor"] = parentComment.Author
		data["baseContent"] = parentComment.Content
		data["replyAuthor"] = comment.Author
		data["replyContent"] = comment.Content
		data["baseAuthorEmail"] = parentComment.Email
		data["replyAuthorEmail"] = comment.Email
		data["status"] = comment.Status
		data["createTime"] = comment.CreateTime
		data["authorUrl"] = comment.AuthorURL
		subject = "You have a new reply in the journal page you comment on " + blogTitle.(string)
	}
	template := "common/mail_template/mail_reply"
	if exist, err := c.ThemeService.TemplateExist(ctx, "mail_template/mail_reply.tmpl"); err != nil && exist {
		t, err := c.ThemeService.Render(ctx, "mail_template/mail_reply")
		if err == nil {
			template = t
		}
	}
	content := bytes.Buffer{}
	err = c.Template.ExecuteTemplate(&content, template, data)
	if err != nil {
		return err
	}
	users, err := c.UserService.GetAllUser(ctx)
	if err != nil {
		return err
	}
	return c.EmailService.SendTemplateEmail(ctx, users[0].Email, subject, content.String())
}
