package service

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type SheetService interface {
	BasePostService
	Page(ctx context.Context, page param.Page, sort *param.Sort) ([]*entity.Post, int64, error)
	Create(ctx context.Context, sheetParam *param.Sheet) (*entity.Post, error)
	Update(ctx context.Context, sheetID int32, sheetParam *param.Sheet) (*entity.Post, error)
	Preview(ctx context.Context, sheetID int32) (string, error)
	CountVisit(ctx context.Context) (int64, error)
	CountLike(ctx context.Context) (int64, error)
	ListIndependentSheets(ctx context.Context) ([]*dto.IndependentSheet, error)
}
