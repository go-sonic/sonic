package listener

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/service"
)

type PostUpdateListener struct {
	CategoryService     service.CategoryService
	PostCategoryService service.PostCategoryService
	PostService         service.PostService
}

func NewPostUpdateListener(bus event.Bus,
	categoryService service.CategoryService,
	postCategoryService service.PostCategoryService,
	postService service.PostService,
) {
	p := &PostUpdateListener{
		PostCategoryService: postCategoryService,
		PostService:         postService,
		CategoryService:     categoryService,
	}

	bus.Subscribe(event.PostUpdateEventName, p.HandlePostUpdateEvent)
}

func (p *PostUpdateListener) HandlePostUpdateEvent(ctx context.Context, postUpdateEvent event.Event) error {
	postID := postUpdateEvent.(*event.PostUpdateEvent).PostID

	categories, err := p.PostCategoryService.ListCategoryByPostID(ctx, postID)
	if err != nil {
		return err
	}

	postStatus := consts.PostStatusPublished
	for _, category := range categories {
		if category.Type == consts.CategoryTypeIntimate {
			postStatus = consts.PostStatusIntimate
		}
	}
	post, err := p.PostService.GetByPostID(ctx, postID)
	if err != nil {
		return err
	}
	if post.Status == consts.PostStatusRecycle || post.Status == consts.PostStatusDraft {
		return nil
	}
	if post.Password != "" {
		postStatus = consts.PostStatusIntimate
	}
	if post.Status == postStatus {
		return nil
	}
	_, err = p.PostService.UpdateStatus(ctx, postID, postStatus)
	return err
}
