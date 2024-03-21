package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/content/authentication"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type CategoryHandler struct {
	PostService            service.PostService
	CategoryService        service.CategoryService
	CategoryAuthentication authentication.CategoryAuthentication
	PostAssembler          assembler.PostAssembler
}

func NewCategoryHandler(postService service.PostService, categoryService service.CategoryService, categoryAuthentication *authentication.CategoryAuthentication, postAssembler assembler.PostAssembler) *CategoryHandler {
	return &CategoryHandler{
		PostService:            postService,
		CategoryService:        categoryService,
		CategoryAuthentication: *categoryAuthentication,
		PostAssembler:          postAssembler,
	}
}

func (c *CategoryHandler) ListCategories(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	categoryQuery := struct {
		*param.Sort
		More *bool `json:"more" form:"more"`
	}{}

	err := ctx.BindAndValidate(&categoryQuery)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if categoryQuery.Sort == nil || len(categoryQuery.Sort.Fields) == 0 {
		categoryQuery.Sort = &param.Sort{Fields: []string{"updateTime,desc"}}
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

func (c *CategoryHandler) ListPosts(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	slug, err := util.ParamString(_ctx, ctx, "slug")
	if err != nil {
		return nil, err
	}
	category, err := c.CategoryService.GetBySlug(_ctx, slug)
	if err != nil {
		return nil, err
	}
	postQueryNoEnum := param.PostQueryNoEnum{}
	err = ctx.BindAndValidate(&postQueryNoEnum)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if postQueryNoEnum.Sort == nil {
		postQueryNoEnum.Sort = &param.Sort{Fields: []string{"topPriority,desc", "updateTime,desc"}}
	}
	password, _ := util.MustGetQueryString(_ctx, ctx, "password")

	if category.Type == consts.CategoryTypeIntimate {
		token := string(ctx.Cookie("authentication"))
		if authenticated, _ := c.CategoryAuthentication.IsAuthenticated(_ctx, token, category.ID); !authenticated {
			token, err := c.CategoryAuthentication.Authenticate(_ctx, token, category.ID, password)
			if err != nil {
				return nil, err
			}
			ctx.SetCookie("authentication", token, 1800, "/", "", 0, false, true)
		}
	}
	postQuery := param.AssertPostQuery(postQueryNoEnum)
	postQuery.WithPassword = util.BoolPtr(false)
	postQuery.Statuses = []*consts.PostStatus{consts.PostStatusPublished.Ptr(), consts.PostStatusIntimate.Ptr()}
	posts, totalCount, err := c.PostService.Page(_ctx, postQuery)
	if err != nil {
		return nil, err
	}
	postVOs, err := c.PostAssembler.ConvertToListVO(_ctx, posts)
	return dto.NewPage(postVOs, totalCount, postQuery.Page), err
}
