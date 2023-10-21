package model

import (
	"context"

	"github.com/go-sonic/sonic/consts"
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
	journalType := consts.JournalTypePublic
	journals, total, err := p.JournalService.ListJournal(ctx, param.JournalQuery{
		Page: param.Page{PageNum: page, PageSize: pageSize},
		Sort: &param.Sort{
			Fields: []string{"createTime,desc"},
		},
		Keyword:     nil,
		JournalType: &journalType,
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
