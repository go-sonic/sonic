package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

func (i *InstallRedirectMiddleware) InstallRedirect() gin.HandlerFunc {
	skipPath := map[string]struct{}{
		"/api/admin/installations":  {},
		"/api/admin/is_installed":   {},
		"/api/admin/login/precheck": {},
	}
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if _, ok := skipPath[path]; ok {
			return
		}
		isInstall, err := i.optionService.GetOrByDefaultWithErr(ctx, property.IsInstalled, false)
		if err != nil {
			abortWithStatusJSON(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if !isInstall.(bool) {
			ctx.Redirect(http.StatusFound, "/admin/#install")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
