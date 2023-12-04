package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type ApplicationPasswordHandler struct {
	ApplicationPasswordService service.ApplicationPasswordService
}

func NewApplicationPasswordHandler(applicationPasswordService service.ApplicationPasswordService) *ApplicationPasswordHandler {
	return &ApplicationPasswordHandler{
		ApplicationPasswordService: applicationPasswordService,
	}
}

func (a *ApplicationPasswordHandler) Create(ctx *gin.Context) (interface{}, error) {
	appPwdParam, err := parseAppPwdParam(ctx)
	if err != nil {
		return nil, err
	}

	return a.ApplicationPasswordService.CreatePwd(ctx, appPwdParam)
}

func (a *ApplicationPasswordHandler) Delete(ctx *gin.Context) (interface{}, error) {
	name, err := util.ParamString(ctx, "name")
	if err != nil {
		return nil, err
	}

	appPwdParam := &param.ApplicationPasswordParam{Name: name}

	return nil, a.ApplicationPasswordService.DeletePwd(ctx, appPwdParam)
}

func (a *ApplicationPasswordHandler) List(ctx *gin.Context) (interface{}, error) {
	return a.ApplicationPasswordService.List(ctx)
}

func parseAppPwdParam(ctx *gin.Context) (*param.ApplicationPasswordParam, error) {
	var appPwdParam param.ApplicationPasswordParam
	err := ctx.ShouldBindJSON(&appPwdParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	return &appPwdParam, nil
}
