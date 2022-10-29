package service

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type JournalService interface {
	Page(ctx context.Context, page param.Page, sort *param.Sort) ([]*entity.Journal, int64, error)
	ListJournal(ctx context.Context, journalQuery param.JournalQuery) ([]*entity.Journal, int64, error)
	ConvertToDTO(journal *entity.Journal) *dto.Journal
	ConvertToWithCommentDTOList(ctx context.Context, journals []*entity.Journal) ([]*dto.JournalWithComment, error)
	Create(ctx context.Context, journalParam *param.Journal) (*entity.Journal, error)
	Update(ctx context.Context, journalID int32, journalParam *param.Journal) (*entity.Journal, error)
	Delete(ctx context.Context, journalID int32) error
	GetByJournalIDs(ctx context.Context, journalIDs []int32) (map[int32]*entity.Journal, error)
	Count(ctx context.Context) (int64, error)
	IncreaseLike(ctx context.Context, journalID int32) error
}
