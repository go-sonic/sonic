package model

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/content/authentication"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/template"
)

func NewCategoryModel(optionService service.OptionService,
	postService service.PostService,
	themeService service.ThemeService,
	postCategoryService service.PostCategoryService,
	categoryService service.CategoryService,
	postTagService service.PostTagService,
	tagService service.TagService,
	postAssembler assembler.PostAssembler,
	metaService service.MetaService,
	categoryAuthentication *authentication.CategoryAuthentication,
) *CategoryModel {
	return &CategoryModel{
		OptionService:          optionService,
		PostService:            postService,
		PostAssembler:          postAssembler,
		ThemeService:           themeService,
		PostCategoryService:    postCategoryService,
		CategoryService:        categoryService,
		PostTagService:         postTagService,
		TagService:             tagService,
		MetaService:            metaService,
		CategoryAuthentication: categoryAuthentication,
	}
}

type CategoryModel struct {
	OptionService          service.OptionService
	PostService            service.PostService
	ThemeService           service.ThemeService
	PostCategoryService    service.PostCategoryService
	CategoryService        service.CategoryService
	PostTagService         service.PostTagService
	TagService             service.TagService
	MetaService            service.MetaService
	PostAssembler          assembler.PostAssembler
	CategoryAuthentication *authentication.CategoryAuthentication
}

func (c *CategoryModel) ListCategories(ctx context.Context, model template.Model) (string, error) {

	seoKeyWords := c.OptionService.GetOrByDefault(ctx, property.SeoKeywords)
	seoDescription := c.OptionService.GetOrByDefault(ctx, property.SeoDescription)

	model["is_categories"] = true
	model["meta_keywords"] = seoKeyWords
	model["meta_description"] = seoDescription
	return c.ThemeService.Render(ctx, "categories")
}

func (c *CategoryModel) CategoryDetail(ctx context.Context, model template.Model, slug string, page int, token string) (string, error) {
	category, err := c.CategoryService.GetBySlug(ctx, slug)
	if err != nil {
		return "", err
	}

	if category.Type == consts.CategoryTypeIntimate {
		if isAuthenticated, err := c.CategoryAuthentication.IsAuthenticated(ctx, token, category.ID); err != nil || !isAuthenticated {
			model["slug"] = category.Slug
			model["type"] = consts.EncryptTypeCategory.Name()
			if exist, err := c.ThemeService.TemplateExist(ctx, "post_password.tmpl"); err == nil && exist {
				return c.ThemeService.Render(ctx, "post_password")
			}
			return "common/template/post_password", nil
		}
	}
	pageSize := c.OptionService.GetOrByDefault(ctx, property.ArchivePageSize).(int)
	sort := c.OptionService.GetPostSort(ctx)
	postQuery := param.PostQuery{
		Page: param.Page{
			PageNum:  page,
			PageSize: pageSize,
		},
		Sort:       &sort,
		Statuses:   []*consts.PostStatus{consts.PostStatusPublished.Ptr()},
		CategoryID: &category.ID,
	}
	if category.Password != "" {
		postQuery.Statuses = append(postQuery.Statuses, consts.PostStatusIntimate.Ptr())
	}
	posts, totalPage, err := c.PostService.Page(ctx, postQuery)
	if err != nil {
		return "", err
	}
	postVOs, err := c.PostAssembler.ConvertToListVO(ctx, posts)
	if err != nil {
		return "", err
	}
	postPage := dto.NewPage(postVOs, totalPage, param.Page{
		PageNum:  page,
		PageSize: pageSize,
	})
	categoryDTO, err := c.CategoryService.ConvertToCategoryDTO(ctx, category)
	if err != nil {
		return "", err
	}
	if categoryDTO.Description != "" {
		model["meta_description"] = categoryDTO.Description
	} else {
		model["meta_description"] = c.OptionService.GetOrByDefault(ctx, property.SeoDescription)
	}
	model["is_category"] = true
	model["posts"] = postPage
	model["category"] = categoryDTO
	model["meta_keywords"] = c.OptionService.GetOrByDefault(ctx, property.SeoKeywords)

	return c.ThemeService.Render(ctx, "category")
}
