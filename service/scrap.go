package service

import "context"

type ScrapService interface {
	QueryMd5List(ctx context.Context) ([]string, error)
}
