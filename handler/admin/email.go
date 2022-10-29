package admin

import (
	"github.com/gin-gonic/gin"

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
func (e *EmailHandler) Test(ctx *gin.Context) (interface{}, error) {
	p := &param.TestEmail{}
	if err := ctx.ShouldBindJSON(p); err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("param error ")
	}
	return nil, e.EmailService.SendTextEmail(ctx, p.To, p.Subject, p.Content)
}
