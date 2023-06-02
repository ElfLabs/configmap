package logger

import (
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger Zap logger config.
type Logger struct {
	zap.Config `json:",inline" yaml:",inline"`
	Rotate     *lumberjack.Logger `json:"rotate,omitempty" yaml:"rotate,omitempty"`
	Console    bool               `json:"console" yaml:"console"`
	Name       string             `json:"name,omitempty" yaml:"name,omitempty"`
}

func NewDevelopmentLogger() *Logger {
	return &Logger{
		Config:  zap.NewDevelopmentConfig(),
		Rotate:  nil,
		Console: true,
	}
}

func NewProductionLogger() *Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return &Logger{
		Config: zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		},
		Rotate: &lumberjack.Logger{
			Compress: true,
		},
		Console: true,
	}
}

func (l *Logger) newRotateWriter() io.Writer {
	if l.Rotate.Filename == "" {
		l.Rotate.Filename = filepath.Base(os.Args[0]) + ".log"
	}
	return l.Rotate
}

func (l *Logger) newRotateZapCore() zapcore.Core {
	var encoder zapcore.Encoder
	switch l.Config.Encoding {
	case "console":
		encoder = zapcore.NewConsoleEncoder(l.Config.EncoderConfig)
	case "json":
		encoder = zapcore.NewJSONEncoder(l.Config.EncoderConfig)
	}
	return zapcore.NewCore(encoder, zapcore.AddSync(l.newRotateWriter()), l.Config.Level)
}

func (l *Logger) enableStderrOutput() bool {
	for _, item := range l.Config.OutputPaths {
		if item == "stderr" {
			return true
		}
	}
	return false
}

func (l *Logger) NewZapLogger(opts ...zap.Option) (*zap.Logger, error) {
	if l.Console && !l.enableStderrOutput() {
		l.Config.OutputPaths = append(l.Config.OutputPaths, "stderr")
	}
	if l.Rotate != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(
				core,
				l.newRotateZapCore(),
			)
		}))
	}
	// build logger
	logger, err := l.Config.Build(opts...)
	if err != nil {
		return nil, err
	}
	if l.Name != "" {
		logger = logger.Named(l.Name)
	}
	return logger, nil
}
