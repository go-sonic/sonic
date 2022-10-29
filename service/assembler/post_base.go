package assembler

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/service"
)

type BasePostAssembler interface {
	ConvertToSimpleDTO(ctx context.Context, post *entity.Post) (*dto.Post, error)
	ConvertToMinimalDTO(ctx context.Context, post *entity.Post) (*dto.PostMinimal, error)
	ConvertToDetailDTO(ctx context.Context, post *entity.Post) (*dto.PostDetail, error)
}

func NewBasePostAssembler(
	basePostService service.BasePostService,
	baseCommentService service.BaseCommentService,
) BasePostAssembler {
	return &basePostAssembler{
		BasePostService:    basePostService,
		BaseCommentService: baseCommentService,
	}
}

type basePostAssembler struct {
	BasePostService    service.BasePostService
	BaseCommentService service.BaseCommentService
}

func (p *basePostAssembler) ConvertToSimpleDTO(ctx context.Context, post *entity.Post) (*dto.Post, error) {
	postMinimal, err := p.ConvertToMinimalDTO(ctx, post)
	if err != nil {
		return nil, err
	}
	postDTO := &dto.Post{
		Summary:         post.Summary,
		Thumbnail:       post.Thumbnail,
		Visits:          post.Visits,
		DisallowComment: post.DisallowComment,
		Password:        post.Password,
		Template:        post.Template,
		TopPriority:     post.TopPriority,
		Likes:           post.Likes,
		WordCount:       post.WordCount,
		Topped:          post.TopPriority > 0,
	}
	postDTO.PostMinimal = *postMinimal

	if post.Summary == "" {
		postDTO.Summary = p.BasePostService.GenerateSummary(ctx, post.FormatContent)
	}
	return postDTO, nil
}

func (p *basePostAssembler) ConvertToMinimalDTO(ctx context.Context, post *entity.Post) (*dto.PostMinimal, error) {
	minimalPost := &dto.PostMinimal{
		ID:              post.ID,
		Title:           post.Title,
		Status:          post.Status,
		Slug:            post.Slug,
		EditorType:      post.EditorType,
		CreateTime:      post.CreateTime.UnixMilli(),
		MetaKeywords:    post.MetaKeywords,
		MetaDescription: post.MetaDescription,
	}
	if post.EditTime != nil {
		minimalPost.EditTime = post.EditTime.UnixMilli()
	}
	if post.UpdateTime != nil {
		minimalPost.UpdateTime = post.UpdateTime.UnixMilli()
	} else {
		minimalPost.UpdateTime = post.CreateTime.UnixMilli()
	}
	fullPath, err := p.BasePostService.BuildFullPath(ctx, post)
	if err != nil {
		return nil, err
	}
	minimalPost.FullPath = fullPath
	return minimalPost, nil
}

func (p *basePostAssembler) ConvertToDetailDTO(ctx context.Context, post *entity.Post) (*dto.PostDetail, error) {
	if post == nil {
		return nil, nil
	}
	postSimple, err := p.ConvertToSimpleDTO(ctx, post)
	if err != nil {
		return nil, err
	}
	postDetailDTO := &dto.PostDetail{
		Post:            *postSimple,
		OriginalContent: post.OriginalContent,
		Content:         post.FormatContent,
	}
	commentCount, err := p.BaseCommentService.CountByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	postDetailDTO.CommentCount = commentCount
	return postDetailDTO, nil
}
