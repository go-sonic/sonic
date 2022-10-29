package impl

import (
	"context"
	"strings"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type tagServiceImpl struct {
	OptionService service.OptionService
}

func NewTagService(optionService service.OptionService) service.TagService {
	return &tagServiceImpl{
		OptionService: optionService,
	}
}

func (t tagServiceImpl) GetByID(ctx context.Context, id int32) (*entity.Tag, error) {
	tagDAL := dal.Use(dal.GetDBByCtx(ctx)).Tag
	tag, err := tagDAL.WithContext(ctx).Where(tagDAL.ID.Eq(id)).First()
	return tag, WrapDBErr(err)
}

func (t tagServiceImpl) GetBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	tagDAL := dal.Use(dal.GetDBByCtx(ctx)).Tag
	tag, err := tagDAL.WithContext(ctx).Where(tagDAL.Slug.Eq(slug)).First()
	return tag, WrapDBErr(err)
}

func (t tagServiceImpl) Create(ctx context.Context, tagParam *param.Tag) (*entity.Tag, error) {
	if tagParam.Slug == "" {
		tagParam.Slug = util.Slug(tagParam.Name)
	} else {
		tagParam.Slug = util.Slug(tagParam.Slug)
	}
	if tagParam.Color == "" {
		tagParam.Color = consts.SonicDefaultTagColor
	}
	tagDAL := dal.Use(dal.GetDBByCtx(ctx)).Tag
	tag := &entity.Tag{
		Name:      tagParam.Name,
		Slug:      tagParam.Slug,
		Thumbnail: tagParam.Thumbnail,
		Color:     tagParam.Color,
	}
	err := tagDAL.WithContext(ctx).Create(tag)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return tag, nil
}

func (t tagServiceImpl) Update(ctx context.Context, id int32, tagParam *param.Tag) (*entity.Tag, error) {
	if tagParam.Slug == "" {
		tagParam.Slug = util.Slug(tagParam.Name)
	} else {
		tagParam.Slug = util.Slug(tagParam.Slug)
	}
	if tagParam.Color == "" {
		tagParam.Color = consts.SonicDefaultTagColor
	}
	tagDAL := dal.Use(dal.GetDBByCtx(ctx)).Tag
	updateResult, err := tagDAL.WithContext(ctx).Where(tagDAL.ID.Eq(id)).UpdateSimple(
		tagDAL.Name.Value(tagParam.Name),
		tagDAL.Slug.Value(tagParam.Slug),
		tagDAL.Thumbnail.Value(tagParam.Thumbnail),
		tagDAL.Color.Value(tagParam.Color),
	)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return nil, xerr.NoType.New("update tag failed id=%v", id).WithStatus(xerr.StatusInternalServerError).WithMsg("update tag failed")
	}
	tag, err := tagDAL.WithContext(ctx).Where(tagDAL.ID.Value(id)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return tag, nil
}

func (t tagServiceImpl) Delete(ctx context.Context, id int32) error {
	tagDAL := dal.Use(dal.GetDBByCtx(ctx)).Tag
	deleteResult, err := tagDAL.WithContext(ctx).Where(tagDAL.ID.Value(id)).Delete()
	if err != nil {
		return WrapDBErr(err)
	}
	if deleteResult.RowsAffected != 1 {
		return xerr.NoType.New("delete tag failed id=%v", id).WithMsg("delete tag failed").WithStatus(xerr.StatusInternalServerError)
	}
	return nil
}

func (t tagServiceImpl) ListAll(ctx context.Context, sort *param.Sort) ([]*entity.Tag, error) {
	tagDAL := dal.Use(dal.GetDBByCtx(ctx)).Tag
	tagDO := tagDAL.WithContext(ctx)
	err := BuildSort(sort, &tagDAL, &tagDO)
	if err != nil {
		return nil, err
	}
	tags, err := tagDAL.WithContext(ctx).Where().Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return tags, nil
}

func (t tagServiceImpl) ListByIDs(ctx context.Context, tagIDs []int32) ([]*entity.Tag, error) {
	if len(tagIDs) == 0 {
		return make([]*entity.Tag, 0), nil
	}
	tagDAL := dal.Use(dal.GetDBByCtx(ctx)).Tag
	tags, err := tagDAL.WithContext(ctx).Where(tagDAL.ID.In(tagIDs...)).Find()
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (t tagServiceImpl) ConvertToDTO(ctx context.Context, tag *entity.Tag) (*dto.Tag, error) {
	tagDTO := &dto.Tag{
		ID:         tag.ID,
		Name:       tag.Name,
		Slug:       tag.Slug,
		Thumbnail:  tag.Thumbnail,
		CreateTime: tag.CreateTime.UnixMilli(),
		Color:      tag.Color,
	}
	fullPath := strings.Builder{}
	isEnabled, err := t.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return nil, err
	}
	if isEnabled {
		blogBaseUrl, err := t.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return nil, err
		}
		fullPath.WriteString(blogBaseUrl)
	}
	fullPath.WriteString("/")

	tagPrefix, err := t.OptionService.GetOrByDefaultWithErr(ctx, property.TagsPrefix, "tags")
	if err != nil {
		return nil, err
	}
	fullPath.WriteString(tagPrefix.(string))
	fullPath.WriteString("/")
	fullPath.WriteString(tag.Slug)
	pathSuffix, err := t.OptionService.GetPathSuffix(ctx)
	if err != nil {
		return nil, err
	}
	fullPath.WriteString(pathSuffix)
	tagDTO.FullPath = fullPath.String()
	return tagDTO, nil
}

func (t tagServiceImpl) ConvertToDTOs(ctx context.Context, tags []*entity.Tag) ([]*dto.Tag, error) {
	isEnabled, err := t.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return nil, err
	}
	var blogBaseUrl string
	if isEnabled {
		blogBaseUrl, err = t.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return nil, err
		}
	}

	tagPrefix, err := t.OptionService.GetOrByDefaultWithErr(ctx, property.TagsPrefix, "tags")
	if err != nil {
		return nil, err
	}
	pathSuffix, err := t.OptionService.GetPathSuffix(ctx)
	if err != nil {
		return nil, err
	}
	tagDTOs := make([]*dto.Tag, 0, len(tags))
	for _, tag := range tags {
		fullPath := strings.Builder{}
		if isEnabled {
			fullPath.WriteString(blogBaseUrl)
		}
		fullPath.WriteString("/")
		fullPath.WriteString(tagPrefix.(string))
		fullPath.WriteString("/")
		fullPath.WriteString(tag.Slug)
		fullPath.WriteString(pathSuffix)
		tagDTO := &dto.Tag{
			ID:         tag.ID,
			Name:       tag.Name,
			Slug:       tag.Slug,
			Thumbnail:  tag.Thumbnail,
			CreateTime: tag.CreateTime.UnixMilli(),
			FullPath:   fullPath.String(),
			Color:      tag.Color,
		}
		tagDTOs = append(tagDTOs, tagDTO)
	}
	return tagDTOs, nil
}

func (t tagServiceImpl) CountAllTag(ctx context.Context) (int64, error) {
	tagDAL := dal.Use(dal.GetDBByCtx(ctx)).Tag
	count, err := tagDAL.WithContext(ctx).Count()
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return count, nil
}
