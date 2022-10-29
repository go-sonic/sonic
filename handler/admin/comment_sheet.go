package admin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/binding"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type SheetCommentHandler struct {
	SheetCommentService   service.SheetCommentService
	BaseCommentService    service.BaseCommentService
	OptionService         service.OptionService
	SheetService          service.SheetService
	SheetAssembler        assembler.SheetAssembler
	SheetCommentAssembler assembler.SheetCommentAssembler
}

func NewSheetCommentHandler(
	sheetCommentService service.SheetCommentService,
	baseCommentService service.BaseCommentService,
	optionService service.OptionService,
	sheetService service.SheetService,
	sheetAssembler assembler.SheetAssembler,
	sheetCommentAssembler assembler.SheetCommentAssembler,
) *SheetCommentHandler {
	return &SheetCommentHandler{
		SheetCommentService:   sheetCommentService,
		BaseCommentService:    baseCommentService,
		OptionService:         optionService,
		SheetService:          sheetService,
		SheetAssembler:        sheetAssembler,
		SheetCommentAssembler: sheetCommentAssembler,
	}
}

func (s *SheetCommentHandler) ListSheetComment(ctx *gin.Context) (interface{}, error) {
	var commentQuery param.CommentQuery
	err := ctx.ShouldBindWith(&commentQuery, binding.CustomFormBinding)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	commentQuery.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	comments, totalCount, err := s.SheetCommentService.Page(ctx, commentQuery, consts.CommentTypeSheet)
	if err != nil {
		return nil, err
	}
	commentDTOs, err := s.ConvertToWithSheet(ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentDTOs, totalCount, commentQuery.Page), nil
}

func (s *SheetCommentHandler) ListSheetCommentLatest(ctx *gin.Context) (interface{}, error) {
	top, err := util.MustGetQueryInt32(ctx, "top")
	if err != nil {
		return nil, err
	}
	commentQuery := param.CommentQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: int(top)},
	}
	comments, _, err := s.SheetCommentService.Page(ctx, commentQuery, consts.CommentTypeSheet)
	if err != nil {
		return nil, err
	}
	return s.ConvertToWithSheet(ctx, comments)
}

func (s *SheetCommentHandler) ListSheetCommentAsTree(ctx *gin.Context) (interface{}, error) {
	postID, err := util.ParamInt32(ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	pageNum, err := util.MustGetQueryInt32(ctx, "page")
	if err != nil {
		return nil, err
	}
	pageSize, err := s.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return nil, err
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: int(pageNum)}

	allComments, err := s.SheetCommentService.GetByContentID(ctx, postID, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return nil, err
	}
	commentVOs, totalCount, err := s.SheetCommentAssembler.PageConvertToVOs(ctx, allComments, page)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentVOs, totalCount, page), nil
}

func (s *SheetCommentHandler) ListSheetCommentWithParent(ctx *gin.Context) (interface{}, error) {
	postID, err := util.ParamInt32(ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	pageNum, err := util.MustGetQueryInt32(ctx, "page")
	if err != nil {
		return nil, err
	}

	pageSize, err := s.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return nil, err
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: int(pageNum)}

	comments, totalCount, err := s.SheetCommentService.Page(ctx, param.CommentQuery{
		ContentID: &postID,
		Page:      page,
		Sort:      &param.Sort{Fields: []string{"createTime,desc"}},
	}, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}

	commentsWithParent, err := s.SheetCommentAssembler.ConvertToWithParentVO(ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentsWithParent, totalCount, page), nil
}

func (s *SheetCommentHandler) CreateSheetComment(ctx *gin.Context) (interface{}, error) {
	var commentParam *param.Comment
	err := ctx.ShouldBindJSON(&commentParam)
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
	commentParam.CommentType = consts.CommentTypeSheet
	comment, err := s.BaseCommentService.CreateBy(ctx, commentParam)
	if err != nil {
		return nil, err
	}
	return s.SheetCommentAssembler.ConvertToDTO(ctx, comment)
}

func (s *SheetCommentHandler) UpdateSheetCommentStatus(ctx *gin.Context) (interface{}, error) {
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
	return s.SheetCommentService.UpdateStatus(ctx, int64(commentID), status)
}

func (s *SheetCommentHandler) UpdateSheetCommentStatusBatch(ctx *gin.Context) (interface{}, error) {
	status, err := util.ParamInt32(ctx, "status")
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	err = ctx.ShouldBindJSON(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}
	comments, err := s.SheetCommentService.UpdateStatusBatch(ctx, ids, consts.CommentStatus(status))
	if err != nil {
		return nil, err
	}
	return s.SheetCommentAssembler.ConvertToDTOList(ctx, comments)
}

func (s *SheetCommentHandler) DeleteSheetComment(ctx *gin.Context) (interface{}, error) {
	commentID, err := util.ParamInt32(ctx, "commentID")
	if err != nil {
		return nil, err
	}
	return nil, s.SheetCommentService.Delete(ctx, int64(commentID))
}

func (s *SheetCommentHandler) DeleteSheetCommentBatch(ctx *gin.Context) (interface{}, error) {
	ids := make([]int64, 0)
	err := ctx.ShouldBindJSON(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}
	return nil, s.SheetCommentService.DeleteBatch(ctx, ids)
}

func (s *SheetCommentHandler) ConvertToWithSheet(ctx context.Context, comments []*entity.Comment) ([]*vo.SheetCommentWithSheet, error) {
	postIDs := make([]int32, 0, len(comments))
	for _, comment := range comments {
		postIDs = append(postIDs, comment.PostID)
	}
	posts, err := s.SheetService.GetByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	result := make([]*vo.SheetCommentWithSheet, 0, len(comments))
	for _, comment := range comments {
		commentDTO, err := s.SheetCommentAssembler.ConvertToDTO(ctx, comment)
		if err != nil {
			return nil, err
		}
		commentWithSheet := &vo.SheetCommentWithSheet{
			Comment: *commentDTO,
		}
		result = append(result, commentWithSheet)
		post, ok := posts[comment.PostID]
		if ok {
			commentWithSheet.PostMinimal, err = s.SheetAssembler.ConvertToMinimalDTO(ctx, post)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
