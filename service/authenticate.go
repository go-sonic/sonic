package service

import (
	"context"

	"github.com/go-sonic/sonic/model/entity"
)

type AuthenticateService interface {
	PostAuthenticate(ctx context.Context, post *entity.Post, password string) (bool, error)
	CategoryAuthenticate(ctx context.Context, categoryID int32, password string) (bool, error)
}
