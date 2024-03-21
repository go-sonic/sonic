package middleware

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
)

type InstallRedirectMiddleware struct {
	optionService service.OptionService
}

func NewInstallRedirectMiddleware(optionService service.OptionService) *InstallRedirectMiddleware {
	return &InstallRedirectMiddleware{
		optionService: optionService,
	}
}

func (i *InstallRedirectMiddleware) InstallRedirect() app.HandlerFunc {
	skipPath := map[string]struct{}{
		"/api/admin/installations":  {},
		"/api/admin/is_installed":   {},
		"/api/admin/login/precheck": {},
	}
	return func(_ctx context.Context, ctx *app.RequestContext) {
		path := string(ctx.URI().Path())
		if _, ok := skipPath[path]; ok {
			return
		}
		isInstall, err := i.optionService.GetOrByDefaultWithErr(_ctx, property.IsInstalled, false)
		if err != nil {
			abortWithStatusJSON(_ctx, ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if !isInstall.(bool) {
			ctx.Redirect(http.StatusFound, []byte("/admin/#install"))
			ctx.Abort()
			return
		}
		ctx.Next(_ctx)

	}
}
