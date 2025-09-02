package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once      sync.Once
	zapLogger *zap.Logger
)

type (
	Level int8
)

func (l Level) getZapLevel() zap.AtomicLevel {
	switch l {
	case DebugLevel:
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case InfoLevel:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case WarnLevel:
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case ErrorLevel:
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case FatalLevel:
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	}

	panic("invalid level")
}

const (
	// DebugLevel development messages
	DebugLevel Level = iota - 1
	// InfoLevel normal information messages for production use
	InfoLevel
	// WarnLevel low priority errors
	WarnLevel
	// ErrorLevel need attention
	ErrorLevel
	// FatalLevel fatal errors - recoverable or not
	FatalLevel
)

type Collector int

const (
	def     Collector = 0
	Elastic Collector = 1
)

type initOptions struct {
	collector Collector
	level     Level
}

type InitOption func(*initOptions)

func WithLevel(level Level) InitOption {
	return func(options *initOptions) {
		options.level = level
	}
}

func WithCollector(c Collector) InitOption {
	return func(options *initOptions) {
		options.collector = c
	}
}

func New(opts ...InitOption) (*zap.Logger, error) {
	var err error
	once.Do(func() {
		zapLogger = new(zap.Logger)
		var (
			zl *zap.Logger
		)

		options := &initOptions{
			level: InfoLevel,
		}

		for _, opt := range opts {
			opt(options)
		}

		if options.collector == Elastic {
			encoderConfig := ecszap.ECSCompatibleEncoderConfig(zap.NewProductionEncoderConfig())
			encoder := zapcore.NewJSONEncoder(encoderConfig)
			core := zapcore.NewCore(encoder, os.Stdout, options.level.getZapLevel())
			zl = zap.New(ecszap.WrapCore(core), zap.AddCaller(), zap.AddStacktrace(zap.PanicLevel))
		} else {
			conf := zap.NewProductionConfig()
			conf.Level = options.level.getZapLevel()
			setCommonConfigParams(&conf)
			zl, err = conf.Build(zap.AddStacktrace(zap.PanicLevel))
			if err != nil {
				err = fmt.Errorf("build logger error: %w", err)
				return
			}
		}

		zapLogger = zl
	})

	return zapLogger, nil
}

func setCommonConfigParams(conf *zap.Config) {
	conf.EncoderConfig.MessageKey = "message"
	conf.EncoderConfig.TimeKey = "msg_time"
	conf.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}
	conf.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(strings.ToUpper(l.String()))
	}
}
