package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// InitLogger initializes the global logger
func InitLogger() {
	if logger != nil {
		panic("Logger is already initialized")
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Colored levels
	config.EncoderConfig.TimeKey = "timestamp"                          // Time key
	config.EncoderConfig.CallerKey = "caller"                           // File and line key
	config.EncoderConfig.MessageKey = "message"                         // Message key
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // Time format
	config.DisableStacktrace = true

	baseLogger, err := config.Build()
	if err != nil {
		panic(err)
	}

	// Skip 1 call frame for correct log source display
	logger = baseLogger.WithOptions(zap.AddCallerSkip(1))
}

// GetLogger returns the current logger instance
func GetLogger() *zap.Logger {
	if logger == nil {
		panic("Logger is not initialized. Call InitLogger() first.")
	}
	return logger
}

// SyncLogger synchronizes the logger buffer (e.g., flush)
func SyncLogger() {
	if logger != nil {
		_ = logger.Sync()
	}
}

// Debug logs a debug-level message
func Debug(msg string, fields ...zap.Field) {
	ensureLoggerInitialized()
	logger.Debug(msg, fields...)
}

// Info logs an info-level message
func Info(msg string, fields ...zap.Field) {
	ensureLoggerInitialized()
	logger.Info(msg, fields...)
}

// Warn logs a warning-level message
func Warn(msg string, fields ...zap.Field) {
	ensureLoggerInitialized()
	logger.Warn(msg, fields...)
}

// Error logs an error-level message
func Error(msg string, fields ...zap.Field) {
	ensureLoggerInitialized()
	logger.Error(msg, fields...)
}

// Fatal logs a fatal-level message and exits
func Fatal(msg string, fields ...zap.Field) {
	ensureLoggerInitialized()
	logger.Fatal(msg, fields...)
}

// ensureLoggerInitialized checks if the logger is initialized
func ensureLoggerInitialized() {
	if logger == nil {
		panic("Logger is not initialized. Call InitLogger() first.")
	}
}
