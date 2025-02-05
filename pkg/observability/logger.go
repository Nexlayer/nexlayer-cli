// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fatih/color"
)

// LogLevel represents the severity of a log entry.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger provides structured logging capabilities.
type Logger struct {
	level    LogLevel
	writer   io.Writer
	jsonMode bool
	maxSize  int64      // maximum size of log file in bytes
	maxAge   int        // maximum number of days to retain old log files
	rotate   *time.Time // next rotation time
}

// LoggerOption allows configuring the logger with functional options.
type LoggerOption func(*Logger)

// WithJSON enables JSON-formatted logging.
func WithJSON() LoggerOption {
	return func(l *Logger) {
		l.jsonMode = true
	}
}

// WithRotation enables log rotation with specified max size and age.
func WithRotation(maxSizeMB int64, maxAgeDays int) LoggerOption {
	return func(l *Logger) {
		l.maxSize = maxSizeMB * 1024 * 1024
		l.maxAge = maxAgeDays
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		l.rotate = &next
	}
}

// NewLogger creates a new logger instance.
func NewLogger(level LogLevel, opts ...LoggerOption) *Logger {
	// Create logs directory if it doesn't exist.
	logDir := filepath.Join(os.Getenv("HOME"), ".nexlayer", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Warning: Could not create log directory: %v\n", err)
		return &Logger{level: level, writer: os.Stdout}
	}

	// Open log file.
	logFile := filepath.Join(logDir, fmt.Sprintf("nexlayer-%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Warning: Could not open log file: %v\n", err)
		return &Logger{level: level, writer: os.Stdout}
	}

	logger := &Logger{
		level:   level,
		writer:  io.MultiWriter(f, os.Stdout),
		maxSize: 50 * 1024 * 1024, // default 50MB
		maxAge:  7,                // default 7 days
	}

	// Apply options.
	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

func (l *Logger) log(ctx context.Context, level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	// Rotate log file if needed.
	if l.rotate != nil && time.Now().After(*l.rotate) {
		l.rotateLogFile()
	}

	// Retrieve caller information.
	_, file, line, _ := runtime.Caller(2)
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(msg, args...)

	var traceID, requestID, sessionID string
	if ctx != nil {
		if id, ok := ctx.Value("trace_id").(string); ok {
			traceID = id
		}
		if id, ok := ctx.Value("request_id").(string); ok {
			requestID = id
		}
		if id, ok := ctx.Value("session_id").(string); ok {
			sessionID = id
		}
	}

	if l.jsonMode {
		entry := map[string]interface{}{
			"timestamp": timestamp,
			"level":     level.String(),
			"file":      filepath.Base(file),
			"line":      line,
			"message":   message,
		}
		if traceID != "" {
			entry["trace_id"] = traceID
		}
		if requestID != "" {
			entry["request_id"] = requestID
		}
		if sessionID != "" {
			entry["session_id"] = sessionID
		}

		jsonBytes, err := json.Marshal(entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling log entry: %v\n", err)
			return
		}
		fmt.Fprintf(l.writer, "%s\n", jsonBytes)
	} else {
		var levelStr string
		switch level {
		case DEBUG:
			levelStr = color.BlueString("DEBUG")
		case INFO:
			levelStr = color.GreenString("INFO")
		case WARN:
			levelStr = color.YellowString("WARN")
		case ERROR:
			levelStr = color.RedString("ERROR")
		}

		logEntry := fmt.Sprintf("[%s] %s %s:%d", timestamp, levelStr, filepath.Base(file), line)
		if traceID != "" {
			logEntry += fmt.Sprintf(" [trace:%s]", traceID)
		}
		if requestID != "" {
			logEntry += fmt.Sprintf(" [req:%s]", requestID)
		}
		if sessionID != "" {
			logEntry += fmt.Sprintf(" [sess:%s]", sessionID)
		}
		logEntry += fmt.Sprintf(" %s\n", message)

		fmt.Fprint(l.writer, logEntry)
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.log(ctx, DEBUG, msg, args...)
}

func (l *Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.log(ctx, INFO, msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.log(ctx, WARN, msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.log(ctx, ERROR, msg, args...)
}

func (level LogLevel) String() string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (l *Logger) rotateLogFile() {
	if fw, ok := l.writer.(*os.File); ok {
		info, err := fw.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting log file info: %v\n", err)
			return
		}

		if info.Size() >= l.maxSize {
			fw.Close()

			logDir := filepath.Dir(fw.Name())
			archiveDir := filepath.Join(logDir, "archive")
			if err := os.MkdirAll(archiveDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating archive directory: %v\n", err)
				return
			}

			timestamp := time.Now().Format("2006-01-02-15-04-05")
			archivePath := filepath.Join(archiveDir, fmt.Sprintf("nexlayer-%s.log", timestamp))
			if err := os.Rename(fw.Name(), archivePath); err != nil {
				fmt.Fprintf(os.Stderr, "Error archiving log file: %v\n", err)
				return
			}

			newFile, err := os.OpenFile(fw.Name(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating new log file: %v\n", err)
				return
			}

			l.writer = io.MultiWriter(newFile, os.Stdout)
			l.cleanOldArchives(archiveDir)
		}
	}

	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	l.rotate = &next
}

func (l *Logger) cleanOldArchives(archiveDir string) {
	cutoff := time.Now().AddDate(0, 0, -l.maxAge)

	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading archive directory: %v\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				path := filepath.Join(archiveDir, entry.Name())
				if err := os.Remove(path); err != nil {
					fmt.Fprintf(os.Stderr, "Error removing old log file %s: %v\n", path, err)
				}
			}
		}
	}
}
