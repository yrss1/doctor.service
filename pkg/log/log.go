package log

import (
	"context"
	"os"
	"strconv"

	"go.elastic.co/apm/module/apmzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
)

func init() {
	defaultLogger = New()
}

type ctxLoggerKey struct{}

func ContextWithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, l)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxLoggerKey{}).(*zap.Logger); ok {
		return l
	}
	return defaultLogger
}

func New() *zap.Logger {
	cfg := zap.NewProductionConfig()

	debugMode, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err == nil && debugMode {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	//cfg.OutputPaths = []string{"stdout", "service.log"}

	log, err := cfg.Build(zap.WrapCore((&apmzap.Core{FatalFlushTimeout: 10000}).WrapCore))
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return log
}

func Sync() {
	if err := defaultLogger.Sync(); err != nil {
		defaultLogger.Error("Error flushing log buffer", zap.Error(err))
	}
}
