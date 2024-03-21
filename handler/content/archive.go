package content

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/cache"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/content/model"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type ArchiveHandler struct {
	OptionService       service.OptionService
	PostService         service.PostService
	PostCategoryService service.PostCategoryService
	CategoryService     service.CategoryService
	PostAssembler       assembler.PostAssembler
	PostModel           *model.PostModel
	Cache               cache.Cache
}

func NewArchiveHandler(
	optionService service.OptionService,
	postService service.PostService,
	categoryService service.CategoryService,
	postCategoryService service.PostCategoryService,
	postAssembler assembler.PostAssembler,
	postModel *model.PostModel,
	cache cache.Cache,
) *ArchiveHandler {
	return &ArchiveHandler{
		OptionService:       optionService,
		PostService:         postService,
		PostCategoryService: postCategoryService,
		CategoryService:     categoryService,
		PostAssembler:       postAssembler,
		PostModel:           postModel,
		Cache:               cache,
	}
}

func (a *ArchiveHandler) Archives(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	return a.PostModel.Archives(_ctx, 0, model)
}

func (a *ArchiveHandler) ArchivesPage(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	page, err := util.ParamInt32(_ctx, ctx, "page")
	if err != nil {
		return "", err
	}
	return a.PostModel.Archives(_ctx, int(page-1), model)
}

func (a *ArchiveHandler) ArchivesBySlug(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	slug, err := util.ParamString(_ctx, ctx, "slug")
	if err != nil {
		return "", err
	}

	postPermalinkType, err := a.OptionService.GetPostPermalinkType(_ctx)
	if err != nil {
		return "", err
	}
	var post *entity.Post
	if postPermalinkType == consts.PostPermalinkTypeDefault {
		post, err = a.PostService.GetBySlug(_ctx, slug)
		if err != nil {
			return "", err
		}
	} else if postPermalinkType == consts.PostPermalinkTypeID {
		postID, err := strconv.ParseInt(slug, 10, 32)
		if err != nil {
			return "", err
		}
		post, err = a.PostService.GetByPostID(_ctx, int32(postID))
		if err != nil {
			return "", err
		}
	}
	token := string(ctx.Cookie("authentication"))
	return a.PostModel.Content(_ctx, post, token, model)
}

// AdminArchivesBySlug It can only be used in the console  to preview articles
func (a *ArchiveHandler) AdminArchivesBySlug(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	slug, err := util.ParamString(_ctx, ctx, "slug")
	if err != nil {
		return "", err
	}
	token, err := util.MustGetQueryString(_ctx, ctx, "token")
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", nil
	}
	_, ok := a.Cache.Get(token)
	if !ok {
		return "", xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("token已过期或者不存在")
	}

	postPermalinkType, err := a.OptionService.GetPostPermalinkType(_ctx)
	if err != nil {
		return "", err
	}
	var post *entity.Post
	if postPermalinkType == consts.PostPermalinkTypeDefault {
		post, err = a.PostService.GetBySlug(_ctx, slug)
		if err != nil {
			return "", err
		}
	} else if postPermalinkType == consts.PostPermalinkTypeID {
		postID, err := strconv.ParseInt(slug, 10, 32)
		if err != nil {
			return "", err
		}
		post, err = a.PostService.GetByPostID(_ctx, int32(postID))
		if err != nil {
			return "", err
		}
	}
	return a.PostModel.AdminPreview(_ctx, post, model)
}
