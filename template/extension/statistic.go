package extension

import (
	"context"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
)

type statisticExtension struct {
	Template         *template.Template
	StatisticService service.StatisticService
}

func RegisterStatisticFunc(template *template.Template, statisticService service.StatisticService) {
	s := &statisticExtension{
		Template:         template,
		StatisticService: statisticService,
	}
	s.addGetStatisticsData()
}

func (s *statisticExtension) addGetStatisticsData() {
	getStatisticsDataFunc := func() (*dto.Statistic, error) {
		ctx := context.Background()
		statistic, err := s.StatisticService.Statistic(ctx)
		if err != nil {
			return nil, err
		}
		return statistic, nil
	}
	s.Template.AddFunc("getStatisticsData", getStatisticsDataFunc)
}
