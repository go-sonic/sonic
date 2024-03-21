package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
)

type OptionHandler struct {
	OptionService service.OptionService
}

func NewOptionHandler(
	optionService service.OptionService,
) *OptionHandler {
	return &OptionHandler{
		OptionService: optionService,
	}
}

func (o *OptionHandler) Comment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	result := make(map[string]interface{})

	result[property.CommentGravatarSource.KeyValue] = o.OptionService.GetOrByDefault(_ctx, property.CommentGravatarSource)
	result[property.CommentGravatarDefault.KeyValue] = o.OptionService.GetOrByDefault(_ctx, property.CommentGravatarDefault)
	result[property.CommentContentPlaceholder.KeyValue] = o.OptionService.GetOrByDefault(_ctx, property.CommentContentPlaceholder)
	return result, nil
}
