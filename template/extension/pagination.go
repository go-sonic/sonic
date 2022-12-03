package extension

import (
	"context"
	"strconv"

	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
)

type paginationExtension struct {
	Template      *template.Template
	OptionService service.OptionService
}

func RegisterPaginationFunc(template *template.Template, optionService service.OptionService) {
	p := &paginationExtension{
		Template:      template,
		OptionService: optionService,
	}
	p.addIndexPagination()
	p.addArchivesPagination()
	p.addCategoryPostsPagination()
	p.addJournalsPagination()
	p.addPhotosPagination()
	p.addSearchPagination()
	p.addTagPostsPagination()
}

func (p *paginationExtension) addIndexPagination() {
	indexPagination := func(page, total, display int) (*vo.Pagination, error) {
		ctx := context.Background()
		suffix, err := p.OptionService.GetPathSuffix(ctx)
		if err != nil {
			return nil, err
		}
		return p.getPagination(ctx, page, total, display, "", suffix)
	}
	p.Template.AddFunc("indexPagination", indexPagination)
}

func (p *paginationExtension) addArchivesPagination() {
	archivesPagination := func(page, total, display int) (*vo.Pagination, error) {
		ctx := context.Background()
		suffix, err := p.OptionService.GetPathSuffix(ctx)
		if err != nil {
			return nil, err
		}
		prefix, err := p.OptionService.GetArchivePrefix(ctx)
		if err != nil {
			return nil, err
		}
		return p.getPagination(ctx, page, total, display, "/"+prefix, suffix)
	}
	p.Template.AddFunc("archivesPagination", archivesPagination)
}

func (p *paginationExtension) addSearchPagination() {
	searchPagination := func(page, total, display int, keyword string) (*vo.Pagination, error) {
		ctx := context.Background()
		suffix, err := p.OptionService.GetPathSuffix(ctx)
		if err != nil {
			return nil, err
		}
		return p.getPagination(ctx, page, total, display, "/search", suffix+"?keyword="+keyword)
	}
	p.Template.AddFunc("searchPagination", searchPagination)
}

func (p *paginationExtension) addTagPostsPagination() {
	tagPostsPagination := func(page, total, display int, slug string) (*vo.Pagination, error) {
		ctx := context.Background()
		suffix, err := p.OptionService.GetPathSuffix(ctx)
		if err != nil {
			return nil, err
		}
		tagPrefix, err := p.OptionService.GetOrByDefaultWithErr(ctx, property.TagsPrefix, property.TagsPrefix.DefaultValue)
		if err != nil {
			return nil, err
		}
		return p.getPagination(ctx, page, total, display, "/"+tagPrefix.(string)+"/"+slug, suffix)
	}
	p.Template.AddFunc("tagPostsPagination", tagPostsPagination)
}

func (p *paginationExtension) addCategoryPostsPagination() {
	categoryPostsPagination := func(page, total, display int, slug string) (*vo.Pagination, error) {
		ctx := context.Background()
		suffix, err := p.OptionService.GetPathSuffix(ctx)
		if err != nil {
			return nil, err
		}
		categoryPrefix, err := p.OptionService.GetOrByDefaultWithErr(ctx, property.CategoriesPrefix, property.CategoriesPrefix.DefaultValue)
		if err != nil {
			return nil, err
		}
		return p.getPagination(ctx, page, total, display, "/"+categoryPrefix.(string)+"/"+slug, suffix)
	}
	p.Template.AddFunc("categoryPostsPagination", categoryPostsPagination)
}

func (p *paginationExtension) addPhotosPagination() {
	photosPagination := func(page, total, display int) (*vo.Pagination, error) {
		ctx := context.Background()
		suffix, err := p.OptionService.GetPathSuffix(ctx)
		if err != nil {
			return nil, err
		}
		prefix, err := p.OptionService.GetPhotoPrefix(ctx)
		if err != nil {
			return nil, err
		}
		return p.getPagination(ctx, page, total, display, "/"+prefix, suffix)
	}
	p.Template.AddFunc("photosPagination", photosPagination)
}

func (p *paginationExtension) addJournalsPagination() {
	journalsPagination := func(page, total, display int) (*vo.Pagination, error) {
		ctx := context.Background()
		suffix, err := p.OptionService.GetPathSuffix(ctx)
		if err != nil {
			return nil, err
		}
		prefix, err := p.OptionService.GetJournalPrefix(ctx)
		if err != nil {
			return nil, err
		}
		return p.getPagination(ctx, page, total, display, "/"+prefix, suffix)
	}
	p.Template.AddFunc("journalsPagination", journalsPagination)
}

func (p *paginationExtension) getPagination(ctx context.Context, page, total, display int, prefix, suffix string) (*vo.Pagination, error) {
	pagination := &vo.Pagination{}

	var nextPageFullPath, prevPageFullPath, fullPath string

	enable, err := p.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return nil, err
	}

	if enable {
		blogBaseURL, err := p.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return nil, err
		}
		prevPageFullPath, nextPageFullPath, fullPath = blogBaseURL, blogBaseURL, blogBaseURL
	}

	rainbow := util.RainbowPage(page+1, total, display)

	rainbowPages := make([]vo.RainbowPage, len(rainbow))

	nextPageFullPath += prefix + "/page/" + strconv.Itoa(page+2) + suffix

	if page == 1 {
		prevPageFullPath += prefix + "/"
	} else {
		prevPageFullPath += prefix + "/page/" + strconv.Itoa(page) + suffix
	}

	fullPath += prefix + "/page/"

	for i, current := range rainbow {
		rainbowPages[i] = vo.RainbowPage{
			Page:      current,
			FullPath:  fullPath + strconv.Itoa(current) + suffix,
			IsCurrent: page+1 == current,
		}
	}
	pagination.NextPageFullPath = nextPageFullPath
	pagination.PrevPageFullPath = prevPageFullPath
	pagination.RainbowPages = rainbowPages
	pagination.HasNext = total != page+1
	pagination.HasPrev = page != 0
	return pagination, nil
}
