package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	// Environment: "development", "production", etc.
	Environment string
	// Debug enables/disables debug-level logging
	Debug bool
	// JSON enables JSON formatting (typically used in production)
	JSON bool
	// Level sets the minimum log level
	Level string
	// Output file path (empty for stdout)
	OutputPath string
}

// Logger interface for application logging
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
	ZapLogger() *zap.Logger
}

// loggerImpl implements the Logger interface with zap
type loggerImpl struct {
	sugar *zap.SugaredLogger
	zap   *zap.Logger
}

const (
	defaultLogLevel = "info"
)

func NewLogger(debug bool) Logger {
	config := LogConfig{
		Environment: "development",
		Debug:       debug,
		JSON:        false,
		Level:       defaultLogLevel,
	}

	// In production, use JSON format
	if os.Getenv("APP_ENV") == "production" {
		config.Environment = "production"
		config.JSON = true
	}

	return NewLoggerWithConfig(config)
}

// NewLoggerWithConfig creates a logger with specific configuration
func NewLoggerWithConfig(config LogConfig) Logger {
	var level zapcore.Level
	switch config.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Override level if debug mode is enabled
	if config.Debug {
		level = zapcore.DebugLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// In development mode, use colorful output for console
	if config.Environment == "development" && !config.JSON {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Define output paths
	outputPaths := []string{"stdout"}
	if config.OutputPath != "" {
		outputPaths = append(outputPaths, config.OutputPath)
	}

	// Create zap configuration
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      config.Environment == "development",
		Encoding:         "console", // Default for development
		EncoderConfig:    encoderConfig,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}

	// Use JSON encoding for production
	if config.JSON {
		zapConfig.Encoding = "json"
	}

	// Build the logger
	zapLogger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		fallback := zap.NewExample()
		fallback.Error("Failed to create logger", zap.Error(err))
		return &loggerImpl{
			sugar: fallback.Sugar(),
			zap:   fallback,
		}
	}

	return &loggerImpl{
		sugar: zapLogger.Sugar(),
		zap:   zapLogger,
	}
}

// Debug logs a message at debug level
func (l *loggerImpl) Debug(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

// Info logs a message at info level
func (l *loggerImpl) Info(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Warn logs a message at warn level
func (l *loggerImpl) Warn(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

// Error logs a message at error level
func (l *loggerImpl) Error(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func (l *loggerImpl) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

// With creates a child logger with the provided context fields
func (l *loggerImpl) With(keysAndValues ...interface{}) Logger {
	return &loggerImpl{
		sugar: l.sugar.With(keysAndValues...),
		zap:   l.zap,
	}
}

// ZapLogger provides direct access to the underlying zap logger
func (l *loggerImpl) ZapLogger() *zap.Logger {
	return l.zap
}
