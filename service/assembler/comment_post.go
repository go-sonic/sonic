package assembler

import (
	"context"

	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
)

type PostCommentAssembler interface {
	BaseCommentAssembler
	ConvertToWithPost(ctx context.Context, comments []*entity.Comment) ([]*vo.PostCommentWithPost, error)
}

func NewPostCommentAssembler(
	optionService service.OptionService,
	baseCommentService service.BaseCommentService,
	baseCommentAssembler BaseCommentAssembler,
	postService service.PostService,
	postAssembler PostAssembler,
) PostCommentAssembler {
	return &postCommentAssembler{
		OptionService:        optionService,
		BaseCommentService:   baseCommentService,
		BaseCommentAssembler: baseCommentAssembler,
		PostAssembler:        postAssembler,
		PostService:          postService,
	}
}

type postCommentAssembler struct {
	OptionService      service.OptionService
	BaseCommentService service.BaseCommentService
	BaseCommentAssembler
	PostAssembler
	PostService service.PostService
}

func (p *postCommentAssembler) ConvertToWithPost(ctx context.Context, comments []*entity.Comment) ([]*vo.PostCommentWithPost, error) {
	postIDs := make([]int32, 0, len(comments))
	for _, comment := range comments {
		postIDs = append(postIDs, comment.PostID)
	}
	posts, err := p.PostService.GetByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	result := make([]*vo.PostCommentWithPost, 0, len(comments))
	for _, comment := range comments {
		commentDTO, err := p.BaseCommentAssembler.ConvertToDTO(ctx, comment)
		if err != nil {
			return nil, err
		}
		commentWithPost := &vo.PostCommentWithPost{
			Comment: *commentDTO,
		}
		result = append(result, commentWithPost)
		post, ok := posts[comment.PostID]
		if ok {
			commentWithPost.Post, err = p.PostAssembler.ConvertToMinimalDTO(ctx, post)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
