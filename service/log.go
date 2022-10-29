package service

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type LogService interface {
	PageLog(ctx context.Context, page param.Page, sort *param.Sort) ([]*entity.Log, int64, error)
	ConvertToDTO(log *entity.Log) *dto.Log
	Clear(ctx context.Context) error
}
