package service

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type PostCommentService interface {
	BaseCommentService
	CreateBy(ctx context.Context, commentParam *param.Comment) (*entity.Comment, error)
	CountByPostID(ctx context.Context, postID int32) (int64, error)
	CountByStatus(ctx context.Context, status consts.CommentStatus) (int64, error)
	UpdateBy(ctx context.Context, commentID int64, commentParam *param.Comment) (*entity.Comment, error)
}
