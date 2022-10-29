package assembler

import (
	"context"
	"sort"
	"time"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
)

type PostAssembler interface {
	BasePostAssembler
	ConvertToListVO(ctx context.Context, posts []*entity.Post) ([]*vo.Post, error)
	ConvertToDetailVO(ctx context.Context, post *entity.Post) (*vo.PostDetailVO, error)
	ConvertToDetailVOs(ctx context.Context, posts []*entity.Post) ([]*vo.PostDetailVO, error)
	ConvertToArchiveYearVOs(ctx context.Context, posts []*entity.Post) ([]*vo.ArchiveYear, error)
	ConvertTOArchiveMonthVOs(ctx context.Context, posts []*entity.Post) ([]*vo.ArchiveMonth, error)
}

func NewPostAssembler(basePostService service.BasePostService,
	baseCommentService service.BaseCommentService,
	postTagService service.PostTagService,
	postCategoryService service.PostCategoryService,
	tagService service.TagService,
	categoryService service.CategoryService,
	postCommentService service.PostCommentService,
	metaService service.MetaService,
	basePostAssembler BasePostAssembler,
) PostAssembler {
	return &postAssembler{
		BasePostAssembler:   basePostAssembler,
		BasePostService:     basePostService,
		BaseCommentService:  baseCommentService,
		PostTagService:      postTagService,
		PostCommentService:  postCommentService,
		PostCategoryService: postCategoryService,
		TagService:          tagService,
		CategoryService:     categoryService,
		MetaService:         metaService,
	}
}

type postAssembler struct {
	BasePostAssembler
	BasePostService     service.BasePostService
	BaseCommentService  service.BaseCommentService
	PostTagService      service.PostTagService
	PostCategoryService service.PostCategoryService
	TagService          service.TagService
	CategoryService     service.CategoryService
	PostCommentService  service.PostCommentService
	MetaService         service.MetaService
}

