package impl

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type journalServiceImpl struct {
	JournalCommentService service.JournalCommentService
}

func (*journalServiceImpl) Page(ctx context.Context, page param.Page, sort *param.Sort) ([]*entity.Journal, int64, error) {
	if page.PageNum < 0 || page.PageSize <= 0 || page.PageSize > 100 {
		return nil, 0, xerr.BadParam.New("").WithStatus(xerr.StatusBadRequest).WithMsg("Paging parameter error")
	}
	journalDAL := dal.Use(dal.GetDBByCtx(ctx)).Journal
	journalDO := journalDAL.WithContext(ctx)
	err := BuildSort(sort, &journalDAL, &journalDO)
	if err != nil {
		return nil, 0, err
	}
	journals, totalCount, err := journalDO.FindByPage(page.PageSize*page.PageNum, page.PageSize)
	if err != nil {
		return nil, 0, WrapDBErr(err)
	}
	return journals, totalCount, nil
}

func NewJournalService(journalCommentService service.JournalCommentService) service.JournalService {
	return &journalServiceImpl{
		JournalCommentService: journalCommentService,
	}
}

func (j *journalServiceImpl) ConvertToDTO(journal *entity.Journal) *dto.Journal {
	return &dto.Journal{
		ID:            journal.ID,
		SourceContent: journal.SourceContent,
		Content:       journal.Content,
		Likes:         journal.Likes,
		CreateTime:    journal.CreateTime.UnixMilli(),
		JournalType:   journal.Type,
	}
}

func (j *journalServiceImpl) ConvertToWithCommentDTOList(ctx context.Context, journals []*entity.Journal) ([]*dto.JournalWithComment, error) {
	journalIDs := make([]int32, 0, len(journals))
	for _, journal := range journals {
		journalIDs = append(journalIDs, journal.ID)
	}
	commentCountMap, err := j.JournalCommentService.CountByStatusAndJournalID(ctx, consts.CommentStatusPublished, journalIDs)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.JournalWithComment, 0, len(journals))
	for _, journal := range journals {
		journalWithCommentCount := &dto.JournalWithComment{
			Journal: *j.ConvertToDTO(journal),
		}
		if commentCount, ok := commentCountMap[journal.ID]; ok {
			journalWithCommentCount.CommentCount = commentCount
		}
		result = append(result, journalWithCommentCount)
	}
	return result, nil
}

func (j *journalServiceImpl) ListJournal(ctx context.Context, journalQuery param.JournalQuery) ([]*entity.Journal, int64, error) {
	if journalQuery.PageNum < 0 || journalQuery.PageSize <= 0 || journalQuery.PageSize > 100 {
		return nil, 0, xerr.BadParam.New("").WithStatus(xerr.StatusBadRequest).WithMsg("Paging parameter error")
	}
	journalDAL := dal.Use(dal.GetDBByCtx(ctx)).Journal
	journalDO := journalDAL.WithContext(ctx)
	err := BuildSort(journalQuery.Sort, &journalDAL, &journalDO)
	if err != nil {
		return nil, 0, err
	}

	if journalQuery.Keyword != nil {
		journalDO.Where(journalDAL.Content.Like(*journalQuery.Keyword))
	}
	if journalQuery.JournalType != nil {
		journalDO.Where(journalDAL.Type.Like(*journalQuery.JournalType))
	}

	posts, totalCount, err := journalDO.FindByPage(journalQuery.PageNum*journalQuery.PageSize, journalQuery.PageSize)
	if err != nil {
		return nil, 0, WrapDBErr(err)
	}
	return posts, totalCount, nil
}

func (j *journalServiceImpl) Create(ctx context.Context, journalParam *param.Journal) (*entity.Journal, error) {
	journal := &entity.Journal{
		Type:          journalParam.Type,
		SourceContent: journalParam.SourceContent,
		Content:       journalParam.Content,
	}
	journalDAL := dal.Use(dal.GetDBByCtx(ctx)).Journal
	err := journalDAL.WithContext(ctx).Create(journal)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return journal, nil
}

func (j *journalServiceImpl) Update(ctx context.Context, journalID int32, journalParam *param.Journal) (*entity.Journal, error) {
	journalDAL := dal.Use(dal.GetDBByCtx(ctx)).Journal
	journal, err := journalDAL.WithContext(ctx).Where(journalDAL.ID.Eq(journalID)).Take()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	journal.SourceContent = journalParam.SourceContent
	journal.Content = journalParam.Content
	if err != nil {
		return nil, err
	}
	updateResult, err := journalDAL.WithContext(ctx).Where(journalDAL.ID.Eq(journalID)).UpdateSimple(journalDAL.SourceContent.Value(journal.SourceContent), journalDAL.Content.Value(journal.Content), journalDAL.Type.Value(journalParam.Type))
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return nil, xerr.NoType.New("").WithMsg("update journal failed")
	}
	return journal, nil
}

func (j *journalServiceImpl) Delete(ctx context.Context, journalID int32) error {
	journalDAL := dal.Use(dal.GetDBByCtx(ctx)).Journal
	deleteResult, err := journalDAL.WithContext(ctx).Where(journalDAL.ID.Eq(journalID)).Delete()
	if err != nil {
		return WrapDBErr(err)
	}
	if deleteResult.RowsAffected != 1 {
		return xerr.NoType.New("journalID=%v", journalID).WithMsg("delete failed")
	}
	return nil
}

func (j *journalServiceImpl) GetByJournalIDs(ctx context.Context, journalIDs []int32) (map[int32]*entity.Journal, error) {
	journalDAL := dal.Use(dal.GetDBByCtx(ctx)).Journal
	journals, err := journalDAL.WithContext(ctx).Where(journalDAL.ID.In(journalIDs...)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	result := make(map[int32]*entity.Journal)
	for _, journal := range journals {
		result[journal.ID] = journal
	}
	return result, nil
}

func (j *journalServiceImpl) Count(ctx context.Context) (int64, error) {
	journalDAL := dal.Use(dal.GetDBByCtx(ctx)).Journal
	count, err := journalDAL.WithContext(ctx).Count()
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return count, nil
}

func (*journalServiceImpl) IncreaseLike(ctx context.Context, journalID int32) error {
	journalDAL := dal.Use(dal.GetDBByCtx(ctx)).Journal
	updateResult, err := journalDAL.WithContext(ctx).Where(journalDAL.ID.Eq(journalID)).UpdateSimple(journalDAL.Likes.Add(1))
	if err != nil {
		return WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return xerr.NoType.New("increase journal like failed id=%v", journalID).WithMsg("increase like failed").WithStatus(xerr.StatusInternalServerError)
	}
	return nil
}
