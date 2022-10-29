package service

import (
	"context"

	"github.com/go-sonic/sonic/consts"
)

type SheetCommentService interface {
	BaseCommentService
	CountByStatus(ctx context.Context, status consts.CommentStatus) (int64, error)
}
