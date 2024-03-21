package admin

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type AdminHandler struct {
	OptionService       service.OptionService
	AdminService        service.AdminService
	TwoFactorMFAService service.TwoFactorTOTPMFAService
}

func NewAdminHandler(optionService service.OptionService, adminService service.AdminService, twoFactorMFA service.TwoFactorTOTPMFAService) *AdminHandler {
	return &AdminHandler{
		OptionService:       optionService,
		AdminService:        adminService,
		TwoFactorMFAService: twoFactorMFA,
	}
}

func (a *AdminHandler) IsInstalled(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return a.OptionService.GetOrByDefaultWithErr(_ctx, property.IsInstalled, false)
}

func (a *AdminHandler) AuthPreCheck(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var loginParam param.LoginParam
	err := ctx.BindAndValidate(&loginParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.BadParam.Wrapf(err, "")
	}

	user, err := a.AdminService.Authenticate(_ctx, loginParam)
	if err != nil {
		return nil, err
	}
	return &dto.LoginPreCheckDTO{NeedMFACode: a.TwoFactorMFAService.UseMFA(user.MfaType)}, nil
}

func (a *AdminHandler) Auth(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var loginParam param.LoginParam
	err := ctx.BindAndValidate(&loginParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.BadParam.Wrapf(err, "").WithStatus(xerr.StatusBadRequest)
	}

	return a.AdminService.Auth(_ctx, loginParam)
}

func (a *AdminHandler) LogOut(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	err := a.AdminService.ClearToken(_ctx)
	return nil, err
}

func (a *AdminHandler) SendResetCode(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var resetPasswordParam param.ResetPasswordParam
	err := ctx.BindAndValidate(&resetPasswordParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.BadParam.Wrapf(err, "").WithStatus(xerr.StatusBadRequest)
	}
	return nil, a.AdminService.SendResetPasswordCode(_ctx, resetPasswordParam)
}

func (a *AdminHandler) RefreshToken(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	refreshToken := ctx.Param("refreshToken")
	if refreshToken == "" {
		return nil, xerr.BadParam.New("refreshToken参数为空").WithStatus(xerr.StatusBadRequest).
			WithMsg("refreshToken 参数不能为空")
	}
	return a.AdminService.RefreshToken(_ctx, refreshToken)
}

func (a *AdminHandler) GetEnvironments(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return a.AdminService.GetEnvironments(_ctx), nil
}

func (a *AdminHandler) GetLogFiles(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	lines, err := util.MustGetQueryInt64(_ctx, ctx, "lines")
	if err != nil {
		return nil, err
	}
	return a.AdminService.GetLogFiles(_ctx, lines)
}
