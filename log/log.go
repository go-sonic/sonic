package log

import (
	"context"

	"go.uber.org/zap"
)

var (
	exportUseLogger      *zap.Logger
	exportUseSugarLogger *zap.SugaredLogger
)

func Debugf(template string, args ...interface{}) {
	exportUseSugarLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	exportUseSugarLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	exportUseSugarLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	exportUseSugarLogger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	exportUseSugarLogger.Fatalf(template, args...)
}

func Debug(msg string, fields ...zap.Field) {
	exportUseLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	exportUseLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	exportUseLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	exportUseLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	exportUseLogger.Fatal(msg, fields...)
}

func CtxDebugf(ctx context.Context, template string, args ...interface{}) {
	exportUseSugarLogger.Debugf(template, args...)
}

func CtxInfof(ctx context.Context, template string, args ...interface{}) {
	exportUseSugarLogger.Infof(template, args...)
}

func CtxWarnf(ctx context.Context, template string, args ...interface{}) {
	exportUseSugarLogger.Warnf(template, args...)
}

func CtxErrorf(ctx context.Context, template string, args ...interface{}) {
	exportUseSugarLogger.Errorf(template, args...)
}

func CtxFatalf(ctx context.Context, template string, args ...interface{}) {
	exportUseSugarLogger.Fatalf(template, args...)
}

func CtxDebug(ctx context.Context, msg string, fields ...zap.Field) {
	exportUseLogger.Debug(msg, fields...)
}

func CtxInfo(ctx context.Context, msg string, fields ...zap.Field) {
	exportUseLogger.Info(msg, fields...)
}

func CtxWarn(ctx context.Context, msg string, fields ...zap.Field) {
	exportUseLogger.Warn(msg, fields...)
}

func CtxError(ctx context.Context, msg string, fields ...zap.Field) {
	exportUseLogger.Error(msg, fields...)
}

func CtxFatal(ctx context.Context, msg string, fields ...zap.Field) {
	exportUseLogger.Fatal(msg, fields...)
}

func Sync() {
	_ = exportUseLogger.Sync()
	_ = exportUseSugarLogger.Sync()
}
