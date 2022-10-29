package extension

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/template"
)

type commentExtension struct {
	PostCommentService   service.PostCommentService
	PostCommentAssembler assembler.PostCommentAssembler
	Template             *template.Template
}

func RegisterCommentFunc(template *template.Template, postCommentService service.PostCommentService, postCommentAssembler assembler.PostCommentAssembler) {
	ce := commentExtension{
		PostCommentService:   postCommentService,
		Template:             template,
		PostCommentAssembler: postCommentAssembler,
	}
	ce.addGetLatestComment()
	ce.addGetCommentCount()
}

func (ce *commentExtension) addGetLatestComment() {
	getLatestComment := func(top int) ([]*vo.PostCommentWithPost, error) {
		commentQuery := param.CommentQuery{
			Sort:          &param.Sort{Fields: []string{"createTime,desc"}},
			Page:          param.Page{PageNum: 0, PageSize: top},
			CommentStatus: consts.CommentStatusPublished.Ptr(),
		}
		comments, _, err := ce.PostCommentService.Page(context.Background(), commentQuery, consts.CommentTypePost)
		if err != nil {
			return nil, err
		}
		return ce.PostCommentAssembler.ConvertToWithPost(context.Background(), comments)
	}
	ce.Template.AddFunc("getLatestComment", getLatestComment)
}

func (ce *commentExtension) addGetCommentCount() {
	getCommentCount := func() (int64, error) {
		count, err := ce.PostCommentService.CountByStatus(context.Background(), consts.CommentStatusPublished)
		return count, err
	}
	ce.Template.AddFunc("getCommentCount", getCommentCount)
}
