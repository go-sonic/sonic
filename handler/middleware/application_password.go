package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/service"
)

var basicAuthRegexp = regexp.MustCompile(`^Basic [a-z\\d/+]*={0,2}`)

type ApplicationPasswordMiddleware struct {
	PasswordService service.ApplicationPasswordService
	UserService     service.UserService
}

func NewApplicationPasswordMiddleware(passwordService service.ApplicationPasswordService, userService service.UserService) *ApplicationPasswordMiddleware {
	m := &ApplicationPasswordMiddleware{
		PasswordService: passwordService,
		UserService:     userService,
	}
	return m
}

func (a *ApplicationPasswordMiddleware) Get() error {
	return nil
}

func (a *ApplicationPasswordMiddleware) GetWrapHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if len(header) == 0 {
			abortUnauthorized(ctx)
			return
		}

		match := verifyHeader(header)
		if !match {
			abortUnauthorized(ctx)
			return
		}

		bytes, err := base64.StdEncoding.DecodeString(header[6:])
		if err != nil {
			abortUnauthorized(ctx)
			return
		}

		userPass := string(bytes)

		if !strings.Contains(userPass, ":") {
			abortUnauthorized(ctx)
			return
		}

		splits := strings.SplitN(userPass, ":", 2)

		userEntity, err := a.UserService.GetByUsername(ctx, splits[0])
		if err != nil {
			abortUnauthorized(ctx)
			return
		}

		pwdEntity, err := a.PasswordService.Verify(ctx, userEntity.ID, splits[1])
		if err != nil || pwdEntity == nil {
			abortUnauthorized(ctx)
			return
		}

		err = a.PasswordService.Update(ctx, pwdEntity.ID, ctx.ClientIP())
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, &dto.BaseDTO{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("Update application password entity error, err=%s", err),
			})
			return
		}
		ctx.Set(consts.AuthorizedUser, userEntity)
	}
}

func abortUnauthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, &dto.BaseDTO{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
	})
}

func verifyHeader(header string) bool {
	return basicAuthRegexp.MatchString(header)
}
