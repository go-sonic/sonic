package content

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/cache"
	"github.com/go-sonic/sonic/handler/content/model"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type SheetHandler struct {
	OptionService service.OptionService
	SheetService  service.SheetService
	SheetModel    *model.SheetModel
	Cache         cache.Cache
}

func NewSheetHandler(
	optionService service.OptionService,
	sheetService service.SheetService,
	sheetModel *model.SheetModel,
	cache cache.Cache,
) *SheetHandler {
	return &SheetHandler{
		OptionService: optionService,
		SheetService:  sheetService,
		SheetModel:    sheetModel,
		Cache:         cache,
	}
}

func (s *SheetHandler) SheetBySlug(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	slug, err := util.ParamString(_ctx, ctx, "slug")
	if err != nil {
		return "", err
	}
	sheet, err := s.SheetService.GetBySlug(_ctx, slug)
	if err != nil {
		return "", err
	}
	token := string(ctx.Cookie("authentication"))
	return s.SheetModel.Content(_ctx, sheet, token, model)
}

func (s *SheetHandler) AdminSheetBySlug(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	slug, err := util.ParamString(_ctx, ctx, "slug")
	if err != nil {
		return "", err
	}

	token, err := util.MustGetQueryString(_ctx, ctx, "token")
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", nil
	}
	_, ok := s.Cache.Get(token)
	if !ok {
		return "", xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("token已过期或者不存在")
	}

	sheet, err := s.SheetService.GetBySlug(_ctx, slug)
	if err != nil {
		return "", err
	}

	return s.SheetModel.AdminPreviewContent(_ctx, sheet, model)
}
