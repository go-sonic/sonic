package api

import (
	"context"
	"html/template"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type PostHandler struct {
	OptionService        service.OptionService
	PostService          service.PostService
	PostCommentService   service.PostCommentService
	PostCommentAssembler assembler.PostCommentAssembler
}

func NewPostHandler(
	optionService service.OptionService,
	postService service.PostService,
	postCommentService service.PostCommentService,
	postCommentAssembler assembler.PostCommentAssembler,
) *PostHandler {
	return &PostHandler{
		OptionService:        optionService,
		PostService:          postService,
		PostCommentService:   postCommentService,
		PostCommentAssembler: postCommentAssembler,
	}
}

func (p *PostHandler) ListTopComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "postID")
	if err != nil {
		return nil, err
	}
	pageSize := p.OptionService.GetOrByDefault(_ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{}
	err = ctx.BindAndValidate(&commentQuery)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if commentQuery.Sort != nil && len(commentQuery.Fields) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &postID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	comments, totalCount, err := p.PostCommentService.Page(_ctx, commentQuery, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}
	_ = p.PostCommentAssembler.ClearSensitiveField(_ctx, comments)
	commenVOs, err := p.PostCommentAssembler.ConvertToWithHasChildren(_ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commenVOs, totalCount, commentQuery.Page), nil
}

func (p *PostHandler) ListChildren(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "postID")
	if err != nil {
		return nil, err
	}
	parentID, err := util.ParamInt32(_ctx, ctx, "parentID")
	if err != nil {
		return nil, err
	}
	children, err := p.PostCommentService.GetChildren(_ctx, parentID, postID, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}
	_ = p.PostCommentAssembler.ClearSensitiveField(_ctx, children)
	return p.PostCommentAssembler.ConvertToDTOList(_ctx, children)
}

func (p *PostHandler) ListCommentTree(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "postID")
	if err != nil {
		return nil, err
	}
	pageSize := p.OptionService.GetOrByDefault(_ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{}
	err = ctx.BindAndValidate(&commentQuery)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if commentQuery.Sort != nil && len(commentQuery.Fields) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &postID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	allComments, err := p.PostCommentService.GetByContentID(_ctx, postID, consts.CommentTypePost, commentQuery.Sort)
	if err != nil {
		return nil, err
	}
	_ = p.PostCommentAssembler.ClearSensitiveField(_ctx, allComments)
	commentVOs, total, err := p.PostCommentAssembler.PageConvertToVOs(_ctx, allComments, commentQuery.Page)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentVOs, total, commentQuery.Page), nil
}

func (p *PostHandler) ListComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "postID")
	if err != nil {
		return nil, err
	}
	pageSize := p.OptionService.GetOrByDefault(_ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{}
	err = ctx.BindAndValidate(&commentQuery)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if commentQuery.Sort != nil && len(commentQuery.Fields) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &postID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	comments, total, err := p.PostCommentService.Page(_ctx, commentQuery, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}
	_ = p.PostCommentAssembler.ClearSensitiveField(_ctx, comments)
	result, err := p.PostCommentAssembler.ConvertToWithParentVO(_ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(result, total, commentQuery.Page), nil
}

func (p *PostHandler) CreateComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	comment := param.Comment{}
	err := ctx.BindAndValidate(&comment)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if comment.AuthorURL != "" {
		err = util.Validate.Var(comment.AuthorURL, "http_url")
		if err != nil {
			return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
		}
	}
	comment.Author = template.HTMLEscapeString(comment.Author)
	comment.AuthorURL = template.HTMLEscapeString(comment.AuthorURL)
	comment.Content = template.HTMLEscapeString(comment.Content)
	comment.Email = template.HTMLEscapeString(comment.Email)
	comment.CommentType = consts.CommentTypePost
	result, err := p.PostCommentService.CreateBy(_ctx, &comment)
	if err != nil {
		return nil, err
	}
	return p.PostCommentAssembler.ConvertToDTO(_ctx, result)
}

func (p *PostHandler) Like(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "postID")
	if err != nil {
		return nil, err
	}
	return nil, p.PostService.IncreaseLike(_ctx, postID)
}
