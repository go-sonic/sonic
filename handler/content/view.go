package content

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/binding"
	"github.com/go-sonic/sonic/handler/content/authentication"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type ViewHandler struct {
	OptionService          service.OptionService
	UserService            service.UserService
	CategoryService        service.CategoryService
	PostService            service.PostService
	ThemeService           service.ThemeService
	CategoryAuthentication *authentication.CategoryAuthentication
	PostAuthentication     *authentication.PostAuthentication
}

func NewViewHandler(
	optionService service.OptionService,
	userService service.UserService,
	categoryService service.CategoryService,
	postService service.PostService,
	themeService service.ThemeService,
	categoryAuthentication *authentication.CategoryAuthentication,
	postAuthentication *authentication.PostAuthentication,
) *ViewHandler {
	return &ViewHandler{
		OptionService:          optionService,
		UserService:            userService,
		CategoryService:        categoryService,
		PostService:            postService,
		ThemeService:           themeService,
		CategoryAuthentication: categoryAuthentication,
		PostAuthentication:     postAuthentication,
	}
}

func (v *ViewHandler) Admin(ctx *gin.Context) (interface{}, error) {
	// TODO
	return nil, nil
}

func (v *ViewHandler) Version(ctx *gin.Context) (interface{}, error) {
	return consts.SonicVersion, nil
}

func (v *ViewHandler) Install(ctx *gin.Context) {
	isInstall := v.OptionService.GetOrByDefault(ctx, property.IsInstalled).(bool)
	if isInstall {
		return
	}
	adminURLPath, _ := v.OptionService.GetAdminURLPath(ctx)
	ctx.Redirect(http.StatusTemporaryRedirect, adminURLPath+"/#install")
}

func (v *ViewHandler) Logo(ctx *gin.Context) (interface{}, error) {
	logo := v.OptionService.GetOrByDefault(ctx, property.BlogLogo).(string)
	if logo != "" {
		ctx.Redirect(http.StatusTemporaryRedirect, logo)
	}
	return nil, nil
}

func (v *ViewHandler) Favicon(ctx *gin.Context) (interface{}, error) {
	favicon := v.OptionService.GetOrByDefault(ctx, property.BlogFavicon).(string)
	if favicon != "" {
		ctx.Redirect(http.StatusTemporaryRedirect, favicon)
	}
	return nil, nil
}

func (v *ViewHandler) Authenticate(ctx *gin.Context, model template.Model) (string, error) {
	contentType, err := util.ParamString(ctx, "type")
	if err != nil {
		return v.authenticateErr(ctx, model, contentType, "", err)
	}
	slug, err := util.ParamString(ctx, "slug")
	if err != nil {
		return v.authenticateErr(ctx, model, contentType, slug, err)
	}

	var authenticationParam param.Authentication
	err = ctx.ShouldBindWith(&authenticationParam, binding.CustomFormBinding)
	if err != nil {
		return v.authenticateErr(ctx, model, "post", slug, err)
	}
	if authenticationParam.Password == "" {
		return v.authenticateErr(ctx, model, "post", slug, xerr.WithMsg(nil, "密码为空"))
	}

	token, _ := ctx.Cookie("authentication")

	switch contentType {
	case consts.EncryptTypeCategory.Name():
		token, err = v.authenticateCategory(ctx, slug, authenticationParam.Password, token)
	case consts.EncryptTypePost.Name():
		token, err = v.authenticatePost(ctx, slug, authenticationParam.Password, token)
	default:
		return v.authenticateErr(ctx, model, "post", slug, xerr.WithStatus(nil, xerr.StatusBadRequest))
	}
	if err != nil {
		return v.authenticateErr(ctx, model, contentType, slug, err)
	}
	ctx.SetCookie("authentication", token, 1800, "/", "", false, true)
	return "", nil
}

func (v *ViewHandler) authenticateCategory(ctx *gin.Context, slug, password, token string) (string, error) {
	category, err := v.CategoryService.GetBySlug(ctx, slug)
	if err != nil {
		return "", err
	}
	categoryDTO, err := v.CategoryService.ConvertToCategoryDTO(ctx, category)
	if err != nil {
		return "", err
	}

	token, err = v.CategoryAuthentication.Authenticate(ctx, token, category.ID, password)
	if err != nil {
		return "", err
	}

	ctx.Redirect(http.StatusFound, categoryDTO.FullPath)
	return token, nil
}

func (v *ViewHandler) authenticatePost(ctx *gin.Context, slug, password, token string) (string, error) {
	post, err := v.PostService.GetBySlug(ctx, slug)
	if err != nil {
		return "", err
	}
	fullPath, err := v.PostService.BuildFullPath(ctx, post)
	if err != nil {
		return "", err
	}
	token, err = v.PostAuthentication.Authenticate(ctx, token, post.ID, password)
	if err != nil {
		return "", err
	}

	ctx.Redirect(http.StatusFound, fullPath)
	return token, nil
}

func (v *ViewHandler) authenticateErr(ctx *gin.Context, model template.Model, aType string, slug string, err error) (string, error) {
	model["type"] = aType
	model["slug"] = slug
	model["errorMsg"] = xerr.GetMessage(err)
	if exist, err := v.ThemeService.TemplateExist(ctx, "post_password.tmpl"); err == nil && exist {
		return v.ThemeService.Render(ctx, "post_password")
	}
	return "common/template/post_password", nil
}
