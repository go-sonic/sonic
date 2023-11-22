package wp

import (
	"github.com/gin-gonic/gin"
	"github.com/go-sonic/sonic/model/dto/wp"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
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
	if err = ctx.ShouldBindJSON(&listParam); err != nil {
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
