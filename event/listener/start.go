package listener

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
)

type StartListener struct {
	db            *gorm.DB
	optionService service.OptionService
	bus           event.Bus
}

func NewStartListener(db *gorm.DB, optionService service.OptionService, bus event.Bus) {
	s := StartListener{
		db:            db,
		optionService: optionService,
		bus:           bus,
	}
	bus.Subscribe(event.StartEventName, s.HandleEvent)
}

func (s *StartListener) HandleEvent(ctx context.Context, startEvent event.Event) error {
	if _, ok := startEvent.(*event.StartEvent); !ok {
		return nil
	}

	err := s.createOptions()
	if err != nil {
		log.Error("create options err", zap.Error(err))
	}
	if dal.DBType == consts.DBTypeMySQL {
		err = dal.DB.Session(&gorm.Session{Context: ctx}).Raw("SELECT VERSION()").Scan(&consts.DatabaseVersion).Error
	} else if dal.DBType == consts.DBTypeSQLite {
		err = dal.DB.Session(&gorm.Session{Context: ctx}).Raw("SELECT SQLITE_VERSION()").Scan(&consts.DatabaseVersion).Error
	}
	if err != nil {
		return err
	}
	_ = s.printStartInfo(ctx)
	return nil
}

func (s *StartListener) createOptions() error {
	ctx := context.Background()

	ctx = dal.SetCtxQuery(ctx, dal.GetQueryByCtx(ctx).ReplaceDB(dal.GetDB().Session(
		&gorm.Session{Logger: dal.DB.Logger.LogMode(logger.Warn)},
	)))

	optionDAL := dal.GetQueryByCtx(ctx).Option
	options, err := optionDAL.WithContext(ctx).Find()
	if err != nil {
		return err
	}
	toCreate := make([]*entity.Option, 0)
out:
	for _, p := range property.AllProperty {
		for _, o := range options {
			if p.KeyValue == o.OptionKey {
				continue out
			}
		}
		toCreate = append(toCreate, p.ConvertToOption())
	}
	return optionDAL.WithContext(ctx).Create(toCreate...)
}

func (s *StartListener) printStartInfo(ctx context.Context) error {
	blogURL, err := s.optionService.GetBlogBaseURL(ctx)
	if err != nil {
		return err
	}
	site := logger.BlueBold + "Sonic started at         " + blogURL + logger.Reset
	log.Info(site)
	fmt.Println(site)

	adminURLPath, err := s.optionService.GetAdminURLPath(ctx)
	if err != nil {
		return err
	}
	adminSite := logger.BlueBold + "Sonic admin started at         " + blogURL + "/" + adminURLPath + logger.Reset
	log.Info(adminSite)
	fmt.Println(adminSite)
	return nil
}
