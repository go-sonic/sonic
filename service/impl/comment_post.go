package impl

import (
	"context"

	"gorm.io/gorm"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type postCommentServiceImpl struct {
	service.BaseCommentService
}

func NewPostCommentService(baseCommentService service.BaseCommentService) service.PostCommentService {
	return &postCommentServiceImpl{
		BaseCommentService: baseCommentService,
	}
}

func (p postCommentServiceImpl) CreateBy(ctx context.Context, commentParam *param.Comment) (*entity.Comment, error) {
	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	post, err := postDAL.WithContext(ctx).Where(postDAL.ID.Eq(commentParam.PostID)).First()
	if err == gorm.ErrRecordNotFound {
		return nil, xerr.WithMsg(nil, "post not found").WithStatus(xerr.StatusBadRequest)
	}
	if err != nil {
		return nil, err
	}
	if post.DisallowComment {
		return nil, xerr.WithMsg(nil, "This post does not allow comments").WithStatus(xerr.StatusBadRequest)
	}
	return p.BaseCommentService.CreateBy(ctx, commentParam)
}

func (p postCommentServiceImpl) UpdateBy(ctx context.Context, commentID int64, commentParam *param.Comment) (*entity.Comment, error) {
	if commentID == 0 {
		return nil, nil
	}
	comment := p.ConvertParam(commentParam)
	comment.ID = commentID
	return p.Update(ctx, comment)
}

func (p postCommentServiceImpl) CountByStatus(ctx context.Context, status consts.CommentStatus) (int64, error) {
	postCommentDAL := dal.Use(dal.GetDBByCtx(ctx)).Comment
	count, err := postCommentDAL.WithContext(ctx).Where(postCommentDAL.Type.Eq(consts.CommentTypePost), postCommentDAL.Status.Eq(status)).Count()
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return count, nil
}
