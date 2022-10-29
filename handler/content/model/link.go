package model

import (
	"context"

	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
)

func NewLinkModel(
	optionService service.OptionService,
	themeService service.ThemeService,
	linkService service.LinkService,
) *LinkModel {
	return &LinkModel{
		OptionService: optionService,
		ThemeService:  themeService,
		LinkService:   linkService,
	}
}

type LinkModel struct {
	LinkService   service.LinkService
	OptionService service.OptionService
	ThemeService  service.ThemeService
}

func (l *LinkModel) Links(ctx context.Context, model template.Model) (string, error) {
	model["is_links"] = true
	model["meta_keywords"] = l.OptionService.GetOrByDefault(ctx, property.SeoKeywords)
	model["meta_description"] = l.OptionService.GetOrByDefault(ctx, property.SeoDescription)
	return l.ThemeService.Render(ctx, "links")
}
