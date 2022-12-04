package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/binding"
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

func (p *PostCommentHandler) ListPostComment(ctx *gin.Context) (interface{}, error) {
	var commentQuery param.CommentQuery
	err := ctx.ShouldBindWith(&commentQuery, binding.CustomFormBinding)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	commentQuery.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	comments, totalCount, err := p.PostCommentService.Page(ctx, commentQuery, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}
	commentDTOs, err := p.PostCommentAssembler.ConvertToWithPost(ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentDTOs, totalCount, commentQuery.Page), nil
}

func (p *PostCommentHandler) ListPostCommentLatest(ctx *gin.Context) (interface{}, error) {
	top, err := util.MustGetQueryInt32(ctx, "top")
	if err != nil {
		return nil, err
	}
	commentQuery := param.CommentQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: int(top)},
	}
	comments, _, err := p.PostCommentService.Page(ctx, commentQuery, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}
	return p.PostCommentAssembler.ConvertToWithPost(ctx, comments)
}

func (p *PostCommentHandler) ListPostCommentAsTree(ctx *gin.Context) (interface{}, error) {
	postID, err := util.ParamInt32(ctx, "postID")
	if err != nil {
		return nil, err
	}
	pageNum, err := util.MustGetQueryInt32(ctx, "page")
	if err != nil {
		return nil, err
	}
	pageSize, err := p.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return nil, err
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: int(pageNum)}
	allComments, err := p.PostCommentService.GetByContentID(ctx, postID, consts.CommentTypePost, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return nil, err
	}
	commentVOs, totalCount, err := p.PostCommentAssembler.PageConvertToVOs(ctx, allComments, page)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentVOs, totalCount, page), nil
}

func (p *PostCommentHandler) ListPostCommentWithParent(ctx *gin.Context) (interface{}, error) {
	postID, err := util.ParamInt32(ctx, "postID")
	if err != nil {
		return nil, err
	}
	pageNum, err := util.MustGetQueryInt32(ctx, "page")
	if err != nil {
		return nil, err
	}

	pageSize, err := p.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return nil, err
	}

	page := param.Page{PageNum: int(pageNum), PageSize: pageSize.(int)}

	comments, totalCount, err := p.PostCommentService.Page(ctx, param.CommentQuery{
		ContentID: &postID,
		Page:      page,
		Sort:      &param.Sort{Fields: []string{"createTime,desc"}},
	}, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}

	commentsWithParent, err := p.PostCommentAssembler.ConvertToWithParentVO(ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentsWithParent, totalCount, page), nil
}

func (p *PostCommentHandler) CreatePostComment(ctx *gin.Context) (interface{}, error) {
	var commentParam *param.AdminComment
	err := ctx.ShouldBindJSON(&commentParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	user, err := impl.MustGetAuthorizedUser(ctx)
	if err != nil || user == nil {
		return nil, err
	}
	blogURL, err := p.OptionService.GetBlogBaseURL(ctx)
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
	comment, err := p.PostCommentService.CreateBy(ctx, &commonParam)
	if err != nil {
		return nil, err
	}
	return p.PostCommentAssembler.ConvertToDTO(ctx, comment)
}

func (p *PostCommentHandler) UpdatePostComment(ctx *gin.Context) (interface{}, error) {
	commentID, err := util.ParamInt32(ctx, "commentID")
	if err != nil {
		return nil, err
	}
	var commentParam *param.Comment
	err = ctx.ShouldBindJSON(&commentParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
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
	comment, err := p.PostCommentService.UpdateBy(ctx, commentID, commentParam)
	if err != nil {
		return nil, err
	}

	return p.PostCommentAssembler.ConvertToDTO(ctx, comment)
}

func (p *PostCommentHandler) UpdatePostCommentStatus(ctx *gin.Context) (interface{}, error) {
	commentID, err := util.ParamInt32(ctx, "commentID")
	if err != nil {
		return nil, err
	}
	strStatus, err := util.ParamString(ctx, "status")
	if err != nil {
		return nil, err
	}
	status, err := consts.CommentStatusFromString(strStatus)
	if err != nil {
		return nil, err
	}
	return p.PostCommentService.UpdateStatus(ctx, commentID, status)
}

func (p *PostCommentHandler) UpdatePostCommentStatusBatch(ctx *gin.Context) (interface{}, error) {
	strStatus, err := util.ParamString(ctx, "status")
	if err != nil {
		return nil, err
	}
	status, err := consts.CommentStatusFromString(strStatus)
	if err != nil {
		return nil, err
	}

	ids := make([]int32, 0)
	err = ctx.ShouldBindJSON(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}
	comments, err := p.PostCommentService.UpdateStatusBatch(ctx, ids, status)
	if err != nil {
		return nil, err
	}
	return p.PostCommentAssembler.ConvertToDTOList(ctx, comments)
}

func (p *PostCommentHandler) DeletePostComment(ctx *gin.Context) (interface{}, error) {
	commentID, err := util.ParamInt32(ctx, "commentID")
	if err != nil {
		return nil, err
	}
	return nil, p.PostCommentService.Delete(ctx, commentID)
}

func (p *PostCommentHandler) DeletePostCommentBatch(ctx *gin.Context) (interface{}, error) {
	ids := make([]int32, 0)
	err := ctx.ShouldBindJSON(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}
	return nil, p.PostCommentService.DeleteBatch(ctx, ids)
}
