package admin

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type InstallHandler struct {
	InstallService service.InstallService
}

func NewInstallHandler(installService service.InstallService) *InstallHandler {
	return &InstallHandler{
		InstallService: installService,
	}
}

func (i *InstallHandler) InstallBlog(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var installParam param.Install
	err := ctx.BindAndValidate(&installParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	err = i.InstallService.InstallBlog(_ctx, installParam)
	if err != nil {
		return nil, err
	}
	return "安装完成", nil
}
