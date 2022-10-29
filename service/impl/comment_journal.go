package impl

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
)

type journalCommentServiceImpl struct {
	service.BaseCommentService
}

func NewJournalCommentService(baseCommentService service.BaseCommentService) service.JournalCommentService {
	return &journalCommentServiceImpl{
		BaseCommentService: baseCommentService,
	}
}

func (j *journalCommentServiceImpl) CountByStatusAndJournalID(ctx context.Context, status consts.CommentStatus, journalIDs []int32) (map[int32]int64, error) {
	return j.CountByStatusAndPostIDs(ctx, status, journalIDs)
}

func (j *journalCommentServiceImpl) UpdateBy(ctx context.Context, commentID int64, commentParam *param.Comment) (*entity.Comment, error) {
	if commentID == 0 {
		return nil, nil
	}
	comment := j.BaseCommentService.ConvertParam(commentParam)
	comment.ID = commentID
	return j.BaseCommentService.Update(ctx, comment)
}

func (j *journalCommentServiceImpl) CountByStatus(ctx context.Context, status consts.CommentStatus) (int64, error) {
	commentDAL := dal.Use(dal.GetDBByCtx(ctx)).Comment
	count, err := commentDAL.WithContext(ctx).Where(commentDAL.Type.Eq(consts.CommentTypeJournal), commentDAL.Status.Eq(status)).Count()
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return count, nil
}
