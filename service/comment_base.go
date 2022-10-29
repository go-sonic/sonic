package service

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type BaseCommentService interface {
	CreateBy(ctx context.Context, commentParam *param.Comment) (*entity.Comment, error)
	Page(ctx context.Context, commentQuery param.CommentQuery, commentType consts.CommentType) ([]*entity.Comment, int64, error)
	GetByID(ctx context.Context, commentID int64) (*entity.Comment, error)
	LGetByIDs(ctx context.Context, commentIDs []int64) ([]*entity.Comment, error)
	GetByContentID(ctx context.Context, contentID int32, sort *param.Sort) ([]*entity.Comment, error)
	Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	UpdateStatus(ctx context.Context, commentID int64, commentStatus consts.CommentStatus) (*entity.Comment, error)
	UpdateStatusBatch(ctx context.Context, commentIDs []int64, commentStatus consts.CommentStatus) ([]*entity.Comment, error)
	Delete(ctx context.Context, commentID int64) error
	DeleteBatch(ctx context.Context, commentIDs []int64) error
	Update(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	BuildAvatarURL(ctx context.Context, gravatarMD5 string, gravatarSource, gravatarDefault *string) (string, error)
	ConvertParam(commentParam *param.Comment) *entity.Comment
	CountByPostID(ctx context.Context, postID int32) (int64, error)
	CountByStatusAndPostIDs(ctx context.Context, status consts.CommentStatus, postIDs []int32) (map[int32]int64, error)
	CountChildren(ctx context.Context, parentCommentIDs []int64) (map[int64]int64, error)
	GetChildren(ctx context.Context, parentCommentID int64, contentID int32) ([]*entity.Comment, error)
}
