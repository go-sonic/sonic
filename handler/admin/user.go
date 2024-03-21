package admin

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/impl"
	"github.com/go-sonic/sonic/util/xerr"
)

type UserHandler struct {
	UserService         service.UserService
	TwoFactorMFAService service.TwoFactorTOTPMFAService
}

func NewUserHandler(userService service.UserService, twoFactorMFAService service.TwoFactorTOTPMFAService) *UserHandler {
	return &UserHandler{
		UserService:         userService,
		TwoFactorMFAService: twoFactorMFAService,
	}
}

func (u *UserHandler) GetCurrentUserProfile(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	user, ok := impl.GetAuthorizedUser(_ctx)
	if !ok {
		return nil, xerr.Forbidden.New("authorized user nil").WithStatus(xerr.StatusForbidden)
	}
	return u.UserService.ConvertToDTO(_ctx, user), nil
}

func (u *UserHandler) UpdateUserProfile(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	userParam := &param.User{}
	err := ctx.BindAndValidate(userParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	user, err := u.UserService.Update(_ctx, userParam)
	if err != nil {
		return nil, err
	}
	return u.UserService.ConvertToDTO(_ctx, user), nil
}

func (u *UserHandler) UpdatePassword(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	type Password struct {
		OldPassword string `json:"oldPassword" form:"oldPassword" binding:"gte=1,lte=100"`
		NewPassword string `json:"newPassword" form:"newPassword" binding:"gte=1,lte=100"`
	}
	passwordParam := &Password{}
	err := ctx.BindAndValidate(passwordParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	return nil, u.UserService.UpdatePassword(_ctx, passwordParam.OldPassword, passwordParam.NewPassword)
}

func (u *UserHandler) GenerateMFAQRCode(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	type Param struct {
		MFAType *consts.MFAType `json:"mfaType"`
	}
	param := &Param{}
	err := ctx.BindAndValidate(param)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	if param.MFAType == nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	user, ok := impl.GetAuthorizedUser(_ctx)
	if !ok || user == nil {
		return nil, xerr.Forbidden.New("").WithMsg("unauthorized").WithStatus(xerr.StatusForbidden)
	}

	mfaFactorAuthDTO := &vo.MFAFactorAuth{}
	if *param.MFAType == consts.MFATFATotp {
		key, url, err := u.TwoFactorMFAService.GenerateOTPKey(_ctx, user.Nickname)
		if err != nil {
			return nil, err
		}
		mfaFactorAuthDTO.MFAType = consts.MFATFATotp
		mfaFactorAuthDTO.OptAuthURL = url
		mfaFactorAuthDTO.MFAKey = key
		qrCode, err := u.TwoFactorMFAService.GenerateMFAQRCode(_ctx, url)
		if err != nil {
			return nil, err
		}
		mfaFactorAuthDTO.QRImage = qrCode
		return mfaFactorAuthDTO, nil
	} else {
		return nil, xerr.WithMsg(nil, "Not supported authentication").WithStatus(xerr.StatusBadRequest)
	}
}

func (u *UserHandler) UpdateMFA(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	type Param struct {
		MFAType  *consts.MFAType `json:"mfaType" form:"mfaType"`
		MFAKey   string          `json:"mfaKey" form:"mfaKey"`
		AuthCode string          `json:"authcode" form:"authcode" binding:"gte=6,lte=6"`
	}
	mfaParam := &Param{}
	err := ctx.BindAndValidate(mfaParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	if mfaParam.MFAType == nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	return nil, u.UserService.UpdateMFA(_ctx, mfaParam.MFAKey, *mfaParam.MFAType, mfaParam.AuthCode)
}
