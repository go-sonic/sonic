package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/cache"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type AuthMiddleware struct {
	OptionService       service.OptionService
	OneTimeTokenService service.OneTimeTokenService
	UserService         service.UserService
	Cache               cache.Cache
}

func NewAuthMiddleware(optionService service.OptionService, oneTimeTokenService service.OneTimeTokenService, cache cache.Cache, userService service.UserService) *AuthMiddleware {
	authMiddleware := &AuthMiddleware{
		OptionService:       optionService,
		OneTimeTokenService: oneTimeTokenService,
		Cache:               cache,
		UserService:         userService,
	}
	return authMiddleware
}

func (a *AuthMiddleware) GetWrapHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isInstalled, err := a.OptionService.GetOrByDefaultWithErr(ctx, property.IsInstalled, false)
		if err != nil {
			abortWithStatusJSON(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if !isInstalled.(bool) {
			abortWithStatusJSON(ctx, http.StatusBadRequest, "Blog is not initialized")
			return
		}

		oneTimeToken, ok := ctx.GetQuery(consts.OneTimeTokenQueryName)
		if ok {
			allowedURL, ok := a.OneTimeTokenService.Get(oneTimeToken)
			if !ok {
				abortWithStatusJSON(ctx, http.StatusBadRequest, "OneTimeToken is not exist or expired")
				return
			}
			currentURL := ctx.Request.URL.Path
			if currentURL != allowedURL {
				abortWithStatusJSON(ctx, http.StatusBadRequest, "The one-time token does not correspond the request uri")
				return
			}
			return
		}

		token := ctx.GetHeader(consts.AdminTokenHeaderName)
		if token == "" {
			abortWithStatusJSON(ctx, http.StatusUnauthorized, "未登录，请登录后访问")
			return
		}
		userID, ok := a.Cache.Get(cache.BuildTokenAccessKey(token))

		if !ok || userID == nil {
			abortWithStatusJSON(ctx, http.StatusUnauthorized, "Token 已过期或不存在")
			return
		}

		user, err := a.UserService.GetByID(ctx, userID.(int32))
		if xerr.GetType(err) == xerr.NoRecord {
			_ = ctx.Error(err)
			abortWithStatusJSON(ctx, http.StatusUnauthorized, "用户不存在")
			return
		}
		if err != nil {
			_ = ctx.Error(err)
			abortWithStatusJSON(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		ctx.Set(consts.AuthorizedUser, user)
	}
}

func abortWithStatusJSON(ctx *gin.Context, status int, message string) {
	ctx.AbortWithStatusJSON(status, &dto.BaseDTO{
		Status:  status,
		Message: message,
	})
}