func (p *postAssembler) ConvertToListVO(ctx context.Context, posts []*entity.Post) ([]*vo.Post, error) {
	postVOs := make([]*vo.Post, 0)
	postIDs := make([]int32, 0)
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}
	postTagMap, err := p.PostTagService.ListTagMapByPostID(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	tagDTOMap := make(map[int32]*dto.Tag)
	for _, tags := range postTagMap {
		for _, tag := range tags {
			if _, ok := tagDTOMap[tag.ID]; !ok {
				tagDTO, err := p.TagService.ConvertToDTO(ctx, tag)
				if err != nil {
					return nil, err
				}
				tagDTOMap[tag.ID] = tagDTO
			}
		}
	}
	categoryMap, err := p.PostCategoryService.ListCategoryMapByPostID(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	categoryDTOMap := make(map[int32]*dto.CategoryDTO, 0)
	for _, categories := range categoryMap {
		for _, category := range categories {
			if _, ok := categoryDTOMap[category.ID]; !ok {
				categoryDTO, err := p.CategoryService.ConvertToCategoryDTO(ctx, category)
				if err != nil {
					return nil, err
				}
				categoryDTOMap[category.ID] = categoryDTO
			}
		}
	}

	commentCountMap, err := p.PostCommentService.CountByStatusAndPostIDs(ctx, consts.CommentStatusPublished, postIDs)
	if err != nil {
		return nil, err
	}

	postMetaMap, err := p.MetaService.GetPostsMeta(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		postVO := &vo.Post{}
		if commentCount, ok := commentCountMap[post.ID]; ok {
			postVO.CommentCount = commentCount
		}
		if categories, ok := categoryMap[post.ID]; ok {
			categoryDTOs := make([]*dto.CategoryDTO, 0)
			for _, category := range categories {
				if categoryDTO, ok := categoryDTOMap[category.ID]; ok {
					categoryDTOs = append(categoryDTOs, categoryDTO)
				}
			}
			postVO.Categories = categoryDTOs
		}
		if tags, ok := postTagMap[post.ID]; ok {
			tagDTOs := make([]*dto.Tag, 0)
			for _, tag := range tags {
				if tagDTO, ok := tagDTOMap[tag.ID]; ok {
					tagDTOs = append(tagDTOs, tagDTO)
				}
			}
			postVO.Tags = tagDTOs
		}
		if metas, ok := postMetaMap[post.ID]; ok {
			metaMap := make(map[string]interface{})
			for _, meta := range metas {
				metaMap[meta.MetaKey] = meta
			}
			postVO.Metas = metaMap
		}
		postDTO, err := p.ConvertToSimpleDTO(ctx, post)
		if err != nil {
			return nil, err
		}
		postVO.Post = *postDTO
		postVOs = append(postVOs, postVO)
	}
	return postVOs, nil
}

func (p *postAssembler) ConvertToDetailVO(ctx context.Context, post *entity.Post) (*vo.PostDetailVO, error) {
	if post == nil {
		return nil, nil
	}
	postDetailVO := &vo.PostDetailVO{}
	postDetailDTO, err := p.ConvertToDetailDTO(ctx, post)
	if err != nil {
		return nil, err
	}
	postDetailVO.PostDetail = *postDetailDTO

	tags, err := p.PostTagService.ListTagByPostID(ctx, post.ID)
	tagIDs := make([]int32, 0)
	tagDTOs := make([]*dto.Tag, 0)
	for _, tag := range tags {
		tagIDs = append(tagIDs, tag.ID)
		tagDTO, err := p.TagService.ConvertToDTO(ctx, tag)
		if err != nil {
			return nil, err
		}
		tagDTOs = append(tagDTOs, tagDTO)
	}
	postDetailVO.TagIDs = tagIDs
	postDetailVO.Tags = tagDTOs
	categories, err := p.PostCategoryService.ListCategoryByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	categoryDTOs := make([]*dto.CategoryDTO, 0, len(categories))
	categoryIDs := make([]int32, 0)
	for _, category := range categories {
		categoryDTO, err := p.CategoryService.ConvertToCategoryDTO(ctx, category)
		if err != nil {
			return nil, err
		}
		categoryDTOs = append(categoryDTOs, categoryDTO)
		categoryIDs = append(categoryIDs, category.ID)
	}
	postDetailVO.CategoryIDs = categoryIDs
	postDetailVO.Categories = categoryDTOs

	postMetas, err := p.MetaService.GetPostMeta(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	metaDTOs := make([]*dto.Meta, 0, len(postMetas))
	metaIDs := make([]int64, 0, len(postMetas))
	for _, meta := range postMetas {
		metaDTOs = append(metaDTOs, p.MetaService.ConvertToMetaDTO(meta))
		metaIDs = append(metaIDs, meta.ID)
	}
	postDetailVO.MetaIDs = metaIDs
	postDetailVO.Metas = metaDTOs

	return postDetailVO, nil
}

func (p *postAssembler) ConvertToDetailVOs(ctx context.Context, posts []*entity.Post) ([]*vo.PostDetailVO, error) {
	postDetailVOs := make([]*vo.PostDetailVO, 0)
	postIDs := make([]int32, 0)
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}
	postTagMap, err := p.PostTagService.ListTagMapByPostID(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	tagDTOMap := make(map[int32]*dto.Tag)
	for _, tags := range postTagMap {
		for _, tag := range tags {
			if _, ok := tagDTOMap[tag.ID]; !ok {
				tagDTO, err := p.TagService.ConvertToDTO(ctx, tag)
				if err != nil {
					return nil, err
				}
				tagDTOMap[tag.ID] = tagDTO
			}
		}
	}
	categoryMap, err := p.PostCategoryService.ListCategoryMapByPostID(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	categoryDTOMap := make(map[int32]*dto.CategoryDTO, 0)
	for _, categories := range categoryMap {
		for _, category := range categories {
			if _, ok := categoryDTOMap[category.ID]; !ok {
				categoryDTO, err := p.CategoryService.ConvertToCategoryDTO(ctx, category)
				if err != nil {
					return nil, err
				}
				categoryDTOMap[category.ID] = categoryDTO
			}
		}
	}

	postMetaMap, err := p.MetaService.GetPostsMeta(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		postDetailVO := &vo.PostDetailVO{}
		if categories, ok := categoryMap[post.ID]; ok {
			categoryDTOs := make([]*dto.CategoryDTO, 0)
			categoryIDs := make([]int32, 0)
			for _, category := range categories {
				if categoryDTO, ok := categoryDTOMap[category.ID]; ok {
					categoryDTOs = append(categoryDTOs, categoryDTO)
					categoryIDs = append(categoryIDs, categoryDTO.ID)
				}
			}
			postDetailVO.Categories = categoryDTOs
			postDetailVO.CategoryIDs = categoryIDs
		}
		if tags, ok := postTagMap[post.ID]; ok {
			tagDTOs := make([]*dto.Tag, 0)
			tagIDs := make([]int32, 0)
			for _, tag := range tags {
				if tagDTO, ok := tagDTOMap[tag.ID]; ok {
					tagDTOs = append(tagDTOs, tagDTO)
					tagIDs = append(tagIDs, tagDTO.ID)
				}
			}
			postDetailVO.Tags = tagDTOs
			postDetailVO.TagIDs = tagIDs
		}
		if metas, ok := postMetaMap[post.ID]; ok {
			metaDTOs := make([]*dto.Meta, 0)
			metaIDs := make([]int64, 0)
			for _, meta := range metas {
				metaDTOs = append(metaDTOs, p.MetaService.ConvertToMetaDTO(meta))
				metaIDs = append(metaIDs, meta.ID)
			}
			postDetailVO.Metas = metaDTOs
			postDetailVO.MetaIDs = metaIDs
		}
		postDetailDTO, err := p.ConvertToDetailDTO(ctx, post)
		if err != nil {
			return nil, err
		}
		postDetailVO.PostDetail = *postDetailDTO
		postDetailVOs = append(postDetailVOs, postDetailVO)
	}
	return postDetailVOs, nil
}

func (p *postAssembler) ConvertToArchiveYearVOs(ctx context.Context, posts []*entity.Post) ([]*vo.ArchiveYear, error) {
	postVos, err := p.ConvertToListVO(ctx, posts)
	if err != nil {
		return nil, err
	}
	archiveYearVos := make([]*vo.ArchiveYear, 0)
	yearToPostMap := make(map[int][]*vo.Post)
	for _, postVo := range postVos {
		yearToPostMap[time.UnixMilli(postVo.CreateTime).Year()] = append(yearToPostMap[time.UnixMilli(postVo.CreateTime).Year()], postVo)
	}
	for year, postVos := range yearToPostMap {
		sort.Slice(postVos, func(i, j int) bool {
			return postVos[i].CreateTime >= postVos[j].CreateTime
		})
		archiveYearVos = append(archiveYearVos, &vo.ArchiveYear{Year: year, Posts: postVos})
	}
	sort.Slice(archiveYearVos, func(i, j int) bool {
		return archiveYearVos[i].Year >= archiveYearVos[j].Year
	})
	return archiveYearVos, nil
}

func (p *postAssembler) ConvertTOArchiveMonthVOs(ctx context.Context, posts []*entity.Post) ([]*vo.ArchiveMonth, error) {
	postVos, err := p.ConvertToListVO(ctx, posts)
	if err != nil {
		return nil, err
	}
	archiveMonthVos := make([]*vo.ArchiveMonth, 0)
	monthToPostMap := make(map[int][]*vo.Post)
	for _, postVo := range postVos {
		monthToPostMap[int(time.UnixMilli(postVo.CreateTime).Month())] = append(monthToPostMap[int(time.UnixMilli(postVo.CreateTime).Month())], postVo)
	}
	for year, postVos := range monthToPostMap {
		sort.Slice(postVos, func(i, j int) bool {
			return postVos[i].CreateTime >= postVos[j].CreateTime
		})
		archiveMonthVos = append(archiveMonthVos, &vo.ArchiveMonth{Month: year, Posts: postVos})
	}
	sort.Slice(archiveMonthVos, func(i, j int) bool {
		return archiveMonthVos[i].Month >= archiveMonthVos[j].Month
	})
	return archiveMonthVos, nil
}
