package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
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
	scrapPageParam := param.ScrapPage{}
	err := ctx.ShouldBindJSON(&scrapPageParam)
	if err != nil {
		return nil, xerr.BadParam.Wrap(err)
	}

	err = handler.ScrapService.Create(ctx, &scrapPageParam)
	if err != nil {
		return nil, err
	}

	return true, nil
}
