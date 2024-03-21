package content

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
)

type FeedHandler struct {
	OptionService       service.OptionService
	PostService         service.PostService
	PostCategoryService service.PostCategoryService
	CategoryService     service.CategoryService
	PostAssembler       assembler.PostAssembler
}

func NewFeedHandler(optionService service.OptionService, postService service.PostService, categoryService service.CategoryService, postCategoryService service.PostCategoryService, postAssembler assembler.PostAssembler) *FeedHandler {
	return &FeedHandler{
		OptionService:       optionService,
		PostService:         postService,
		CategoryService:     categoryService,
		PostCategoryService: postCategoryService,
		PostAssembler:       postAssembler,
	}
}

func (f *FeedHandler) Feed(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	_, err := f.Atom(_ctx, ctx, model)
	if err != nil {
		return "", err
	}
	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	return "common/web/rss", nil
}

func (f *FeedHandler) CategoryFeed(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	_, err := f.CategoryAtom(_ctx, ctx, model)
	if err != nil {
		return "", err
	}
	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	return "common/web/rss", nil
}

func (f *FeedHandler) Atom(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	rssPageSize := f.OptionService.GetOrByDefault(_ctx, property.RssPageSize).(int)
	postQuery := param.PostQuery{
		Page:     param.Page{PageNum: 0, PageSize: rssPageSize},
		Sort:     &param.Sort{Fields: []string{"createTime,desc"}},
		Statuses: []*consts.PostStatus{consts.PostStatusPublished.Ptr()},
	}
	posts, _, err := f.PostService.Page(_ctx, postQuery)
	if err != nil {
		return "", err
	}
	postDetailVOs, err := f.buildPost(_ctx, posts)
	if err != nil {
		return "", err
	}
	lastModified := f.getLastModifiedTime(posts)
	ctx.Header("Last-Modified", lastModified.Format(http.TimeFormat))
	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	model["lastModified"] = lastModified
	model["posts"] = postDetailVOs
	return "common/web/atom", nil
}

func (f *FeedHandler) CategoryAtom(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	slug, err := util.ParamString(_ctx, ctx, "slug")
	if err != nil {
		return "", err
	}
	slug = strings.TrimSuffix(slug, ".xml")
	category, err := f.CategoryService.GetBySlug(_ctx, slug)
	if err != nil {
		return "", err
	}
	categoryDTO, err := f.CategoryService.ConvertToCategoryDTO(_ctx, category)
	if err != nil {
		return "", err
	}

	posts, err := f.PostCategoryService.ListByCategoryID(_ctx, category.ID, consts.PostStatusPublished)
	if err != nil {
		return "", err
	}

	postDetailVOs, err := f.buildPost(_ctx, posts)
	if err != nil {
		return "", err
	}
	lastModified := f.getLastModifiedTime(posts)

	model["category"] = categoryDTO
	model["posts"] = postDetailVOs
	model["lastModified"] = lastModified
	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	return "common/web/atom", nil
}

func (f *FeedHandler) Robots(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	ctx.Header("Content-Type", "text/plain;charset=utf-8")
	return "common/web/robots", nil
}

func (f *FeedHandler) SitemapXML(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	posts, _, err := f.PostService.Page(_ctx, param.PostQuery{
		Page:     param.Page{PageNum: 0, PageSize: int(^uint(0) >> 1)},
		Sort:     &param.Sort{Fields: []string{"createTime,desc"}},
		Statuses: []*consts.PostStatus{consts.PostStatusPublished.Ptr()},
	})
	if err != nil {
		return "", err
	}
	postDetailVOs, err := f.buildPost(_ctx, posts)
	if err != nil {
		return "", err
	}
	model["posts"] = postDetailVOs

	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	return "common/web/sitemap_xml", nil
}

func (f *FeedHandler) SitemapHTML(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	posts, _, err := f.PostService.Page(_ctx, param.PostQuery{
		Page:     param.Page{PageNum: 0, PageSize: int(^uint(0) >> 1)},
		Sort:     &param.Sort{Fields: []string{"createTime,desc"}},
		Statuses: []*consts.PostStatus{consts.PostStatusPublished.Ptr()},
	})
	if err != nil {
		return "", err
	}
	postDetailVOs, err := f.buildPost(_ctx, posts)
	if err != nil {
		return "", err
	}
	model["posts"] = postDetailVOs
	return "common/web/sitemap_html", nil
}

func (f *FeedHandler) getLastModifiedTime(posts []*entity.Post) time.Time {
	lastModifiedTime := time.Time{}
	for _, post := range posts {
		if post.EditTime != nil {
			if post.EditTime.After(lastModifiedTime) {
				lastModifiedTime = *post.EditTime
			}
		} else {
			if post.CreateTime.After(lastModifiedTime) {
				lastModifiedTime = post.CreateTime
			}
		}
	}
	if lastModifiedTime == (time.Time{}) {
		lastModifiedTime = time.Now()
	}
	return lastModifiedTime
}

var xmlInValidChar = regexp.MustCompile("[\x00-\x1F\x7F]")

func (f *FeedHandler) buildPost(ctx context.Context, posts []*entity.Post) ([]*vo.PostDetailVO, error) {
	postDetailVOs, err := f.PostAssembler.ConvertToDetailVOs(ctx, posts)
	if err != nil {
		return nil, err
	}
	for _, postDetailVO := range postDetailVOs {
		postDetailVO.Content = xmlInValidChar.ReplaceAllString(postDetailVO.Content, "")
		postDetailVO.Summary = xmlInValidChar.ReplaceAllString(postDetailVO.Summary, "")
	}
	return postDetailVOs, nil
}
