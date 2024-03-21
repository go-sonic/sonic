package content

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/consts"
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

func (v *ViewHandler) Admin(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	// TODO
	return nil, nil
}

func (v *ViewHandler) Version(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return consts.SonicVersion, nil
}

func (v *ViewHandler) Install(_ctx context.Context, ctx *app.RequestContext) {
	isInstall := v.OptionService.GetOrByDefault(_ctx, property.IsInstalled).(bool)
	if isInstall {
		return
	}
	adminURLPath, _ := v.OptionService.GetAdminURLPath(_ctx)
	ctx.Redirect(http.StatusTemporaryRedirect, []byte(adminURLPath+"/#install"))
}

func (v *ViewHandler) Logo(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	logo := v.OptionService.GetOrByDefault(_ctx, property.BlogLogo).(string)
	if logo != "" {
		ctx.Redirect(http.StatusTemporaryRedirect, []byte(logo))
	}
	return nil, nil
}

func (v *ViewHandler) Favicon(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	favicon := v.OptionService.GetOrByDefault(_ctx, property.BlogFavicon).(string)
	if favicon != "" {
		ctx.Redirect(http.StatusTemporaryRedirect, []byte(favicon))
	}
	return nil, nil
}

func (v *ViewHandler) Authenticate(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	contentType, err := util.ParamString(_ctx, ctx, "type")
	if err != nil {
		return v.authenticateErr(_ctx, ctx, model, contentType, "", err)
	}
	slug, err := util.ParamString(_ctx, ctx, "slug")
	if err != nil {
		return v.authenticateErr(_ctx, ctx, model, contentType, slug, err)
	}

	var authenticationParam param.Authentication
	err = ctx.BindAndValidate(&authenticationParam)
	if err != nil {
		return v.authenticateErr(_ctx, ctx, model, "post", slug, err)
	}
	if authenticationParam.Password == "" {
		return v.authenticateErr(_ctx, ctx, model, "post", slug, xerr.WithMsg(nil, "密码为空"))
	}

	token := string(ctx.Cookie("authentication"))

	switch contentType {
	case consts.EncryptTypeCategory.Name():
		token, err = v.authenticateCategory(_ctx, ctx, slug, authenticationParam.Password, token)
	case consts.EncryptTypePost.Name():
		token, err = v.authenticatePost(_ctx, ctx, slug, authenticationParam.Password, token)
	default:
		return v.authenticateErr(_ctx, ctx, model, "post", slug, xerr.WithStatus(nil, xerr.StatusBadRequest))
	}
	if err != nil {
		return v.authenticateErr(_ctx, ctx, model, contentType, slug, err)
	}
	ctx.SetCookie("authentication", token, 1800, "/", "", 0, false, true)
	return "", nil
}

func (v *ViewHandler) authenticateCategory(_ctx context.Context, ctx *app.RequestContext, slug, password, token string) (string, error) {
	category, err := v.CategoryService.GetBySlug(_ctx, slug)
	if err != nil {
		return "", err
	}
	categoryDTO, err := v.CategoryService.ConvertToCategoryDTO(_ctx, category)
	if err != nil {
		return "", err
	}

	token, err = v.CategoryAuthentication.Authenticate(_ctx, token, category.ID, password)
	if err != nil {
		return "", err
	}
	ctx.Redirect(http.StatusFound, []byte(categoryDTO.FullPath))
	return token, nil
}

func (v *ViewHandler) authenticatePost(_ctx context.Context, ctx *app.RequestContext, slug, password, token string) (string, error) {
	post, err := v.PostService.GetBySlug(_ctx, slug)
	if err != nil {
		return "", err
	}
	fullPath, err := v.PostService.BuildFullPath(_ctx, post)
	if err != nil {
		return "", err
	}
	token, err = v.PostAuthentication.Authenticate(_ctx, token, post.ID, password)
	if err != nil {
		return "", err
	}
	ctx.Redirect(http.StatusFound, []byte(fullPath))
	return token, nil
}

func (v *ViewHandler) authenticateErr(_ctx context.Context, ctx *app.RequestContext, model template.Model, aType string, slug string, err error) (string, error) {
	model["type"] = aType
	model["slug"] = slug
	model["errorMsg"] = xerr.GetMessage(err)
	if exist, err := v.ThemeService.TemplateExist(_ctx, "post_password.tmpl"); err == nil && exist {
		return v.ThemeService.Render(_ctx, "post_password")
	}
	return "common/template/post_password", nil
}
