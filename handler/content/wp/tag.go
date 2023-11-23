package wp

import (
	"github.com/gin-gonic/gin"
	"github.com/go-sonic/sonic/model/dto/wp"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
	"strings"
)

type TagHandler struct {
	TagService service.TagService
}

func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{
		TagService: tagService,
	}
}

func (handler *TagHandler) List(ctx *gin.Context) (interface{}, error) {
	var err error
	var listParam param.TagListParam
	if err = ctx.ShouldBind(&listParam); err != nil {
		return nil, err
	}

	entities, err := handler.TagService.ListByOption(ctx, &listParam)
	if err != nil {
		return nil, err
	}

	tagDTOList := make([]*wp.TagDTO, 0, len(entities))
	for _, tagEntity := range entities {
		tagDTOList = append(tagDTOList, convertToWpTag(tagEntity))
	}

	return tagDTOList, nil
}

func (handler *TagHandler) Create(ctx *gin.Context) (interface{}, error) {
	var err error
	var createParam param.TagCreateParam
	if err = ctx.ShouldBindJSON(&createParam); err != nil {
		return nil, err
	}

	createParam.Name = strings.TrimSpace(createParam.Name)
	if createParam.Name == "" {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("blank name")
	}

	tagEntity, err := handler.TagService.GetByName(ctx, createParam.Name)
	if err != nil {
		return nil, err
	}

	if tagEntity != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("tag exists")
	}

	tagParam := &param.Tag{
		Name:      createParam.Name,
		Slug:      createParam.Slug,
		Thumbnail: "",
		Color:     "",
	}
	create, err := handler.TagService.Create(ctx, tagParam)
	if err != nil {
		return nil, err
	}
	return convertToWpTag(create), nil
}

func convertToWpTag(tagEntity *entity.Tag) *wp.TagDTO {
	tagDTO := &wp.TagDTO{
		ID:          tagEntity.ID,
		Count:       0,
		Description: "",
		Link:        "",
		Name:        tagEntity.Name,
		Slug:        tagEntity.Slug,
		Taxonomy:    "",
		Meta:        nil,
	}
	return tagDTO
}
