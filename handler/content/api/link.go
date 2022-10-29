package api

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
)

type LinkHandler struct {
	LinkService service.LinkService
}

func NewLinkHandler(linkService service.LinkService) *LinkHandler {
	return &LinkHandler{
		LinkService: linkService,
	}
}

type linkParam struct {
	*param.Sort
}

func (l *LinkHandler) ListLinks(ctx *gin.Context) (interface{}, error) {
	p := linkParam{}
	if err := ctx.ShouldBindQuery(&p); err != nil {
		return nil, err
	}

	if p.Sort == nil || len(p.Sort.Fields) == 0 {
		p.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	links, err := l.LinkService.List(ctx, p.Sort)
	if err != nil {
		return nil, err
	}
	return l.LinkService.ConvertToDTOs(ctx, links), nil
}

func (l *LinkHandler) LinkTeamVO(ctx *gin.Context) (interface{}, error) {
	p := linkParam{}
	if err := ctx.ShouldBindQuery(&p); err != nil {
		return nil, err
	}

	if p.Sort == nil || len(p.Sort.Fields) == 0 {
		p.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	links, err := l.LinkService.List(ctx, p.Sort)
	if err != nil {
		return nil, err
	}
	return l.LinkService.ConvertToLinkTeamVO(ctx, links), nil
}
