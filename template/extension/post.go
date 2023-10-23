package extension

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/template"
)

type postExtension struct {
	Template            *template.Template
	PostService         service.PostService
	PostTagService      service.PostTagService
	PostCategoryService service.PostCategoryService
	CategoryService     service.CategoryService
	TagService          service.TagService
	PostAssembler       assembler.PostAssembler
}

func RegisterPostFunc(template *template.Template, postService service.PostService, postTagService service.PostTagService, postCategoryService service.PostCategoryService, categoryService service.CategoryService, postAssembler assembler.PostAssembler, tagService service.TagService) {
	p := &postExtension{
		Template:            template,
		PostService:         postService,
		PostTagService:      postTagService,
		PostCategoryService: postCategoryService,
		CategoryService:     categoryService,
		PostAssembler:       postAssembler,
		TagService:          tagService,
	}
	p.addListLatestPost()
	p.addGetPostCount()
	p.addGetPostArchiveYear()
	p.addGetPostArchiveMonth()
	p.addListPostByCategoryID()
	p.addListPostByCategorySlug()
	p.addListPostByTagID()
	p.addListPostByTagSlug()
	p.addListMostPopularPost()
}

func (p *postExtension) addListLatestPost() {
	listLatestPostFunc := func(top int) ([]*vo.Post, error) {
		ctx := context.Background()
		posts, _, err := p.PostService.Page(ctx, param.PostQuery{
			Page: param.Page{
				PageNum:  0,
				PageSize: top,
			},
			Sort: &param.Sort{
				Fields: []string{"createTime,desc"},
			},
			Statuses: []*consts.PostStatus{consts.PostStatusPublished.Ptr()},
		})
		if err != nil {
			return nil, err
		}
		return p.PostAssembler.ConvertToListVO(ctx, posts)
	}
	p.Template.AddFunc("listLatestPost", listLatestPostFunc)
}

func (p *postExtension) addListMostPopularPost() {
	listMostPopularPost := func(top int) ([]*vo.Post, error) {
		ctx := context.Background()
		posts, _, err := p.PostService.Page(ctx, param.PostQuery{
			Page: param.Page{
				PageNum:  0,
				PageSize: top,
			},
			Sort: &param.Sort{
				Fields: []string{"visits,desc"},
			},
			Statuses: []*consts.PostStatus{consts.PostStatusPublished.Ptr()},
		})
		if err != nil {
			return nil, err
		}
		return p.PostAssembler.ConvertToListVO(ctx, posts)
	}
	p.Template.AddFunc("listMostPopularPost", listMostPopularPost)
}

func (p *postExtension) addGetPostCount() {
	getPostCountFunc := func() (int64, error) {
		ctx := context.Background()
		return p.PostService.CountByStatus(ctx, consts.PostStatusPublished)
	}
	p.Template.AddFunc("getPostCount", getPostCountFunc)
}

func (p *postExtension) addGetPostArchiveYear() {
	getPostArchiveYearFunc := func() ([]*vo.ArchiveYear, error) {
		ctx := context.Background()
		posts, err := p.PostService.GetByStatus(ctx, []consts.PostStatus{consts.PostStatusPublished}, consts.PostTypePost, nil)
		if err != nil {
			return nil, err
		}
		return p.PostAssembler.ConvertToArchiveYearVOs(ctx, posts)
	}
	p.Template.AddFunc("listYearArchives", getPostArchiveYearFunc)
}

func (p *postExtension) addGetPostArchiveMonth() {
	getPostArchiveMonthFunc := func() ([]*vo.ArchiveMonth, error) {
		ctx := context.Background()
		posts, err := p.PostService.GetByStatus(ctx, []consts.PostStatus{consts.PostStatusPublished}, consts.PostTypePost, nil)
		if err != nil {
			return nil, err
		}
		return p.PostAssembler.ConvertTOArchiveMonthVOs(ctx, posts)
	}
	p.Template.AddFunc("listMonthArchives", getPostArchiveMonthFunc)
}

func (p *postExtension) addListPostByCategoryID() {
	listPostByCategoryID := func(categoryID int32) ([]*vo.Post, error) {
		ctx := context.Background()
		posts, err := p.PostCategoryService.ListByCategoryID(ctx, categoryID, consts.PostStatusPublished)
		if err != nil {
			return nil, err
		}
		return p.PostAssembler.ConvertToListVO(ctx, posts)
	}
	p.Template.AddFunc("listPostByCategoryID", listPostByCategoryID)
}

func (p *postExtension) addListPostByCategorySlug() {
	listPostByCategorySlug := func(slug string) ([]*vo.Post, error) {
		ctx := context.Background()
		category, err := p.CategoryService.GetBySlug(ctx, slug)
		if err != nil {
			return nil, err
		}
		posts, err := p.PostCategoryService.ListByCategoryID(ctx, category.ID, consts.PostStatusPublished)
		if err != nil {
			return nil, err
		}
		return p.PostAssembler.ConvertToListVO(ctx, posts)
	}
	p.Template.AddFunc("listPostByCategorySlug", listPostByCategorySlug)
}

func (p *postExtension) addListPostByTagID() {
	listPostByTagID := func(tagID int32) ([]*vo.Post, error) {
		ctx := context.Background()
		posts, err := p.PostTagService.ListPostByTagID(ctx, tagID, consts.PostStatusPublished)
		if err != nil {
			return nil, err
		}
		return p.PostAssembler.ConvertToListVO(ctx, posts)
	}
	p.Template.AddFunc("listPostByTagID", listPostByTagID)
}

func (p *postExtension) addListPostByTagSlug() {
	listPostByTagSlug := func(slug string) ([]*vo.Post, error) {
		ctx := context.Background()
		tag, err := p.TagService.GetBySlug(ctx, slug)
		if err != nil {
			return nil, err
		}
		posts, err := p.PostTagService.ListPostByTagID(ctx, tag.ID, consts.PostStatusPublished)
		if err != nil {
			return nil, err
		}
		return p.PostAssembler.ConvertToListVO(ctx, posts)
	}
	p.Template.AddFunc("listPostByTagSlug", listPostByTagSlug)
}
