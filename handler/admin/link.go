package admin

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type LinkHandler struct {
	LinkService service.LinkService
}

func NewLinkHandler(linkService service.LinkService) *LinkHandler {
	return &LinkHandler{
		LinkService: linkService,
	}
}

func (l *LinkHandler) ListLinks(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	sort := param.Sort{}
	err := ctx.BindAndValidate(&sort)
	if err != nil {
		return nil, xerr.WithMsg(err, "sort parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "team,desc", "priority,asc")
	} else {
		sort.Fields = append(sort.Fields, "priority,asc")
	}
	links, err := l.LinkService.List(_ctx, &sort)
	if err != nil {
		return nil, err
	}
	return l.LinkService.ConvertToDTOs(_ctx, links), nil
}

func (l *LinkHandler) GetLinkByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	link, err := l.LinkService.GetByID(_ctx, id)
	if err != nil {
		return nil, err
	}
	return l.LinkService.ConvertToDTO(_ctx, link), nil
}

func (l *LinkHandler) CreateLink(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	linkParam := &param.Link{}
	err := ctx.BindAndValidate(linkParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	link, err := l.LinkService.Create(_ctx, linkParam)
	if err != nil {
		return nil, err
	}
	return l.LinkService.ConvertToDTO(_ctx, link), nil
}

func (l *LinkHandler) UpdateLink(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	linkParam := &param.Link{}
	err = ctx.BindAndValidate(linkParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	link, err := l.LinkService.Update(_ctx, id, linkParam)
	if err != nil {
		return nil, err
	}
	return l.LinkService.ConvertToDTO(_ctx, link), nil
}

func (l *LinkHandler) DeleteLink(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	return nil, l.LinkService.Delete(_ctx, id)
}

func (l *LinkHandler) ListLinkTeams(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return l.LinkService.ListTeams(_ctx)
}
