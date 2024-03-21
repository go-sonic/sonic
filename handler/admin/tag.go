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

type TagHandler struct {
	PostTagService service.PostTagService
	TagService     service.TagService
}

func NewTagHandler(postTagService service.PostTagService, tagService service.TagService) *TagHandler {
	return &TagHandler{
		PostTagService: postTagService,
		TagService:     tagService,
	}
}

func (t *TagHandler) ListTags(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	sort := param.Sort{}
	err := ctx.BindAndValidate(&sort)
	if err != nil {
		return nil, xerr.WithMsg(err, "sort parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "createTime,desc")
	}
	more, _ := util.MustGetQueryBool(_ctx, ctx, "more")
	if more {
		return t.PostTagService.ListAllTagWithPostCount(_ctx, &sort)
	}
	tags, err := t.TagService.ListAll(_ctx, &sort)
	if err != nil {
		return nil, err
	}
	return t.TagService.ConvertToDTOs(_ctx, tags)
}

func (t *TagHandler) GetTagByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	tag, err := t.TagService.GetByID(_ctx, id)
	if err != nil {
		return nil, err
	}
	return t.TagService.ConvertToDTO(_ctx, tag)
}

func (t *TagHandler) CreateTag(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	tagParam := &param.Tag{}
	err := ctx.BindAndValidate(tagParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	tag, err := t.TagService.Create(_ctx, tagParam)
	if err != nil {
		return nil, err
	}
	return t.TagService.ConvertToDTO(_ctx, tag)
}

func (t *TagHandler) UpdateTag(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	tagParam := &param.Tag{}
	err = ctx.BindAndValidate(tagParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	tag, err := t.TagService.Update(_ctx, id, tagParam)
	if err != nil {
		return nil, err
	}
	return t.TagService.ConvertToDTO(_ctx, tag)
}

func (t *TagHandler) DeleteTag(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	return nil, t.TagService.Delete(_ctx, id)
}
