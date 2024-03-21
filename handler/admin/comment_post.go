package admin

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/service/impl"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type PostCommentHandler struct {
	PostCommentService   service.PostCommentService
	OptionService        service.OptionService
	PostService          service.PostService
	PostAssembler        assembler.PostAssembler
	PostCommentAssembler assembler.PostCommentAssembler
}

func NewPostCommentHandler(
	postCommentHandler service.PostCommentService,
	optionService service.OptionService,
	postService service.PostService,
	postAssembler assembler.PostAssembler,
	postCommentAssembler assembler.PostCommentAssembler,
) *PostCommentHandler {
	return &PostCommentHandler{
		PostCommentService:   postCommentHandler,
		OptionService:        optionService,
		PostService:          postService,
		PostAssembler:        postAssembler,
		PostCommentAssembler: postCommentAssembler,
	}
}

func (p *PostCommentHandler) ListPostComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var commentQuery param.CommentQueryNoEnum
	err := ctx.BindAndValidate(&commentQuery)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	commentQuery.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	comments, totalCount, err := p.PostCommentService.Page(_ctx, param.AssertCommentQuery(commentQuery), consts.CommentTypePost)
	if err != nil {
		return nil, err
	}
	commentDTOs, err := p.PostCommentAssembler.ConvertToWithPost(_ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentDTOs, totalCount, commentQuery.Page), nil
}

func (p *PostCommentHandler) ListPostCommentLatest(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	top, err := util.MustGetQueryInt32(_ctx, ctx, "top")
	if err != nil {
		return nil, err
	}
	commentQuery := param.CommentQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: int(top)},
	}
	comments, _, err := p.PostCommentService.Page(_ctx, commentQuery, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}
	return p.PostCommentAssembler.ConvertToWithPost(_ctx, comments)
}

func (p *PostCommentHandler) ListPostCommentAsTree(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "postID")
	if err != nil {
		return nil, err
	}
	pageNum, err := util.MustGetQueryInt32(_ctx, ctx, "page")
	if err != nil {
		return nil, err
	}
	pageSize, err := p.OptionService.GetOrByDefaultWithErr(_ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return nil, err
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: int(pageNum)}
	allComments, err := p.PostCommentService.GetByContentID(_ctx, postID, consts.CommentTypePost, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return nil, err
	}
	commentVOs, totalCount, err := p.PostCommentAssembler.PageConvertToVOs(_ctx, allComments, page)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentVOs, totalCount, page), nil
}

func (p *PostCommentHandler) ListPostCommentWithParent(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "postID")
	if err != nil {
		return nil, err
	}
	pageNum, err := util.MustGetQueryInt32(_ctx, ctx, "page")
	if err != nil {
		return nil, err
	}

	pageSize, err := p.OptionService.GetOrByDefaultWithErr(_ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return nil, err
	}

	page := param.Page{PageNum: int(pageNum), PageSize: pageSize.(int)}

	comments, totalCount, err := p.PostCommentService.Page(_ctx, param.CommentQuery{
		ContentID: &postID,
		Page:      page,
		Sort:      &param.Sort{Fields: []string{"createTime,desc"}},
	}, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}

	commentsWithParent, err := p.PostCommentAssembler.ConvertToWithParentVO(_ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentsWithParent, totalCount, page), nil
}

func (p *PostCommentHandler) CreatePostComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var commentParam *param.AdminComment
	err := ctx.BindAndValidate(&commentParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	user, err := impl.MustGetAuthorizedUser(_ctx)
	if err != nil || user == nil {
		return nil, err
	}
	blogURL, err := p.OptionService.GetBlogBaseURL(_ctx)
	if err != nil {
		return nil, err
	}
	commonParam := param.Comment{
		Author:            user.Username,
		Email:             user.Email,
		AuthorURL:         blogURL,
		Content:           commentParam.Content,
		PostID:            commentParam.PostID,
		ParentID:          commentParam.ParentID,
		AllowNotification: true,
		CommentType:       consts.CommentTypePost,
	}
	comment, err := p.PostCommentService.CreateBy(_ctx, &commonParam)
	if err != nil {
		return nil, err
	}
	return p.PostCommentAssembler.ConvertToDTO(_ctx, comment)
}

func (p *PostCommentHandler) UpdatePostComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	commentID, err := util.ParamInt32(_ctx, ctx, "commentID")
	if err != nil {
		return nil, err
	}
	var commentParam *param.Comment
	err = ctx.BindAndValidate(&commentParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	if commentParam.AuthorURL != "" {
		err = util.Validate.Var(commentParam.AuthorURL, "url")
		if err != nil {
			return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("url is not available")
		}
	}
	comment, err := p.PostCommentService.UpdateBy(_ctx, commentID, commentParam)
	if err != nil {
		return nil, err
	}

	return p.PostCommentAssembler.ConvertToDTO(_ctx, comment)
}

func (p *PostCommentHandler) UpdatePostCommentStatus(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	commentID, err := util.ParamInt32(_ctx, ctx, "commentID")
	if err != nil {
		return nil, err
	}
	strStatus, err := util.ParamString(_ctx, ctx, "status")
	if err != nil {
		return nil, err
	}
	status, err := consts.CommentStatusFromString(strStatus)
	if err != nil {
		return nil, err
	}
	return p.PostCommentService.UpdateStatus(_ctx, commentID, status)
}

func (p *PostCommentHandler) UpdatePostCommentStatusBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	strStatus, err := util.ParamString(_ctx, ctx, "status")
	if err != nil {
		return nil, err
	}
	status, err := consts.CommentStatusFromString(strStatus)
	if err != nil {
		return nil, err
	}

	ids := make([]int32, 0)
	err = ctx.BindAndValidate(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}
	comments, err := p.PostCommentService.UpdateStatusBatch(_ctx, ids, status)
	if err != nil {
		return nil, err
	}
	return p.PostCommentAssembler.ConvertToDTOList(_ctx, comments)
}

func (p *PostCommentHandler) DeletePostComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	commentID, err := util.ParamInt32(_ctx, ctx, "commentID")
	if err != nil {
		return nil, err
	}
	return nil, p.PostCommentService.Delete(_ctx, commentID)
}

func (p *PostCommentHandler) DeletePostCommentBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	ids := make([]int32, 0)
	err := ctx.BindAndValidate(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}
	return nil, p.PostCommentService.DeleteBatch(_ctx, ids)
}
