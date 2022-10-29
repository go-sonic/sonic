package assembler

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
)

type SheetAssembler interface {
	BasePostAssembler
	ConvertToDetailVO(ctx context.Context, sheet *entity.Post) (*vo.SheetDetail, error)
	ConvertToListVO(ctx context.Context, sheets []*entity.Post) ([]*vo.SheetList, error)
}

func NewSheetAssembler(
	basePostService service.BasePostService,
	metaService service.MetaService,
	basePostAssembler BasePostAssembler,
	sheetCommentService service.SheetCommentService,
) SheetAssembler {
	return &sheetAssembler{
		BasePostAssembler:   basePostAssembler,
		MetaService:         metaService,
		SheetCommentService: sheetCommentService,
	}
}

type sheetAssembler struct {
	BasePostAssembler
	SheetCommentService service.SheetCommentService
	MetaService         service.MetaService
}

func (s *sheetAssembler) ConvertToDetailVO(ctx context.Context, sheet *entity.Post) (*vo.SheetDetail, error) {
	var sheetDetailVO vo.SheetDetail

	detailDTO, err := s.ConvertToDetailDTO(ctx, sheet)
	if err != nil {
		return nil, err
	}
	metas, err := s.MetaService.GetPostMeta(ctx, sheet.ID)
	if err != nil {
		return nil, err
	}
	metaIDs := make([]int64, 0, len(metas))
	metaDTOs := make([]*dto.Meta, 0, len(metas))
	for _, meta := range metas {
		metaIDs = append(metaIDs, meta.ID)
		metaDTOs = append(metaDTOs, s.MetaService.ConvertToMetaDTO(meta))
	}

	sheetDetailVO.PostDetail = *detailDTO
	sheetDetailVO.MetaIDs = metaIDs
	sheetDetailVO.Metas = metaDTOs
	return &sheetDetailVO, nil
}

func (s *sheetAssembler) ConvertToListVO(ctx context.Context, sheets []*entity.Post) ([]*vo.SheetList, error) {
	sheetListVOs := make([]*vo.SheetList, 0, len(sheets))

	for _, sheet := range sheets {
		var sheetListVO vo.SheetList
		postDTO, err := s.ConvertToSimpleDTO(ctx, sheet)
		if err != nil {
			return nil, err
		}
		commentCount, err := s.SheetCommentService.CountByPostID(ctx, sheet.ID)
		if err != nil {
			return nil, err
		}
		sheetListVO.CommentCount = commentCount
		sheetListVO.Post = *postDTO
		sheetListVOs = append(sheetListVOs, &sheetListVO)
	}
	return sheetListVOs, nil
}
