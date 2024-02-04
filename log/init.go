package log

import (
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/go-sonic/sonic/config"
)

func NewLogger(conf *config.Config) *zap.Logger {
	_, err := os.Stat(conf.Sonic.LogDir)
	if err != nil {
		if os.IsNotExist(err) && !config.LogToConsole() {
			err := os.MkdirAll(conf.Sonic.LogDir, os.ModePerm)
			if err != nil {
				panic("mkdir failed![%v]")
			}
		}
	}

	var core zapcore.Core

	if config.LogToConsole() {
		core = zapcore.NewCore(getDevEncoder(), os.Stdout, getLogLevel(conf.Log.Levels.App))
	} else {
		core = zapcore.NewCore(getProdEncoder(), getWriter(conf), zap.DebugLevel)
	}

	// 传入 zap.AddCaller() 显示打日志点的文件名和行数
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))

	exportUseLogger = logger.WithOptions(zap.AddCallerSkip(1))
	exportUseSugarLogger = exportUseLogger.Sugar()
	return logger
}

// getWriter 自定义Writer,分割日志
func getWriter(conf *config.Config) zapcore.WriteSyncer {
	rotatingLogger := &lumberjack.Logger{
		Filename: filepath.Join(conf.Sonic.LogDir, conf.Log.FileName),
		MaxSize:  conf.Log.MaxSize,
		MaxAge:   conf.Log.MaxAge,
		Compress: conf.Log.Compress,
	}
	return zapcore.AddSync(rotatingLogger)
}

// getProdEncoder 自定义日志编码器
func getProdEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getDevEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		panic("log level error")
	}
}
