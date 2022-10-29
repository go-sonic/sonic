package assembler

import (
	"context"

	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
)

type SheetCommentAssembler interface {
	BaseCommentAssembler
	ConvertToWithSheet(ctx context.Context, comments []*entity.Comment) ([]*vo.SheetCommentWithSheet, error)
}

func NewSheetCommentAssembler(
	optionService service.OptionService,
	baseCommentService service.BaseCommentService,
	baseCommentAssembler BaseCommentAssembler,
	sheetService service.SheetService,
	sheetAssembler SheetAssembler,
) SheetCommentAssembler {
	return &sheetCommentAssembler{
		OptionService:        optionService,
		BaseCommentService:   baseCommentService,
		BaseCommentAssembler: baseCommentAssembler,
		SheetAssembler:       sheetAssembler,
		SheetService:         sheetService,
	}
}

type sheetCommentAssembler struct {
	OptionService      service.OptionService
	BaseCommentService service.BaseCommentService
	BaseCommentAssembler
	SheetAssembler
	SheetService service.SheetService
}

func (p *sheetCommentAssembler) ConvertToWithSheet(ctx context.Context, comments []*entity.Comment) ([]*vo.SheetCommentWithSheet, error) {
	sheetIDs := make([]int32, 0, len(comments))
	for _, comment := range comments {
		sheetIDs = append(sheetIDs, comment.PostID)
	}
	sheets, err := p.SheetService.GetByPostIDs(ctx, sheetIDs)
	if err != nil {
		return nil, err
	}
	result := make([]*vo.SheetCommentWithSheet, 0, len(comments))
	for _, comment := range comments {
		commentDTO, err := p.BaseCommentAssembler.ConvertToDTO(ctx, comment)
		if err != nil {
			return nil, err
		}
		commentWithSheet := &vo.SheetCommentWithSheet{
			Comment: *commentDTO,
		}
		result = append(result, commentWithSheet)
		sheet, ok := sheets[comment.PostID]
		if ok {
			commentWithSheet.PostMinimal, err = p.SheetAssembler.ConvertToMinimalDTO(ctx, sheet)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
