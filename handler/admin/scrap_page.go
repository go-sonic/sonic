package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/handler/binding"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type ScrapPageHandler struct {
	ScrapService service.ScrapService
}

func NewScrapPageHandler(scrapService service.ScrapService) *ScrapPageHandler {
	return &ScrapPageHandler{
		ScrapService: scrapService,
	}
}

func (handler *ScrapPageHandler) QueryMd5List(ctx *gin.Context) (interface{}, error) {
	return handler.ScrapService.QueryMd5List(ctx)
}

func (handler *ScrapPageHandler) Create(ctx *gin.Context) (interface{}, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.BadParam.New("empty file").WithStatus(xerr.StatusBadRequest).WithMsg("empty files")
	}

	scrapPageParam := param.ScrapPage{}
	err = ctx.ShouldBindWith(&scrapPageParam, binding.CustomFormBinding)
	if err != nil {
		return nil, xerr.BadParam.Wrap(err)
	}

	pageDTO, err := handler.ScrapService.Create(ctx, &scrapPageParam, file)
	if err != nil {
		return nil, err
	}

	return pageDTO, nil
}

func (handler *ScrapPageHandler) Get(ctx *gin.Context) {
	pageID, err := util.ParamInt32(ctx, "id")
	if err != nil {
		ctx.String(400, "invalid pageId")
		return
	}
	pageDTO, err := handler.ScrapService.Get(ctx, pageID)
	if err != nil {
		ctx.String(500, "fail ")
		return
	}
	if len(pageDTO.Content) > 0 {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.String(200, pageDTO.Content)
	} else {
		ctx.String(200, "<p>empty</p>")
	}
}

func (handler *ScrapPageHandler) Query(ctx *gin.Context) (interface{}, error) {
	query := param.ScrapPageQuery{}
	err := ctx.ShouldBindWith(&query, binding.CustomFormBinding)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}

	result, total, err := handler.ScrapService.Query(ctx, &query)
	if err != nil {
		return nil, err
	}

	return dto.NewPage(result, total, query.Page), nil
}
