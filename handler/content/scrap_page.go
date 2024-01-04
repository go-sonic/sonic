package content

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type ScrapHandler struct {
	ScrapService  service.ScrapService
	ThemeService  service.ThemeService
	OptionService service.OptionService
}

func NewScrapHandler(scrapService service.ScrapService, themeService service.ThemeService, optionService service.OptionService) *ScrapHandler {
	return &ScrapHandler{
		ScrapService:  scrapService,
		ThemeService:  themeService,
		OptionService: optionService,
	}
}

func (handler *ScrapHandler) Index(ctx *gin.Context, model template.Model) (string, error) {
	pageSize := handler.OptionService.GetIndexPageSize(ctx)

	page, err := util.GetQueryInt32(ctx, "page", 1)
	if err != nil {
		return "", xerr.WithStatus(nil, int(xerr.StatusBadRequest)).WithMsg("查询不到文章信息")
	}
	query := &param.ScrapPageQuery{
		Page: param.Page{
			PageNum:  int(page),
			PageSize: pageSize,
		},
	}
	pageList, total, err := handler.ScrapService.Query(ctx, query)
	if err != nil {
		return "", err
	}
	model["pages"] = dto.NewPage(pageList, total, param.Page{
		PageNum:  int(page),
		PageSize: pageSize,
	})
	return handler.ThemeService.Render(ctx, "scrap")
}
