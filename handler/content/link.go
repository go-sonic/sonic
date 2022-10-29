package content

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/handler/content/model"
	"github.com/go-sonic/sonic/template"
)

type LinkHandler struct {
	LinkModel *model.LinkModel
}

func NewLinkHandler(
	linkModel *model.LinkModel,
) *LinkHandler {
	return &LinkHandler{
		LinkModel: linkModel,
	}
}

func (t *LinkHandler) Link(ctx *gin.Context, model template.Model) (string, error) {
	return t.LinkModel.Links(ctx, model)
}
