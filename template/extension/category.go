package extension

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
)

type categoryExtension struct {
	Template            *template.Template
	CategoryService     service.CategoryService
	PostCategoryService service.PostCategoryService
}

func RegisterCategoryFunc(t *template.Template, categoryService service.CategoryService, postCategoryService service.PostCategoryService) {
	ce := &categoryExtension{
		Template:            t,
		CategoryService:     categoryService,
		PostCategoryService: postCategoryService,
	}
	ce.addListCategoryFunc()
	ce.addListCategoryAsTreeFunc()
	ce.addGetCategoryCountFunc()
	ce.addListCategoryByPostIDFunc()
}

func (ce *categoryExtension) addListCategoryFunc() {
	listCategory := func() ([]*dto.CategoryWithPostCount, error) {
		sort := param.Sort{
			Fields: []string{"priority,asc"},
		}
		return ce.CategoryService.ListCategoryWithPostCountDTO(context.Background(), &sort)
	}
	ce.Template.AddFunc("listCategory", listCategory)
}

func (ce *categoryExtension) addListCategoryAsTreeFunc() {
	listCategoryAsTree := func() ([]*vo.CategoryVO, error) {
		sort := param.Sort{
			Fields: []string{"priority,asc"},
		}
		return ce.CategoryService.ListAsTree(context.Background(), &sort, false)
	}
	ce.Template.AddFunc("listCategoryAsTree", listCategoryAsTree)
}

func (ce *categoryExtension) addListCategoryByPostIDFunc() {
	listCategoryByPostID := func(postID int) ([]*dto.CategoryDTO, error) {
		categories, err := ce.PostCategoryService.ListCategoryByPostID(context.Background(), int32(postID))
		if err != nil {
			return nil, err
		}
		return ce.CategoryService.ConvertToCategoryDTOs(context.Background(), categories)
	}
	ce.Template.AddFunc("listCategoryByPostID", listCategoryByPostID)
}

func (ce *categoryExtension) addGetCategoryCountFunc() {
	getCategoryCount := func() (int64, error) {
		return ce.CategoryService.Count(context.Background())
	}
	ce.Template.AddFunc("getCategoryCount", getCategoryCount)
}
