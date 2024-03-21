package admin

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/service/impl"
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

func (s *SheetCommentHandler) ListSheetComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var commentQueryNoEnum param.CommentQueryNoEnum
	err := ctx.BindAndValidate(&commentQueryNoEnum)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	commentQueryNoEnum.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	commentQuery := param.AssertCommentQuery(commentQueryNoEnum)
	comments, totalCount, err := s.SheetCommentService.Page(_ctx, commentQuery, consts.CommentTypeSheet)
	if err != nil {
		return nil, err
	}
	commentDTOs, err := s.ConvertToWithSheet(_ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentDTOs, totalCount, commentQuery.Page), nil
}

func (s *SheetCommentHandler) ListSheetCommentLatest(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	top, err := util.MustGetQueryInt32(_ctx, ctx, "top")
	if err != nil {
		return nil, err
	}
	commentQuery := param.CommentQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: int(top)},
	}
	comments, _, err := s.SheetCommentService.Page(_ctx, commentQuery, consts.CommentTypeSheet)
	if err != nil {
		return nil, err
	}
	return s.ConvertToWithSheet(_ctx, comments)
}

func (s *SheetCommentHandler) ListSheetCommentAsTree(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	pageNum, err := util.MustGetQueryInt32(_ctx, ctx, "page")
	if err != nil {
		return nil, err
	}
	pageSize, err := s.OptionService.GetOrByDefaultWithErr(_ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return nil, err
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: int(pageNum)}

	allComments, err := s.SheetCommentService.GetByContentID(_ctx, postID, consts.CommentTypeSheet, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return nil, err
	}
	commentVOs, totalCount, err := s.SheetCommentAssembler.PageConvertToVOs(_ctx, allComments, page)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentVOs, totalCount, page), nil
}

func (s *SheetCommentHandler) ListSheetCommentWithParent(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	postID, err := util.ParamInt32(_ctx, ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	pageNum, err := util.MustGetQueryInt32(_ctx, ctx, "page")
	if err != nil {
		return nil, err
	}

	pageSize, err := s.OptionService.GetOrByDefaultWithErr(_ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return nil, err
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: int(pageNum)}

	comments, totalCount, err := s.SheetCommentService.Page(_ctx, param.CommentQuery{
		ContentID: &postID,
		Page:      page,
		Sort:      &param.Sort{Fields: []string{"createTime,desc"}},
	}, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}

	commentsWithParent, err := s.SheetCommentAssembler.ConvertToWithParentVO(_ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentsWithParent, totalCount, page), nil
}

func (s *SheetCommentHandler) CreateSheetComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
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
	blogURL, err := s.OptionService.GetBlogBaseURL(_ctx)
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
		CommentType:       consts.CommentTypeSheet,
	}
	comment, err := s.BaseCommentService.CreateBy(_ctx, &commonParam)
	if err != nil {
		return nil, err
	}
	return s.SheetCommentAssembler.ConvertToDTO(_ctx, comment)
}

func (s *SheetCommentHandler) UpdateSheetCommentStatus(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
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
	return s.SheetCommentService.UpdateStatus(_ctx, commentID, status)
}

func (s *SheetCommentHandler) UpdateSheetCommentStatusBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	status, err := util.ParamInt32(_ctx, ctx, "status")
	if err != nil {
		return nil, err
	}

	ids := make([]int32, 0)
	err = ctx.BindAndValidate(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}
	comments, err := s.SheetCommentService.UpdateStatusBatch(_ctx, ids, consts.CommentStatus(status))
	if err != nil {
		return nil, err
	}
	return s.SheetCommentAssembler.ConvertToDTOList(_ctx, comments)
}

func (s *SheetCommentHandler) DeleteSheetComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	commentID, err := util.ParamInt32(_ctx, ctx, "commentID")
	if err != nil {
		return nil, err
	}
	return nil, s.SheetCommentService.Delete(_ctx, commentID)
}

func (s *SheetCommentHandler) DeleteSheetCommentBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	ids := make([]int32, 0)
	err := ctx.BindAndValidate(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}
	return nil, s.SheetCommentService.DeleteBatch(_ctx, ids)
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
