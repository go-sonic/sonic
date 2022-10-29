package authentication

import (
	"context"

	"github.com/golang-jwt/jwt"
)

type ContentAuthentication interface {
	Authenticate(ctx context.Context, token string, id int32, password string) (string, error)
	IsAuthenticated(ctx context.Context, token string, id int32) (bool, error)
}

type customClaims struct {
	CategoryIDs []int32 `json:"category_ids"`
	PostIDs     []int32 `json:"post_ids"`
	jwt.StandardClaims
}
