package impl

import (
	"context"

	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/service"
)

type metaServiceImpl struct{}

func NewMetaService() service.MetaService {
	return &metaServiceImpl{}
}

func (m *metaServiceImpl) GetPostsMeta(ctx context.Context, postIDs []int32) (map[int32][]*entity.Meta, error) {
	metaDAL := dal.Use(dal.GetDBByCtx(ctx)).Meta
	metas, err := metaDAL.WithContext(ctx).Where(metaDAL.PostID.In(postIDs...)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	result := make(map[int32][]*entity.Meta, 0)
	for _, meta := range metas {
		postMetas, ok := result[meta.PostID]
		if !ok {
			postMetas = make([]*entity.Meta, 0)
		}
		postMetas = append(postMetas, meta)
		result[meta.PostID] = postMetas
	}
	return result, nil
}

func (m *metaServiceImpl) ConvertToMetaDTO(meta *entity.Meta) *dto.Meta {
	if meta == nil {
		return nil
	}
	metaDTO := &dto.Meta{
		ID:         meta.ID,
		PostID:     meta.PostID,
		Key:        meta.MetaKey,
		Value:      meta.MetaValue,
		CreateTime: meta.CreateTime.UnixMilli(),
	}
	return metaDTO
}

func (m *metaServiceImpl) GetPostMeta(ctx context.Context, postID int32) ([]*entity.Meta, error) {
	metaDAL := dal.Use(dal.GetDBByCtx(ctx)).Meta
	metas, err := metaDAL.WithContext(ctx).Where(metaDAL.PostID.Eq(postID)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return metas, nil
}

func (m *metaServiceImpl) ConvertToMetaDTOs(metas []*entity.Meta) []*dto.Meta {
	result := make([]*dto.Meta, len(metas))
	for i, meta := range metas {
		result[i] = m.ConvertToMetaDTO(meta)
	}
	return result
}
