package logger

import (
	"os"

	"github.com/natefinch/lumberjack"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level      string
	Filename   string
	MaxSize    int
	StSkip     int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// NewLogger creates a new logger instance
func NewLogger(cfg Config) *zap.Logger {
	logLevel := getLogLevel(cfg.Level)

	// Ensure log directory exists
	if err := os.MkdirAll("./logs", 0755); err != nil {
		panic(err)
	}

	// File writer with rotation
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	// Console writer
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Core
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, logLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, logLevel),
	)

	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}
	if cfg.StSkip > 0 {
		opts = append(opts, zap.AddCallerSkip(cfg.StSkip))
	}

	logger := zap.New(core, opts...)

	return logger
}

// NewStdoutLogger creates a structured logger for OpenTelemetry/container environments
// that outputs structured JSON logs to stdout without file rotation.
func NewStdoutLogger(cfg Config) *zap.Logger {
	logLevel := getLogLevel(cfg.Level)

	stdoutWriter := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), stdoutWriter, logLevel)

	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}
	if cfg.StSkip > 0 {
		opts = append(opts, zap.AddCallerSkip(cfg.StSkip))
	}

	return zap.New(core, opts...)
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
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
