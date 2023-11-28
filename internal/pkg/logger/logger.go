package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func NewLogger() {
	core := zapcore.NewCore(getEncoder(), getLogWriter(), zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log/run.log", // 日志文件的位置
		MaxSize:    10,              // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 100,             // 保留旧文件的最大个数
		MaxAge:     365,             // 保留旧文件的最大天数
		Compress:   false,           // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}
