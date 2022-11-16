package listener

import (
	"context"

	"gorm.io/gorm"

	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/model/entity"
)

type LogEventListener struct {
	db *gorm.DB
}

func NewLogEventListener(db *gorm.DB, bus event.Bus) {
	l := &LogEventListener{
		db: db,
	}
	bus.Subscribe(event.LogEventName, l.HandleEvent)
}

func (l *LogEventListener) HandleEvent(ctx context.Context, logEvent event.Event) error {
	log, ok := logEvent.(*event.LogEvent)
	if !ok {
		return nil
	}
	logDAL := dal.GetQueryByCtx(ctx).Log
	logEntity := &entity.Log{
		Content:   log.Content,
		IPAddress: log.IpAddress,
		LogKey:    log.LogKey,
		Type:      log.LogType,
	}
	return logDAL.WithContext(ctx).Create(logEntity)
}
