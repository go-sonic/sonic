package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/service"
)

type StatisticHandler struct {
	StatisticService service.StatisticService
}

func NewStatisticHandler(l service.StatisticService) *StatisticHandler {
	return &StatisticHandler{
		StatisticService: l,
	}
}

func (s *StatisticHandler) Statistics(ctx *gin.Context) (interface{}, error) {
	return s.StatisticService.Statistic(ctx)
}

func (s *StatisticHandler) StatisticsWithUser(ctx *gin.Context) (interface{}, error) {
	return s.StatisticService.StatisticWithUser(ctx)
}
