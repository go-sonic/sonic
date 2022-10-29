package assembler

import (
	"context"

	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
)

type JournalCommentAssembler interface {
	BaseCommentAssembler
	ConvertToWithJournal(ctx context.Context, comments []*entity.Comment) ([]*vo.JournalCommentWithJournal, error)
}

func NewJournalCommentAssembler(
	optionService service.OptionService,
	baseCommentService service.BaseCommentService,
	baseCommentAssembler BaseCommentAssembler,
	journalService service.JournalService,
) JournalCommentAssembler {
	return &journalCommentAssembler{
		OptionService:        optionService,
		BaseCommentService:   baseCommentService,
		BaseCommentAssembler: baseCommentAssembler,
		JournalService:       journalService,
	}
}

type journalCommentAssembler struct {
	OptionService      service.OptionService
	BaseCommentService service.BaseCommentService
	BaseCommentAssembler
	JournalService service.JournalService
}

func (j *journalCommentAssembler) ConvertToWithJournal(ctx context.Context, comments []*entity.Comment) ([]*vo.JournalCommentWithJournal, error) {
	journalIDs := make([]int32, 0, len(comments))
	for _, comment := range comments {
		journalIDs = append(journalIDs, comment.PostID)
	}
	journals, err := j.JournalService.GetByJournalIDs(ctx, journalIDs)
	if err != nil {
		return nil, err
	}
	result := make([]*vo.JournalCommentWithJournal, 0, len(comments))
	for _, comment := range comments {
		commentDTO, err := j.BaseCommentAssembler.ConvertToDTO(ctx, comment)
		if err != nil {
			return nil, err
		}
		commentWithJournal := &vo.JournalCommentWithJournal{
			Comment: *commentDTO,
		}
		result = append(result, commentWithJournal)
		journal, ok := journals[comment.PostID]
		if ok {
			commentWithJournal.Journal = j.JournalService.ConvertToDTO(journal)
		}
	}
	return result, nil
}
