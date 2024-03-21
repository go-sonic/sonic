package admin

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
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

func (s *StatisticHandler) Statistics(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return s.StatisticService.Statistic(_ctx)
}

func (s *StatisticHandler) StatisticsWithUser(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return s.StatisticService.StatisticWithUser(_ctx)
}
