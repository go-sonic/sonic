package content

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
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

func (t *LinkHandler) Link(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	return t.LinkModel.Links(_ctx, model)
}
