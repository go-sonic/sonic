package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
)

type PhotoHandler struct {
	PhotoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		PhotoService: photoService,
	}
}

func (p *PhotoHandler) Like(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "photoID")
	if err != nil {
		return nil, err
	}
	return nil, p.PhotoService.IncreaseLike(_ctx, id)
}
