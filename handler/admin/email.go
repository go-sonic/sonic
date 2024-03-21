package admin

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type EmailHandler struct {
	EmailService service.EmailService
}

func NewEmailHandler(emailService service.EmailService) *EmailHandler {
	return &EmailHandler{
		EmailService: emailService,
	}
}

func (e *EmailHandler) Test(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	p := &param.TestEmail{}
	if err := ctx.BindAndValidate(p); err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("param error ")
	}
	return nil, e.EmailService.SendTextEmail(_ctx, p.To, p.Subject, p.Content)
}
