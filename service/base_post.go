package service

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type BasePostService interface {
	GetByStatus(ctx context.Context, status []consts.PostStatus, postType consts.PostType, sort *param.Sort) ([]*entity.Post, error)
	GetByPostID(ctx context.Context, postID int32) (*entity.Post, error)
	GetByPostIDs(ctx context.Context, postIDs []int32) (map[int32]*entity.Post, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Post, error)
	GenerateSummary(ctx context.Context, htmlContent string) string
	BuildFullPath(ctx context.Context, post *entity.Post) (string, error)
	Delete(ctx context.Context, postID int32) error
	DeleteBatch(ctx context.Context, postIDs []int32) error
	UpdateDraftContent(ctx context.Context, postID int32, content, originalContent string) (*entity.Post, error)
	UpdateStatus(ctx context.Context, postID int32, status consts.PostStatus) (*entity.Post, error)
	UpdateStatusBatch(ctx context.Context, status consts.PostStatus, postIDs []int32) ([]*entity.Post, error)
	CreateOrUpdate(ctx context.Context, post *entity.Post, categoryIDs, tagIDs []int32, metas []param.Meta) (*entity.Post, error)
	IncreaseVisit(ctx context.Context, postID int32)
}
