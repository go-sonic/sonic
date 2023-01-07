package model

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
)

func NewJournalModel(optionService service.OptionService,
	themeService service.ThemeService,
	journalService service.JournalService,
) *JournalModel {
	return &JournalModel{
		OptionService:  optionService,
		ThemeService:   themeService,
		JournalService: journalService,
	}
}

type JournalModel struct {
	JournalService service.JournalService
	OptionService  service.OptionService
	ThemeService   service.ThemeService
}

func (p *JournalModel) Journals(ctx context.Context, model template.Model, page int) (string, error) {
	pageSize := p.OptionService.GetOrByDefault(ctx, property.JournalPageSize).(int)
	journals, total, err := p.JournalService.Page(ctx,
		param.Page{
			PageNum:  page,
			PageSize: pageSize,
		},
		&param.Sort{
			Fields: []string{"createTime,desc"},
		})
	if err != nil {
		return "", err
	}
	journalDTOs, err := p.JournalService.ConvertToWithCommentDTOList(ctx, journals)
	if err != nil {
		return "", err
	}
	journalPage := dto.NewPage(journalDTOs, total, param.Page{PageNum: page, PageSize: pageSize})
	model["is_journals"] = true
	model["journals"] = journalPage
	model["meta_keywords"] = p.OptionService.GetOrByDefault(ctx, property.SeoKeywords)
	model["meta_description"] = p.OptionService.GetOrByDefault(ctx, property.SeoDescription)
	return p.ThemeService.Render(ctx, "journals")
}
