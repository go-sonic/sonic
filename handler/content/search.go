package content

import (
	"context"
	"html"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type SearchHandler struct {
	PostAssembler assembler.PostAssembler
	PostService   service.PostService
	OptionService service.OptionService
	ThemeService  service.ThemeService
}

func NewSearchHandler(
	postAssembler assembler.PostAssembler,
	postService service.PostService,
	optionService service.OptionService,
	themeService service.ThemeService,
) *SearchHandler {
	return &SearchHandler{
		PostAssembler: postAssembler,
		PostService:   postService,
		OptionService: optionService,
		ThemeService:  themeService,
	}
}

func (s *SearchHandler) Search(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	return s.search(_ctx, ctx, 0, model)
}

func (s *SearchHandler) PageSearch(_ctx context.Context, ctx *app.RequestContext, model template.Model) (string, error) {
	page, err := util.ParamInt32(_ctx, ctx, "page")
	if err != nil {
		return "", err
	}
	return s.search(_ctx, ctx, int(page)-1, model)
}

func (s *SearchHandler) search(_ctx context.Context, ctx *app.RequestContext, pageNum int, model template.Model) (string, error) {
	keyword, err := util.MustGetQueryString(_ctx, ctx, "keyword")
	if err != nil {
		return "", err
	}
	sort := param.Sort{}
	err = ctx.BindAndValidate(&sort)
	if err != nil {
		return "", xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if len(sort.Fields) == 0 {
		sort = s.OptionService.GetPostSort(_ctx)
	}
	defaultPageSize := s.OptionService.GetIndexPageSize(_ctx)
	page := param.Page{
		PageNum:  pageNum,
		PageSize: defaultPageSize,
	}
	postQuery := param.PostQuery{
		Page:     page,
		Sort:     &sort,
		Keyword:  &keyword,
		Statuses: []*consts.PostStatus{consts.PostStatusPublished.Ptr()},
	}
	posts, total, err := s.PostService.Page(_ctx, postQuery)
	if err != nil {
		return "", err
	}
	postVOs, err := s.PostAssembler.ConvertToListVO(_ctx, posts)
	if err != nil {
		return "", err
	}
	model["is_search"] = true
	model["keyword"] = html.EscapeString(keyword)
	model["posts"] = dto.NewPage(postVOs, total, page)
	model["meta_keywords"] = s.OptionService.GetOrByDefault(_ctx, property.SeoKeywords)
	model["meta_description"] = s.OptionService.GetOrByDefault(_ctx, property.SeoDescription)
	return s.ThemeService.Render(_ctx, "search")
}
