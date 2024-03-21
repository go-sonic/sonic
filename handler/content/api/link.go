package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
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

func (l *LinkHandler) ListLinks(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	p := linkParam{}
	if err := ctx.BindAndValidate(&p); err != nil {
		return nil, err
	}

	if p.Sort == nil || len(p.Sort.Fields) == 0 {
		p.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	links, err := l.LinkService.List(_ctx, p.Sort)
	if err != nil {
		return nil, err
	}
	return l.LinkService.ConvertToDTOs(_ctx, links), nil
}

func (l *LinkHandler) LinkTeamVO(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	p := linkParam{}
	if err := ctx.BindAndValidate(&p); err != nil {
		return nil, err
	}

	if p.Sort == nil || len(p.Sort.Fields) == 0 {
		p.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	links, err := l.LinkService.List(_ctx, p.Sort)
	if err != nil {
		return nil, err
	}
	return l.LinkService.ConvertToLinkTeamVO(_ctx, links), nil
}
