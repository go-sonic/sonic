package admin

import (
	"context"
	"errors"
	"net/http"

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

func (s *SheetHandler) GetSheetByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	sheetID, err := util.ParamInt32(_ctx, ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	sheet, err := s.SheetService.GetByPostID(_ctx, sheetID)
	if err != nil {
		return nil, err
	}
	return s.SheetAssembler.ConvertToDetailVO(_ctx, sheet)
}

func (s *SheetHandler) ListSheet(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	type SheetParam struct {
		param.Page
		Sort string `json:"sort"`
	}
	var sheetParam SheetParam
	err := ctx.BindAndValidate(&sheetParam)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	sheets, totalCount, err := s.SheetService.Page(_ctx, sheetParam.Page, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return nil, err
	}
	sheetVOs, err := s.SheetAssembler.ConvertToListVO(_ctx, sheets)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(sheetVOs, totalCount, sheetParam.Page), nil
}

func (s *SheetHandler) IndependentSheets(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return s.SheetService.ListIndependentSheets(_ctx)
}

func (s *SheetHandler) CreateSheet(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var sheetParam param.Sheet
	err := ctx.BindAndValidate(&sheetParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	sheet, err := s.SheetService.Create(_ctx, &sheetParam)
	if err != nil {
		return nil, err
	}
	sheetDetailVO, err := s.SheetAssembler.ConvertToDetailVO(_ctx, sheet)
	if err != nil {
		return nil, err
	}
	return sheetDetailVO, nil
}

func (s *SheetHandler) UpdateSheet(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var sheetParam param.Sheet
	err := ctx.BindAndValidate(&sheetParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}

	sheetID, err := util.ParamInt32(_ctx, ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	postDetailVO, err := s.SheetService.Update(_ctx, sheetID, &sheetParam)
	if err != nil {
		return nil, err
	}
	return postDetailVO, nil
}

func (s *SheetHandler) UpdateSheetStatus(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	sheetID, err := util.ParamInt32(_ctx, ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	statusStr, err := util.ParamString(_ctx, ctx, "status")
	if err != nil {
		return nil, err
	}
	status, err := consts.PostStatusFromString(statusStr)
	if err != nil {
		return nil, err
	}
	if status < consts.PostStatusPublished || status > consts.PostStatusIntimate {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("status error")
	}
	return s.SheetService.UpdateStatus(_ctx, sheetID, status)
}

func (s *SheetHandler) UpdateSheetDraft(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	sheetID, err := util.ParamInt32(_ctx, ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	var postContentParam param.PostContent
	err = ctx.BindAndValidate(&postContentParam)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("content param error")
	}
	post, err := s.SheetService.UpdateDraftContent(_ctx, sheetID, postContentParam.Content, postContentParam.OriginalContent)
	if err != nil {
		return nil, err
	}
	return s.SheetAssembler.ConvertToDetailDTO(_ctx, post)
}

func (s *SheetHandler) DeleteSheet(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	sheetID, err := util.ParamInt32(_ctx, ctx, "sheetID")
	if err != nil {
		return nil, err
	}
	return nil, s.SheetService.Delete(_ctx, sheetID)
}

func (s *SheetHandler) PreviewSheet(_ctx context.Context, ctx *app.RequestContext) {
	sheetID, err := util.ParamInt32(_ctx, ctx, "sheetID")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		_ = ctx.Error(err)
		return
	}

	previewPath, err := s.SheetService.Preview(_ctx, sheetID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		_ = ctx.Error(err)
		return
	}
	ctx.String(http.StatusOK, previewPath)
}
