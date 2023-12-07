package service

import (
	"context"

	"github.com/go-sonic/sonic/model/param"
)

type ScrapService interface {
	QueryMd5List(ctx context.Context) ([]string, error)
	Create(ctx context.Context, pageParam *param.ScrapPage) error
}
