package wp

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/model/dto/wp"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
)

type CategoryHandler struct {
	CategoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: categoryService,
	}
}

func (c *CategoryHandler) List(ctx *gin.Context) (interface{}, error) {
	sort := &param.Sort{
		Fields: []string{"name,desc"},
	}
	categoryEntities, err := c.CategoryService.ListAll(ctx, sort)
	if err != nil {
		return nil, err
	}

	categoryDTOList := make([]*wp.CategoryDTO, 0, len(categoryEntities))
	for _, categoryEntity := range categoryEntities {
		categoryDTOList = append(categoryDTOList, convertToCategoryDTO(categoryEntity))
	}

	return categoryDTOList, nil
}

func convertToCategoryDTO(categoryEntity *entity.Category) *wp.CategoryDTO {
	categoryDTO := &wp.CategoryDTO{
		ID:          categoryEntity.ID,
		Count:       0,
		Description: categoryEntity.Description,
		Link:        "",
		Name:        categoryEntity.Name,
		Slug:        categoryEntity.Slug,
		Taxonomy:    "",
		Parent:      categoryEntity.ParentID,
		Meta:        nil,
	}
	return categoryDTO
}
