package service

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/vo"
)

type LinkService interface {
	List(ctx context.Context, sort *param.Sort) ([]*entity.Link, error)
	GetByID(ctx context.Context, id int32) (*entity.Link, error)
	Create(ctx context.Context, linkParam *param.Link) (*entity.Link, error)
	Update(ctx context.Context, id int32, linkParam *param.Link) (*entity.Link, error)
	Delete(ctx context.Context, id int32) error
	ConvertToDTO(ctx context.Context, link *entity.Link) *dto.Link
	ConvertToDTOs(ctx context.Context, links []*entity.Link) []*dto.Link
	ConvertToLinkTeamVO(ctx context.Context, links []*entity.Link) []*vo.LinkTeamVO
	ListTeams(ctx context.Context) ([]string, error)
	Count(ctx context.Context) (int64, error)
}
