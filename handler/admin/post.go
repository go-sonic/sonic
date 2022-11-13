package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/binding"
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

func (p *PostHandler) ListPosts(ctx *gin.Context) (interface{}, error) {
	postQuery := param.PostQuery{}
	err := ctx.ShouldBindWith(&postQuery, binding.CustomFormBinding)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if postQuery.Sort == nil {
		postQuery.Sort = &param.Sort{Fields: []string{"topPriority,desc", "createTime,desc"}}
	}
	posts, totalCount, err := p.PostService.Page(ctx, postQuery)
	if err != nil {
		return nil, err
	}
	if postQuery.More == nil || *postQuery.More {
		postVOs, err := p.PostAssembler.ConvertToListVO(ctx, posts)
		return dto.NewPage(postVOs, totalCount, postQuery.Page), err
	}
	postDTOs := make([]*dto.Post, 0)
	for _, post := range posts {
		postDTO, err := p.PostAssembler.ConvertToSimpleDTO(ctx, post)
		if err != nil {
			return nil, err
		}
		postDTOs = append(postDTOs, postDTO)
	}
	return dto.NewPage(postDTOs, totalCount, postQuery.Page), nil
}

func (p *PostHandler) ListLatestPosts(ctx *gin.Context) (interface{}, error) {
	top, err := util.MustGetQueryInt32(ctx, "top")
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
	posts, _, err := p.PostService.Page(ctx, postQuery)
	if err != nil {
		return nil, err
	}
	postMinimals := make([]*dto.PostMinimal, 0, len(posts))

	for _, post := range posts {
		postMinimal, err := p.PostAssembler.ConvertToMinimalDTO(ctx, post)
		if err != nil {
			return nil, err
		}
		postMinimals = append(postMinimals, postMinimal)
	}
	return postMinimals, nil
}

func (p *PostHandler) ListPostsByStatus(ctx *gin.Context) (interface{}, error) {
	var postQuery param.PostQuery
	err := ctx.ShouldBindWith(&postQuery, binding.CustomFormBinding)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if postQuery.Sort == nil {
		postQuery.Sort = &param.Sort{Fields: []string{"createTime,desc"}}
	}

	status, err := util.ParamInt32(ctx, "status")
	if err != nil {
		return nil, err
	}
	postQuery.Statuses = make([]*consts.PostStatus, 0)
	statusType := consts.PostStatus(status)
	postQuery.Statuses = append(postQuery.Statuses, &statusType)

	posts, totalCount, err := p.PostService.Page(ctx, postQuery)
	if err != nil {
		return nil, err
	}
	if postQuery.More == nil {
		*postQuery.More = false
	}
	if postQuery.More == nil {
		postVOs, err := p.PostAssembler.ConvertToListVO(ctx, posts)
		return dto.NewPage(postVOs, totalCount, postQuery.Page), err
	}

	postDTOs := make([]*dto.Post, 0)
	for _, post := range posts {
		postDTO, err := p.PostAssembler.ConvertToSimpleDTO(ctx, post)
		if err != nil {
			return nil, err
		}
		postDTOs = append(postDTOs, postDTO)
	}

	return dto.NewPage(postDTOs, totalCount, postQuery.Page), nil
}

func (p *PostHandler) GetByPostID(ctx *gin.Context) (interface{}, error) {
	postIDStr := ctx.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 32)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	post, err := p.PostService.GetByPostID(ctx, int32(postID))
	if err != nil {
		return nil, err
	}
	postDetailVO, err := p.PostAssembler.ConvertToDetailVO(ctx, post)
	if err != nil {
		return nil, err
	}
	return postDetailVO, nil
}

func (p *PostHandler) LikePost(ctx *gin.Context) (interface{}, error) {
	postID, err := util.ParamInt32(ctx, "postID")
	if err != nil {
		return nil, err
	}
	return nil, p.PostService.IncreaseLike(ctx, int32(postID))
}

func (p *PostHandler) CreatePost(ctx *gin.Context) (interface{}, error) {
	var postParam param.Post
	err := ctx.ShouldBindJSON(&postParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}

	post, err := p.PostService.Create(ctx, &postParam)
	if err != nil {
		return nil, err
	}
	return p.PostAssembler.ConvertToDetailVO(ctx, post)
}

func (p *PostHandler) UpdatePost(ctx *gin.Context) (interface{}, error) {
	var postParam param.Post
	err := ctx.ShouldBindJSON(&postParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}

	postIDStr := ctx.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 32)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}

	postDetailVO, err := p.PostService.Update(ctx, int32(postID), &postParam)
	if err != nil {
		return nil, err
	}
	return postDetailVO, nil
}

func (p *PostHandler) UpdatePostStatus(ctx *gin.Context) (interface{}, error) {
	postIDStr := ctx.Param("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 32)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	statusStr, err := util.ParamString(ctx, "status")
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
	post, err := p.PostService.UpdateStatus(ctx, int32(postID), consts.PostStatus(status))
	if err != nil {
		return nil, err
	}
	return p.PostAssembler.ConvertToMinimalDTO(ctx, post)
}

func (p *PostHandler) UpdatePostStatusBatch(ctx *gin.Context) (interface{}, error) {
	statusStr, err := util.ParamString(ctx, "status")
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
	err = ctx.ShouldBind(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}

	return p.PostService.UpdateStatusBatch(ctx, consts.PostStatus(status), ids)
}

func (p *PostHandler) UpdatePostDraft(ctx *gin.Context) (interface{}, error) {
	postID, err := util.ParamInt32(ctx, "postID")
	if err != nil {
		return nil, err
	}
	var postContentParam param.PostContent
	err = ctx.ShouldBindJSON(&postContentParam)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("content param error")
	}
	post, err := p.PostService.UpdateDraftContent(ctx, int32(postID), postContentParam.Content)
	if err != nil {
		return nil, err
	}
	return p.PostAssembler.ConvertToDetailDTO(ctx, post)
}

func (p *PostHandler) DeletePost(ctx *gin.Context) (interface{}, error) {
	postID, err := util.ParamInt32(ctx, "postID")
	if err != nil {
		return nil, err
	}
	return nil, p.PostService.Delete(ctx, postID)
}

func (p *PostHandler) DeletePostBatch(ctx *gin.Context) (interface{}, error) {
	postIDs := make([]int32, 0)
	err := ctx.ShouldBind(&postIDs)
	if err != nil {
		return nil, xerr.WithMsg(err, "postIDs error").WithStatus(xerr.StatusBadRequest)
	}
	return nil, p.PostService.DeleteBatch(ctx, postIDs)
}

func (p *PostHandler) PreviewPost(ctx *gin.Context) (interface{}, error) {
	postID, err := util.ParamInt32(ctx, "postID")
	if err != nil {
		return nil, err
	}
	return p.PostService.Preview(ctx, int32(postID))
}
