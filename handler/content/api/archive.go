package api

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
)

type ArchiveHandler struct {
	PostService   service.PostService
	PostAssembler assembler.PostAssembler
}

func NewArchiveHandler(postService service.PostService, postAssemeber assembler.PostAssembler) *ArchiveHandler {
	return &ArchiveHandler{
		PostService:   postService,
		PostAssembler: postAssemeber,
	}
}

func (a *ArchiveHandler) ListYearArchives(ctx *gin.Context) (interface{}, error) {
	posts, err := a.PostService.GetByStatus(ctx, []consts.PostStatus{consts.PostStatusPublished}, consts.PostTypePost, nil)
	if err != nil {
		return nil, err
	}
	return a.PostAssembler.ConvertToArchiveYearVOs(ctx, posts)
}

func (a *ArchiveHandler) ListMonthArchives(ctx *gin.Context) (interface{}, error) {
	posts, err := a.PostService.GetByStatus(ctx, []consts.PostStatus{consts.PostStatusPublished}, consts.PostTypePost, nil)
	if err != nil {
		return nil, err
	}
	return a.PostAssembler.ConvertTOArchiveMonthVOs(ctx, posts)
}
