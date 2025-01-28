package observability

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fatih/color"
)

// LogLevel represents the severity of a log entry
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger provides structured logging capabilities
type Logger struct {
	level  LogLevel
	writer io.Writer
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel) *Logger {
	// Create logs directory if it doesn't exist
	logDir := filepath.Join(os.Getenv("HOME"), ".nexlayer", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Warning: Could not create log directory: %v\n", err)
		return &Logger{level: level, writer: os.Stdout}
	}

	// Open log file
	logFile := filepath.Join(logDir, fmt.Sprintf("nexlayer-%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Warning: Could not open log file: %v\n", err)
		return &Logger{level: level, writer: os.Stdout}
	}

	// Use MultiWriter to write to both file and stdout
	return &Logger{
		level:  level,
		writer: io.MultiWriter(f, os.Stdout),
	}
}

func (l *Logger) log(ctx context.Context, level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	// Get caller information
	_, file, line, _ := runtime.Caller(2)

	// Format timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// Format message
	message := fmt.Sprintf(msg, args...)

	// Get trace ID from context if available
	var traceID string
	if ctx != nil {
		if id, ok := ctx.Value("trace_id").(string); ok {
			traceID = id
		}
	}

	// Format log entry
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
		logEntry += fmt.Sprintf(" [%s]", traceID)
	}
	logEntry += fmt.Sprintf(" %s\n", message)

	fmt.Fprint(l.writer, logEntry)
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
