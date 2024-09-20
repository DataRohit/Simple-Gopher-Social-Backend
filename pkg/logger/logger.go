package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once            sync.Once
	LOG_OUTPUT_PATH = os.Getenv("LOG_OUTPUT_PATH")
)

type Logger struct {
	zapLogger *zap.Logger
}

var instance *Logger

func InitializeLogger(logLevel zapcore.Level, outputPaths []string) error {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(logLevel)
	config.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "ts",
		NameKey:      "logger",
		CallerKey:    "",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.RFC3339TimeEncoder,
		EncodeCaller: zapcore.FullCallerEncoder,
	}

	var cores []zapcore.Core

	for _, path := range outputPaths {
		var sink zapcore.WriteSyncer
		if filepath.Ext(path) == ".txt" {
			file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("failed to open log file: %v", err)
			}
			sink = zapcore.AddSync(file)
		} else if path == "stdout" {
			sink = zapcore.AddSync(os.Stdout)
		} else {
			return fmt.Errorf("unsupported output path: %s", path)
		}

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			sink,
			config.Level,
		)
		cores = append(cores, core)
	}

	var zapLogger *zap.Logger
	if len(cores) == 1 {
		zapLogger = zap.New(cores[0])
	} else {
		zapLogger = zap.New(zapcore.NewTee(cores...))
	}

	instance = &Logger{zapLogger: zapLogger}
	return nil
}

func GetLogger() *Logger {
	once.Do(func() {
		if err := InitializeLogger(zapcore.DebugLevel, []string{LOG_OUTPUT_PATH}); err != nil {
			log.Fatalf("error initializing logger: %v", err)
		}
		instance.Info("logger initialized successfully")
	})
	return instance
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
	os.Exit(1)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}
