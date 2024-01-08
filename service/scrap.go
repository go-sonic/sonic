package service

import (
	"context"
	"mime/multipart"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
)

type ScrapService interface {
	QueryMd5List(ctx context.Context) ([]string, error)
	Create(ctx context.Context, pageParam *param.ScrapPage, file *multipart.FileHeader) (*dto.ScrapPageDTO, error)
	Get(ctx context.Context, pageID int32) (*dto.ScrapPageDTO, error)
	Query(ctx context.Context, query *param.ScrapPageQuery) ([]*dto.ScrapPageDTO, int64, error)
	Update(ctx context.Context, pageId int32, pageParam *param.ScrapPage) (*dto.ScrapPageDTO, error)
}
