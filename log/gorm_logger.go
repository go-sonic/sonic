package log

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/go-sonic/sonic/config"
)

type gormLogger struct {
	logger.Config
	traceStr     string
	traceWarnStr string
	traceErrStr  string
	zapLogger    *zap.Logger
}

func NewGormLogger(conf *config.Config, zapLogger *zap.Logger) logger.Interface {
	logConfig := logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  GetGormLogLevel(conf.Log.Levels.Gorm),
		IgnoreRecordNotFoundError: true,
		Colorful:                  config.LogToConsole(),
	}
	gl := &gormLogger{
		Config:       logConfig,
		traceStr:     "[%.3fms] [rows:%v] %s",
		traceWarnStr: "%s [%.3fms] [rows:%v] %s",
		traceErrStr:  "%s [%.3fms] [rows:%v] %s",
		zapLogger:    zapLogger,
	}
	if logConfig.Colorful {
		gl.traceStr = logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		gl.traceWarnStr = "%s " + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		gl.traceErrStr = logger.MagentaBold + "%s " + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}
	return gl
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

const level = 2

func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.zapLogger.WithOptions(zap.AddCallerSkip(getCallerSkip()-level)).Sugar().Infof(msg, data...)
}

func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.zapLogger.WithOptions(zap.AddCallerSkip(getCallerSkip()-level)).Sugar().Warnf(msg, data...)
}

func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.zapLogger.WithOptions(zap.AddCallerSkip(getCallerSkip()-level)).Sugar().Errorf(msg, data...)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.zapLogger.WithOptions(zap.AddCallerSkip(getCallerSkip()-level)).Sugar().Errorf(l.traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.zapLogger.WithOptions(zap.AddCallerSkip(getCallerSkip()-level)).Sugar().Errorf(l.traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.zapLogger.WithOptions(zap.AddCallerSkip(getCallerSkip()-level)).Sugar().Warnf(l.traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.zapLogger.WithOptions(zap.AddCallerSkip(getCallerSkip()-level)).Sugar().Warnf(l.traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.zapLogger.WithOptions(zap.AddCallerSkip(getCallerSkip()-level)).Sugar().Infof(l.traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.zapLogger.WithOptions(zap.AddCallerSkip(getCallerSkip()-level)).Sugar().Infof(l.traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func getCallerSkip() int {
	for i := 3; i < 15; i++ {
		pc := make([]uintptr, 1)
		numFrames := runtime.Callers(i, pc)
		if numFrames < 1 {
			return i
		}
		frame, _ := runtime.CallersFrames(pc).Next()
		if !strings.Contains(frame.Function, "gorm.io") && !strings.Contains(frame.Function, "github.com/go-sonic/sonic/dal") {
			return i
		}
	}
	return 0
}

func GetGormLogLevel(level string) logger.LogLevel {
	switch level {
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	case "silent":
		return logger.Silent
	default:
		panic("log level error")
	}
}
