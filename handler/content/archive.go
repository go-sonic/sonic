package content

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/content/model"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
)

type ArchiveHandler struct {
	OptionService       service.OptionService
	PostService         service.PostService
	PostCategoryService service.PostCategoryService
	CategoryService     service.CategoryService
	PostAssembler       assembler.PostAssembler
	PostModel           *model.PostModel
}

func NewArchiveHandler(
	optionService service.OptionService,
	postService service.PostService,
	categoryService service.CategoryService,
	postCategoryService service.PostCategoryService,
	postAssembler assembler.PostAssembler,
	postModel *model.PostModel,
) *ArchiveHandler {
	return &ArchiveHandler{
		OptionService:       optionService,
		PostService:         postService,
		PostCategoryService: postCategoryService,
		CategoryService:     categoryService,
		PostAssembler:       postAssembler,
		PostModel:           postModel,
	}
}

func (a *ArchiveHandler) Archives(ctx *gin.Context, model template.Model) (string, error) {
	return a.PostModel.Archives(ctx, 0, model)
}

func (a *ArchiveHandler) ArchivesPage(ctx *gin.Context, model template.Model) (string, error) {
	page, err := util.ParamInt32(ctx, "page")
	if err != nil {
		return "", err
	}
	return a.PostModel.Archives(ctx, int(page-1), model)
}

func (a *ArchiveHandler) ArchivesBySlug(ctx *gin.Context, model template.Model) (string, error) {
	slug, err := util.ParamString(ctx, "slug")
	if err != nil {
		return "", err
	}

	postPermalinkType, err := a.OptionService.GetPostPermalinkType(ctx)
	if err != nil {
		return "", err
	}
	var post *entity.Post
	if postPermalinkType == consts.PostPermalinkTypeDefault {
		post, err = a.PostService.GetBySlug(ctx, slug)
		if err != nil {
			return "", err
		}
	} else if postPermalinkType == consts.PostPermalinkTypeID {
		postID, err := strconv.ParseInt(slug, 10, 32)
		if err != nil {
			return "", err
		}
		post, err = a.PostService.GetByPostID(ctx, int32(postID))
		if err != nil {
			return "", err
		}
	}
	token, _ := ctx.Cookie("authentication")
	return a.PostModel.Content(ctx, post, token, model)
}
