package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
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

func (c *CategoryHandler) GetCategoryByID(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "categoryID")
	if err != nil {
		return nil, err
	}
	category, err := c.CategoryService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return c.CategoryService.ConvertToCategoryDTO(ctx, category)
}

func (c *CategoryHandler) ListAllCategory(ctx *gin.Context) (interface{}, error) {
	categoryQuery := struct {
		*param.Sort
		More *bool `json:"more" form:"more"`
	}{}

	err := ctx.ShouldBindQuery(&categoryQuery)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if categoryQuery.Sort == nil || len(categoryQuery.Sort.Fields) == 0 {
		categoryQuery.Sort = &param.Sort{Fields: []string{"priority,asc"}}
	}
	if categoryQuery.More != nil && *categoryQuery.More {
		return c.CategoryService.ListCategoryWithPostCountDTO(ctx, categoryQuery.Sort)
	}
	categories, err := c.CategoryService.ListAll(ctx, categoryQuery.Sort)
	if err != nil {
		return nil, err
	}
	return c.CategoryService.ConvertToCategoryDTOs(ctx, categories)
}

func (c *CategoryHandler) ListAsTree(ctx *gin.Context) (interface{}, error) {
	var sort param.Sort
	err := ctx.ShouldBindQuery(&sort)
	if err != nil {
		return nil, err
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "priority,asc")
	}
	return c.CategoryService.ListAsTree(ctx, &sort, false)
}

func (c *CategoryHandler) CreateCategory(ctx *gin.Context) (interface{}, error) {
	var categoryParam param.Category
	err := ctx.ShouldBindJSON(&categoryParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	category, err := c.CategoryService.Create(ctx, &categoryParam)
	if err != nil {
		return nil, err
	}
	return c.CategoryService.ConvertToCategoryDTO(ctx, category)
}

func (c *CategoryHandler) UpdateCategory(ctx *gin.Context) (interface{}, error) {
	var categoryParam param.Category
	err := ctx.ShouldBindJSON(&categoryParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	categoryID, err := util.ParamInt32(ctx, "categoryID")
	if err != nil {
		return nil, err
	}
	categoryParam.ID = categoryID
	category, err := c.CategoryService.Update(ctx, &categoryParam)
	if err != nil {
		return nil, err
	}
	return c.CategoryService.ConvertToCategoryDTO(ctx, category)
}

func (c *CategoryHandler) UpdateCategoryBatch(ctx *gin.Context) (interface{}, error) {
	categoryParams := make([]*param.Category, 0)
	err := ctx.ShouldBindJSON(&categoryParams)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	categories, err := c.CategoryService.UpdateBatch(ctx, categoryParams)
	if err != nil {
		return nil, err
	}
	return c.CategoryService.ConvertToCategoryDTOs(ctx, categories)
}

func (c *CategoryHandler) DeleteCategory(ctx *gin.Context) (interface{}, error) {
	categoryID, err := util.ParamInt32(ctx, "categoryID")
	if err != nil {
		return nil, err
	}
	return nil, c.CategoryService.Delete(ctx, categoryID)
}
