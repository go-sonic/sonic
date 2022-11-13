package service

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type TagService interface {
	ListAll(ctx context.Context, sort *param.Sort) ([]*entity.Tag, error)
	ListByIDs(ctx context.Context, tagIDs []int32) ([]*entity.Tag, error)
	GetByID(ctx context.Context, id int32) (*entity.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Tag, error)
	GetByName(ctx context.Context, name string) (*entity.Tag, error)
	ConvertToDTO(ctx context.Context, tag *entity.Tag) (*dto.Tag, error)
	ConvertToDTOs(ctx context.Context, tags []*entity.Tag) ([]*dto.Tag, error)
	Create(ctx context.Context, tagParam *param.Tag) (*entity.Tag, error)
	Update(ctx context.Context, id int32, tagParam *param.Tag) (*entity.Tag, error)
	Delete(ctx context.Context, id int32) error
	CountAllTag(ctx context.Context) (int64, error)
}
