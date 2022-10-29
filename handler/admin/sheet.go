package admin

import (
	"github.com/gin-gonic/gin"
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

type SheetHandler struct {
	SheetService   service.SheetService
	PostService    service.PostService
	SheetAssembler assembler.SheetAssembler
}

func NewSheetHandler(sheetService service.SheetService, postService service.PostService, sheetAssembler assembler.SheetAssembler) *SheetHandler {
	return &SheetHandler{
		SheetService:   sheetService,
		PostService:    postService,
		SheetAssembler: sheetAssembler,
	}
}

func (s *SheetHandler) GetSheetByID(ctx *gin.Context) (interface{}, error) {
	sheetID, err := util.ParamInt32(ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	sheet, err := s.SheetService.GetByPostID(ctx, sheetID)
	if err != nil {
		return nil, err
	}
	return s.SheetAssembler.ConvertToDetailVO(ctx, sheet)
}

func (s *SheetHandler) ListSheet(ctx *gin.Context) (interface{}, error) {
	type SheetParam struct {
		param.Page
		Sort string `json:"sort"`
	}
	var sheetParam SheetParam
	err := ctx.ShouldBind(&sheetParam)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	sheets, totalCount, err := s.SheetService.Page(ctx, sheetParam.Page, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return nil, err
	}
	sheetVOs, err := s.SheetAssembler.ConvertToListVO(ctx, sheets)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(sheetVOs, totalCount, sheetParam.Page), nil
}

func (s *SheetHandler) IndependentSheets(ctx *gin.Context) (interface{}, error) {
	return s.SheetService.ListIndependentSheets(ctx)
}

func (s *SheetHandler) CreateSheet(ctx *gin.Context) (interface{}, error) {
	var sheetParam param.Sheet
	err := ctx.ShouldBindJSON(&sheetParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	sheet, err := s.SheetService.Create(ctx, &sheetParam)
	if err != nil {
		return nil, err
	}
	sheetDetailVO, err := s.SheetAssembler.ConvertToDetailVO(ctx, sheet)
	if err != nil {
		return nil, err
	}
	return sheetDetailVO, nil
}

func (s *SheetHandler) UpdateSheet(ctx *gin.Context) (interface{}, error) {
	var sheetParam param.Sheet
	err := ctx.ShouldBindJSON(&sheetParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}

	sheetID, err := util.ParamInt32(ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	postDetailVO, err := s.SheetService.Update(ctx, sheetID, &sheetParam)
	if err != nil {
		return nil, err
	}
	return postDetailVO, nil
}

func (s *SheetHandler) UpdateSheetStatus(ctx *gin.Context) (interface{}, error) {
	sheetID, err := util.MustGetQueryInt32(ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	status, err := util.ParamInt32(ctx, "status")
	if err != nil {
		return nil, err
	}
	if status < int32(consts.PostStatusPublished) || status > int32(consts.PostStatusIntimate) {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("status error")
	}
	return s.SheetService.UpdateStatus(ctx, sheetID, consts.PostStatus(status))
}

func (s *SheetHandler) UpdateSheetDraft(ctx *gin.Context) (interface{}, error) {
	sheetID, err := util.MustGetQueryInt32(ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	var postContentParam param.PostContent
	err = ctx.ShouldBindJSON(&postContentParam)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("content param error")
	}
	post, err := s.SheetService.UpdateDraftContent(ctx, int32(sheetID), postContentParam.Content)
	if err != nil {
		return nil, err
	}
	return s.SheetAssembler.ConvertToDetailDTO(ctx, post)
}

func (s *SheetHandler) DeleteSheet(ctx *gin.Context) (interface{}, error) {
	sheetID, err := util.MustGetQueryInt32(ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	return nil, s.SheetService.Delete(ctx, int32(sheetID))
}

func (s *SheetHandler) PreviewSheet(ctx *gin.Context) (interface{}, error) {
	sheetID, err := util.ParamInt32(ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	return s.PostService.Preview(ctx, int32(sheetID))
}
