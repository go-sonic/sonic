package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
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

func (t *TagHandler) ListTags(ctx *gin.Context) (interface{}, error) {
	sort := param.Sort{}
	err := ctx.ShouldBindQuery(&sort)
	if err != nil {
		return nil, xerr.WithMsg(err, "sort parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "createTime,desc")
	}
	more, _ := util.MustGetQueryBool(ctx, "more")
	if more {
		return t.PostTagService.ListAllTagWithPostCount(ctx, &sort)
	}
	tags, err := t.TagService.ListAll(ctx, &sort)
	if err != nil {
		return nil, err
	}
	return t.TagService.ConvertToDTOs(ctx, tags)
}

func (t *TagHandler) GetTagByID(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "id")
	if err != nil {
		return nil, err
	}
	tag, err := t.TagService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return t.TagService.ConvertToDTO(ctx, tag)
}

func (t *TagHandler) CreateTag(ctx *gin.Context) (interface{}, error) {
	tagParam := &param.Tag{}
	err := ctx.ShouldBindJSON(tagParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	tag, err := t.TagService.Create(ctx, tagParam)
	if err != nil {
		return nil, err
	}
	return t.TagService.ConvertToDTO(ctx, tag)
}

func (t *TagHandler) UpdateTag(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "id")
	if err != nil {
		return nil, err
	}
	tagParam := &param.Tag{}
	err = ctx.ShouldBindJSON(tagParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	tag, err := t.TagService.Update(ctx, id, tagParam)
	if err != nil {
		return nil, err
	}
	return t.TagService.ConvertToDTO(ctx, tag)
}

func (t *TagHandler) DeleteTag(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "id")
	if err != nil {
		return nil, err
	}
	return nil, t.TagService.Delete(ctx, id)
}
