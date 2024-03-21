package admin

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type CategoryHandler struct {
	CategoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: categoryService,
	}
}

func (c *CategoryHandler) GetCategoryByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "categoryID")
	if err != nil {
		return nil, err
	}
	category, err := c.CategoryService.GetByID(_ctx, id)
	if err != nil {
		return nil, err
	}
	return c.CategoryService.ConvertToCategoryDTO(_ctx, category)
}

func (c *CategoryHandler) ListAllCategory(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	categoryQuery := struct {
		*param.Sort
		More *bool `json:"more" form:"more"`
	}{}

	err := ctx.BindAndValidate(&categoryQuery)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if categoryQuery.Sort == nil || len(categoryQuery.Sort.Fields) == 0 {
		categoryQuery.Sort = &param.Sort{Fields: []string{"priority,asc"}}
	}
	if categoryQuery.More != nil && *categoryQuery.More {
		return c.CategoryService.ListCategoryWithPostCountDTO(_ctx, categoryQuery.Sort)
	}
	categories, err := c.CategoryService.ListAll(_ctx, categoryQuery.Sort)
	if err != nil {
		return nil, err
	}
	return c.CategoryService.ConvertToCategoryDTOs(_ctx, categories)
}

func (c *CategoryHandler) ListAsTree(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var sort param.Sort
	err := ctx.BindAndValidate(&sort)
	if err != nil {
		return nil, err
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "priority,asc")
	}
	return c.CategoryService.ListAsTree(_ctx, &sort, false)
}

func (c *CategoryHandler) CreateCategory(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var categoryParam param.Category
	err := ctx.BindAndValidate(&categoryParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	category, err := c.CategoryService.Create(_ctx, &categoryParam)
	if err != nil {
		return nil, err
	}
	return c.CategoryService.ConvertToCategoryDTO(_ctx, category)
}

func (c *CategoryHandler) UpdateCategory(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var categoryParam param.Category
	err := ctx.BindAndValidate(&categoryParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	categoryID, err := util.ParamInt32(_ctx, ctx, "categoryID")
	if err != nil {
		return nil, err
	}
	categoryParam.ID = categoryID
	category, err := c.CategoryService.Update(_ctx, &categoryParam)
	if err != nil {
		return nil, err
	}
	return c.CategoryService.ConvertToCategoryDTO(_ctx, category)
}

func (c *CategoryHandler) UpdateCategoryBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	categoryParams := make([]*param.Category, 0)
	err := ctx.BindAndValidate(&categoryParams)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	categories, err := c.CategoryService.UpdateBatch(_ctx, categoryParams)
	if err != nil {
		return nil, err
	}
	return c.CategoryService.ConvertToCategoryDTOs(_ctx, categories)
}

func (c *CategoryHandler) DeleteCategory(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	categoryID, err := util.ParamInt32(_ctx, ctx, "categoryID")
	if err != nil {
		return nil, err
	}
	return nil, c.CategoryService.Delete(_ctx, categoryID)
}
