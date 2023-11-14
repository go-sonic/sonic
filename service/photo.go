package service

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type PhotoService interface {
	List(ctx context.Context, sort *param.Sort) ([]*entity.Photo, error)
	Page(ctx context.Context, page param.Page, sort *param.Sort) ([]*entity.Photo, int64, error)
	GetByID(ctx context.Context, id int32) (*entity.Photo, error)
	Create(ctx context.Context, photoParam *param.Photo) (*entity.Photo, error)
	CreateBatch(ctx context.Context, photosParam []*param.Photo) ([]*entity.Photo, error)
	Update(ctx context.Context, id int32, photoParam *param.Photo) (*entity.Photo, error)
	Delete(ctx context.Context, id int32) error
	ConvertToDTO(ctx context.Context, photo *entity.Photo) *dto.Photo
	ConvertToDTOs(ctx context.Context, photos []*entity.Photo) []*dto.Photo
	ListTeams(ctx context.Context) ([]string, error)
	ListByTeam(ctx context.Context, team string, sort *param.Sort) ([]*entity.Photo, error)
	IncreaseLike(ctx context.Context, id int32) error
	GetPhotoCount(ctx context.Context) (int64, error)
}
