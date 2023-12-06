package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/service"
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
}
