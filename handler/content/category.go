package content

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/handler/content/model"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
)

type CategoryHandler struct {
	OptionService       service.OptionService
	PostService         service.PostService
	PostCategoryService service.PostCategoryService
	CategoryService     service.CategoryService
	PostAssembler       assembler.PostAssembler
	PostModel           *model.PostModel
	CategoryModel       *model.CategoryModel
}

func NewCategoryHandler(
	optionService service.OptionService,
	postService service.PostService,
	categoryService service.CategoryService,
	postCategoryService service.PostCategoryService,
	postAssembler assembler.PostAssembler,
	postModel *model.PostModel,
	categoryModel *model.CategoryModel,
) *CategoryHandler {
	return &CategoryHandler{
		OptionService:       optionService,
		PostService:         postService,
		PostCategoryService: postCategoryService,
		CategoryService:     categoryService,
		PostAssembler:       postAssembler,
		PostModel:           postModel,
		CategoryModel:       categoryModel,
	}
}

func (c *CategoryHandler) Categories(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	return c.CategoryModel.ListCategories(_ctx, model)
}

func (c *CategoryHandler) CategoryDetail(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	slug, err := util.ParamString(_ctx, ctx, "slug")
	if err != nil {
		return "", err
	}
	token := string(ctx.Cookie("authentication"))
	return c.CategoryModel.CategoryDetail(_ctx, model, slug, 0, token)
}

func (c *CategoryHandler) CategoryDetailPage(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	slug, err := util.ParamString(_ctx, ctx, "slug")
	if err != nil {
		return "", err
	}

	page, err := util.ParamInt32(_ctx, ctx, "page")
	if err != nil {
		return "", err
	}
	token := string(ctx.Cookie("authentication"))
	return c.CategoryModel.CategoryDetail(_ctx, model, slug, int(page-1), token)
}
