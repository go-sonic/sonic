package admin

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type PostHandler struct {
	PostService   service.PostService
	PostAssembler assembler.PostAssembler
}

func NewPostHandler(postService service.PostService, postAssembler assembler.PostAssembler) *PostHandler {
	return &PostHandler{
		PostService:   postService,
		PostAssembler: postAssembler,
	}
}

func (p *PostHandler) ListPosts(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postQuery := param.PostQueryNoEnum{}
	err := ctx.BindAndValidate(&postQuery)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if postQuery.Sort == nil {
		postQuery.Sort = &param.Sort{Fields: []string{"topPriority,desc", "createTime,desc"}}
	}
	posts, totalCount, err := p.PostService.Page(_ctx, param.AssertPostQuery(postQuery))
	if err != nil {
		return nil, err
	}
	if postQuery.More == nil || *postQuery.More {
		postVOs, err := p.PostAssembler.ConvertToListVO(_ctx, posts)
		return dto.NewPage(postVOs, totalCount, postQuery.Page), err
	}
	postDTOs := make([]*dto.Post, 0)
	for _, post := range posts {
		postDTO, err := p.PostAssembler.ConvertToSimpleDTO(_ctx, post)
		if err != nil {
			return nil, err
		}
		postDTOs = append(postDTOs, postDTO)
	}
	return dto.NewPage(postDTOs, totalCount, postQuery.Page), nil
}

func (p *PostHandler) ListLatestPosts(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	top, err := util.MustGetQueryInt32(_ctx, ctx, "top")
	if err != nil {
		top = 10
	}
	postQuery := param.PostQuery{
		Page: param.Page{
			PageSize: int(top),
			PageNum:  0,
		},
		Sort: &param.Sort{
			Fields: []string{"createTime,desc"},
		},
		Keyword:    nil,
		CategoryID: nil,
		More:       util.BoolPtr(false),
	}
	posts, _, err := p.PostService.Page(_ctx, postQuery)
	if err != nil {
		return nil, err
	}
	postMinimals := make([]*dto.PostMinimal, 0, len(posts))

	for _, post := range posts {
		postMinimal, err := p.PostAssembler.ConvertToMinimalDTO(_ctx, post)
		if err != nil {
			return nil, err
		}
		postMinimals = append(postMinimals, postMinimal)
	}
	return postMinimals, nil
}

func (p *PostHandler) ListPostsByStatus(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var postQueryNoEnum param.PostQueryNoEnum
	err := ctx.BindByContentType(&postQueryNoEnum)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if postQueryNoEnum.Sort == nil {
		postQueryNoEnum.Sort = &param.Sort{Fields: []string{"createTime,desc"}}
	}

	status, err := util.ParamInt32(_ctx, ctx, "status")
	if err != nil {
		return nil, err
	}
	postQuery := param.AssertPostQuery(postQueryNoEnum)
	postQuery.Statuses = make([]*consts.PostStatus, 0)
	statusType := consts.PostStatus(status)
	postQuery.Statuses = append(postQuery.Statuses, &statusType)

	posts, totalCount, err := p.PostService.Page(_ctx, postQuery)
	if err != nil {
		return nil, err
	}
	if postQuery.More == nil {
		*postQuery.More = false
	}
	if postQuery.More == nil {
		postVOs, err := p.PostAssembler.ConvertToListVO(_ctx, posts)
		return dto.NewPage(postVOs, totalCount, postQuery.Page), err
	}

	postDTOs := make([]*dto.Post, 0)
	for _, post := range posts {
		postDTO, err := p.PostAssembler.ConvertToSimpleDTO(_ctx, post)
		if err != nil {
			return nil, err
		}
		postDTOs = append(postDTOs, postDTO)
	}

	return dto.NewPage(postDTOs, totalCount, postQuery.Page), nil
}

func (p *PostHandler) GetByPostID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postIDStr := ctx.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 32)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	post, err := p.PostService.GetByPostID(_ctx, int32(postID))
	if err != nil {
		return nil, err
	}
	postDetailVO, err := p.PostAssembler.ConvertToDetailVO(_ctx, post)
	if err != nil {
		return nil, err
	}
	return postDetailVO, nil
}

func (p *PostHandler) CreatePost(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var postParam param.Post
	err := ctx.BindAndValidate(&postParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}

	post, err := p.PostService.Create(_ctx, &postParam)
	if err != nil {
		return nil, err
	}
	return p.PostAssembler.ConvertToDetailVO(_ctx, post)
}

func (p *PostHandler) UpdatePost(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var postParam param.Post
	err := ctx.BindAndValidate(&postParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}

	postIDStr := ctx.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 32)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}

	postDetailVO, err := p.PostService.Update(_ctx, int32(postID), &postParam)
	if err != nil {
		return nil, err
	}
	return postDetailVO, nil
}

func (p *PostHandler) UpdatePostStatus(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postIDStr := ctx.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 32)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	statusStr, err := util.ParamString(_ctx, ctx, "status")
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	status, err := consts.PostStatusFromString(statusStr)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if int32(status) < int32(consts.PostStatusPublished) || int32(status) > int32(consts.PostStatusIntimate) {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("status error")
	}
	post, err := p.PostService.UpdateStatus(_ctx, int32(postID), status)
	if err != nil {
		return nil, err
	}
	return p.PostAssembler.ConvertToMinimalDTO(_ctx, post)
}

func (p *PostHandler) UpdatePostStatusBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	statusStr, err := util.ParamString(_ctx, ctx, "status")
	if err != nil {
		return nil, err
	}
	status, err := consts.PostStatusFromString(statusStr)
	if err != nil {
		return nil, err
	}
	if int32(status) < int32(consts.PostStatusPublished) || int32(status) > int32(consts.PostStatusIntimate) {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("status error")
	}
	ids := make([]int32, 0)
	err = ctx.BindAndValidate(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}

	return p.PostService.UpdateStatusBatch(_ctx, status, ids)
}

func (p *PostHandler) UpdatePostDraft(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "postID")
	if err != nil {
		return nil, err
	}
	var postContentParam param.PostContent
	err = ctx.BindAndValidate(&postContentParam)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("content param error")
	}
	post, err := p.PostService.UpdateDraftContent(_ctx, postID, postContentParam.Content, postContentParam.OriginalContent)
	if err != nil {
		return nil, err
	}
	return p.PostAssembler.ConvertToDetailDTO(_ctx, post)
}

func (p *PostHandler) DeletePost(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "postID")
	if err != nil {
		return nil, err
	}
	return nil, p.PostService.Delete(_ctx, postID)
}

func (p *PostHandler) DeletePostBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postIDs := make([]int32, 0)
	err := ctx.BindAndValidate(&postIDs)
	if err != nil {
		return nil, xerr.WithMsg(err, "postIDs error").WithStatus(xerr.StatusBadRequest)
	}
	return nil, p.PostService.DeleteBatch(_ctx, postIDs)
}

func (p *PostHandler) PreviewPost(_ctx context.Context, ctx *app.RequestContext) {
	postID, err := util.ParamInt32(_ctx, ctx, "postID")
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		_ = ctx.Error(err)
		return
	}
	previewPath, err := p.PostService.Preview(_ctx, postID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		_ = ctx.Error(err)
		return
	}
	ctx.String(http.StatusOK, previewPath)
}
