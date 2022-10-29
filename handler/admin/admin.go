package admin

import (
	"github.com/gin-gonic/gin"
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

func (a *AdminHandler) IsInstalled(ctx *gin.Context) (interface{}, error) {
	return a.OptionService.GetOrByDefaultWithErr(ctx, property.IsInstalled, false)
}

func (a *AdminHandler) AuthPreCheck(ctx *gin.Context) (interface{}, error) {
	var loginParam param.LoginParam
	err := ctx.ShouldBindJSON(&loginParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.BadParam.Wrapf(err, "")
	}

	user, err := a.AdminService.Authenticate(ctx, loginParam)
	if err != nil {
		return nil, err
	}
	return &dto.LoginPreCheckDTO{NeedMFACode: a.TwoFactorMFAService.UseMFA(user.MfaType)}, nil
}

func (a *AdminHandler) Auth(ctx *gin.Context) (interface{}, error) {
	var loginParam param.LoginParam
	err := ctx.ShouldBindJSON(&loginParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.BadParam.Wrapf(err, "").WithStatus(xerr.StatusBadRequest)
	}

	return a.AdminService.Auth(ctx, loginParam)
}

func (a *AdminHandler) LogOut(ctx *gin.Context) (interface{}, error) {
	err := a.AdminService.ClearToken(ctx)
	return nil, err
}

func (a *AdminHandler) SendResetCode(ctx *gin.Context) (interface{}, error) {
	var resetPasswordParam param.ResetPasswordParam
	err := ctx.ShouldBindJSON(&resetPasswordParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.BadParam.Wrapf(err, "").WithStatus(xerr.StatusBadRequest)
	}
	return nil, a.AdminService.SendResetPasswordCode(ctx, resetPasswordParam)
}

func (a *AdminHandler) RefreshToken(ctx *gin.Context) (interface{}, error) {
	refreshToken := ctx.Param("refreshToken")
	if refreshToken == "" {
		return nil, xerr.BadParam.New("refreshToken参数为空").WithStatus(xerr.StatusBadRequest).
			WithMsg("refreshToken 参数不能为空")
	}
	return a.AdminService.RefreshToken(ctx, refreshToken)
}

func (a *AdminHandler) GetEnvironments(ctx *gin.Context) (interface{}, error) {
	return a.AdminService.GetEnvironments(ctx), nil
}

func (a *AdminHandler) GetLogFiles(ctx *gin.Context) (interface{}, error) {
	lines, err := util.MustGetQueryInt64(ctx, "lines")
	if err != nil {
		return nil, err
	}
	return a.AdminService.GetLogFiles(ctx, lines)
}
