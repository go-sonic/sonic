package service

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type PostService interface {
	BasePostService
	Page(ctx context.Context, postQuery param.PostQuery) ([]*entity.Post, int64, error)
	IncreaseLike(ctx context.Context, postID int32) error
	GetPrevPosts(ctx context.Context, post *entity.Post, size int) ([]*entity.Post, error)
	GetNextPosts(ctx context.Context, post *entity.Post, size int) ([]*entity.Post, error)
	Create(ctx context.Context, postParam *param.Post) (*entity.Post, error)
	Update(ctx context.Context, postID int32, postParam *param.Post) (*entity.Post, error)
	CountByStatus(ctx context.Context, status consts.PostStatus) (int64, error)
	CountVisit(ctx context.Context) (int64, error)
	Preview(ctx context.Context, postID int32) (string, error)
	CountLike(ctx context.Context) (int64, error)
}
