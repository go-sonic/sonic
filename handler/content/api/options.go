package api

import (
	"github.com/gin-gonic/gin"

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

func (o *OptionHandler) Comment(ctx *gin.Context) (interface{}, error) {
	result := make(map[string]interface{})

	result[property.CommentGravatarSource.KeyValue] = o.OptionService.GetOrByDefault(ctx, property.CommentGravatarSource)
	result[property.CommentGravatarDefault.KeyValue] = o.OptionService.GetOrByDefault(ctx, property.CommentGravatarDefault)
	result[property.CommentContentPlaceholder.KeyValue] = o.OptionService.GetOrByDefault(ctx, property.CommentContentPlaceholder)
	return result, nil
}
