package content

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/handler/content/model"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
)

type IndexHandler struct {
	PostModel *model.PostModel
}

func NewIndexHandler(postModel *model.PostModel) *IndexHandler {
	return &IndexHandler{
		PostModel: postModel,
	}
}

func (h *IndexHandler) Index(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	return h.PostModel.List(_ctx, 0, model)
}

func (h *IndexHandler) IndexPage(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	page, err := util.ParamInt32(_ctx, ctx, "page")
	if err != nil {
		return "", err
	}
	return h.PostModel.List(_ctx, int(page)-1, model)
}
