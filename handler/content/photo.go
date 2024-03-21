package content

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/handler/content/model"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
)

type PhotoHandler struct {
	OptionService service.OptionService
	PhotoService  service.PhotoService
	PhotoModel    *model.PhotoModel
}

func NewPhotoHandler(
	optionService service.OptionService,
	photoService service.PhotoService,
	photoModel *model.PhotoModel,
) *PhotoHandler {
	return &PhotoHandler{
		OptionService: optionService,
		PhotoService:  photoService,
		PhotoModel:    photoModel,
	}
}

func (p *PhotoHandler) PhotosPage(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	page, err := util.ParamInt32(_ctx, ctx, "page")
	if err != nil {
		return "", err
	}
	return p.PhotoModel.Photos(_ctx, model, int(page-1))
}

func (p *PhotoHandler) Phtotos(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	return p.PhotoModel.Photos(_ctx, model, 0)
}
