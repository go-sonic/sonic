package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
)

type CommentHandler struct {
	BaseCommentService service.BaseCommentService
}

func NewCommentHandler(baseCommentService service.BaseCommentService) *CommentHandler {
	return &CommentHandler{
		BaseCommentService: baseCommentService,
	}
}

func (c *CommentHandler) Like(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	commentID, err := util.ParamInt32(_ctx, ctx, "commentID")
	if err != nil {
		return nil, err
	}
	return nil, c.BaseCommentService.IncreaseLike(_ctx, commentID)
}
