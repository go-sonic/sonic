package event

import (
	"context"
	"reflect"
	"sync"

	"go.uber.org/zap"

	"github.com/go-sonic/sonic/log"
)

type Listener func(ctx context.Context, event Event) error

type Bus interface {
	Publish(ctx context.Context, event Event)
	Subscribe(eventType string, listener Listener)
	UnSubscribe(eventType string, listener Listener)
}

type syncLocalBus struct {
	listeners sync.Map
	logger    *zap.Logger
}

func NewSyncEventBus(logger *zap.Logger) Bus {
	return &syncLocalBus{
		logger: logger,
	}
}

func (e *syncLocalBus) Publish(ctx context.Context, event Event) {
	defer func() {
		if err := recover(); err != nil {
			log.CtxError(ctx, "event panic", zap.String("event", event.EventType()), zap.Stack("stack"), zap.Any("err", err))
		}
	}()
	if listeners, ok := e.listeners.Load(event.EventType()); ok {
		for _, listener := range listeners.([]Listener) {
			err := listener(ctx, event)
			if err != nil {
				e.logger.Error("error in event listener", zap.Any("event", event.EventType()), zap.Error(err))
			}
		}
	}
}

func (e *syncLocalBus) Subscribe(eventType string, listener Listener) {
	if listeners, ok := e.listeners.Load(eventType); ok {
		listeners = append(listeners.([]Listener), listener)
		e.listeners.Store(eventType, listeners)
	} else {
		listeners := make([]Listener, 0)
		listeners = append(listeners, listener)
		e.listeners.Store(eventType, listeners)
	}
}

func (e *syncLocalBus) UnSubscribe(eventType string, listener Listener) {
	if listeners, ok := e.listeners.Load(eventType); ok && len(listeners.([]Listener)) > 0 {
		target := reflect.ValueOf(listener).Pointer()
		var filtered []Listener
		for _, i := range listeners.([]Listener) {
			if reflect.ValueOf(i).Pointer() != target {
				filtered = append(filtered, i)
			}
		}
		e.listeners.Store(eventType, filtered)
	}
}
