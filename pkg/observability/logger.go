// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
// Package observability provides logging and metrics for the Nexlayer CLI.
package observability

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel represents the logging level
type LogLevel string

const (
	// Log levels
	DEBUG LogLevel = "debug"
	INFO  LogLevel = "info"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
)

// LogOption represents logger configuration options
type LogOption func(*Logger)

// Logger wraps zap logger with CLI-specific methods
type Logger struct {
	*zap.SugaredLogger
	level LogLevel
	opts  []LogOption
}

// WithJSON enables JSON logging format
func WithJSON() LogOption {
	return func(l *Logger) {
		l.opts = append(l.opts, func(l *Logger) {
			config := zap.NewProductionConfig()
			logger, _ := config.Build()
			l.SugaredLogger = logger.Sugar()
		})
	}
}

// WithRotation enables log rotation with specified size and days
func WithRotation(maxSizeMB int, maxDays int) LogOption {
	return func(l *Logger) {
		l.opts = append(l.opts, func(l *Logger) {
			// Configure log rotation
		})
	}
}

// WithLevel sets the logging level
func WithLevel(level LogLevel) LogOption {
	return func(l *Logger) {
		l.level = level
	}
}

// NewLogger creates a new CLI logger
func NewLogger(level LogLevel, opts ...LogOption) *Logger {
	logger := &Logger{
		level: level,
	}

	// Apply options
	for _, opt := range opts {
		opt(logger)
	}

	// Create default logger if none configured
	if logger.SugaredLogger == nil {
		config := zap.NewDevelopmentConfig()
		switch level {
		case DEBUG:
			config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		case INFO:
			config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		case WARN:
			config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		case ERROR:
			config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		}

		baseLogger, _ := config.Build()
		logger.SugaredLogger = baseLogger.Sugar()
	}

	return logger
}

// Debug logs a debug message with context
func (l *Logger) Debug(ctx context.Context, format string, args ...interface{}) {
	l.SugaredLogger.Debugw(fmt.Sprintf(format, args...),
		"timestamp", time.Now(),
		"context", getContextFields(ctx),
	)
}

// Info logs an info message with context
func (l *Logger) Info(ctx context.Context, format string, args ...interface{}) {
	l.SugaredLogger.Infow(fmt.Sprintf(format, args...),
		"timestamp", time.Now(),
		"context", getContextFields(ctx),
	)
}

// Warn logs a warning message with context
func (l *Logger) Warn(ctx context.Context, format string, args ...interface{}) {
	l.SugaredLogger.Warnw(fmt.Sprintf(format, args...),
		"timestamp", time.Now(),
		"context", getContextFields(ctx),
	)
}

// Error logs an error message with context
func (l *Logger) Error(ctx context.Context, format string, args ...interface{}) {
	l.SugaredLogger.Errorw(fmt.Sprintf(format, args...),
		"timestamp", time.Now(),
		"context", getContextFields(ctx),
	)
}

// Fatal logs a fatal message with context and exits
func (l *Logger) Fatal(ctx context.Context, format string, args ...interface{}) {
	l.SugaredLogger.Fatalw(fmt.Sprintf(format, args...),
		"timestamp", time.Now(),
		"context", getContextFields(ctx),
	)
}

// getContextFields extracts fields from context
func getContextFields(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})

	// Add standard fields
	fields["pid"] = os.Getpid()
	fields["command"] = filepath.Base(os.Args[0])

	// Add custom fields from context
	if ctx != nil {
		if requestID := ctx.Value("request_id"); requestID != nil {
			fields["request_id"] = requestID
		}
		if userID := ctx.Value("user_id"); userID != nil {
			fields["user_id"] = userID
		}
		if command := ctx.Value("command"); command != nil {
			fields["command"] = command
		}
	}

	return fields
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(key, value),
		level:         l.level,
		opts:          l.opts,
	}
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(args...),
		level:         l.level,
		opts:          l.opts,
	}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}
