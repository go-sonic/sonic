package model

import (
	"context"
	"strings"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/content/authentication"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util/xerr"
)

func NewSheetModel(optionService service.OptionService,
	themeService service.ThemeService,
	postTagService service.PostTagService,
	tagService service.TagService,
	metaService service.MetaService,
	sheetAssembler assembler.SheetAssembler,
	sheetService service.SheetService,
	postAuthentication *authentication.PostAuthentication,
) *SheetModel {
	return &SheetModel{
		OptionService:      optionService,
		ThemeService:       themeService,
		PostTagService:     postTagService,
		TagService:         tagService,
		MetaService:        metaService,
		SheetAssembler:     sheetAssembler,
		SheetService:       sheetService,
		PostAuthentication: postAuthentication,
	}
}

type SheetModel struct {
	SheetService       service.SheetService
	OptionService      service.OptionService
	ThemeService       service.ThemeService
	PostTagService     service.PostTagService
	TagService         service.TagService
	MetaService        service.MetaService
	SheetAssembler     assembler.SheetAssembler
	PostAuthentication *authentication.PostAuthentication
}

func (s *SheetModel) Content(ctx context.Context, sheet *entity.Post, token string, model template.Model) (string, error) {
	if sheet == nil {
		return "", xerr.WithStatus(nil, int(xerr.StatusBadRequest)).WithMsg("查询不到文章信息")
	}
	if sheet.Status == consts.PostStatusRecycle || sheet.Status == consts.PostStatusDraft {
		return "", xerr.WithStatus(nil, xerr.StatusNotFound).WithMsg("查询不到文章信息")
	} else if sheet.Status == consts.PostStatusIntimate {
		if isAuthenticated, err := s.PostAuthentication.IsAuthenticated(ctx, token, sheet.ID); err != nil || !isAuthenticated {
			model["slug"] = sheet.Slug
			model["type"] = consts.EncryptTypePost.Name()
			if exist, err := s.ThemeService.TemplateExist(ctx, "post_password.tmpl"); err == nil && exist {
				return s.ThemeService.Render(ctx, "post_password")
			}
			return "common/template/post_password", nil
		}
	}

	sheetVO, err := s.SheetAssembler.ConvertToDetailVO(ctx, sheet)
	if err != nil {
		return "", err
	}
	model["target"] = sheetVO
	model["type"] = "sheet"
	model["post"] = sheetVO
	model["sheet"] = sheetVO
	model["is_sheet"] = true

	metas, err := s.MetaService.GetPostMeta(ctx, sheet.ID)
	if err != nil {
		return "", err
	}
	model["metas"] = s.MetaService.ConvertToMetaDTOs(metas)

	tags, err := s.PostTagService.ListTagByPostID(ctx, sheet.ID)
	if err != nil {
		return "", err
	}
	model["tags"], _ = s.TagService.ConvertToDTOs(ctx, tags)

	if sheet.MetaDescription != "" {
		model["meta_description"] = sheet.MetaDescription
	} else {
		model["meta_description"] = sheet.Summary
	}
	if sheet.MetaKeywords != "" {
		model["meta_keywords"] = sheet.MetaKeywords
	} else if len(tags) > 0 {
		metaKeywords := strings.Builder{}
		metaKeywords.Write([]byte(tags[0].Name))
		for _, tag := range tags[1:] {
			metaKeywords.Write([]byte(","))
			metaKeywords.Write([]byte(tag.Name))
		}
		model["meta_keywords"] = metaKeywords.String()
	}

	s.SheetService.IncreaseVisit(ctx, sheet.ID)

	return s.ThemeService.Render(ctx, "sheet")
}
