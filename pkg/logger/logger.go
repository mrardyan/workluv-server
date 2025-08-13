package logger

import (
	"fmt"
	"os"
	"time"
	"workluv/pkg/config"
)

// Logger represents a structured logger
type Logger struct {
	level  string
	output string
}

// New creates a new logger instance
func New(config config.LogConfig) *Logger {
	return &Logger{
		level:  config.Level,
		output: "stdout",
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.log("INFO", msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.log("ERROR", msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.log("WARN", msg, fields...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.log("DEBUG", msg, fields...)
}

// log logs a message with the given level
func (l *Logger) log(level, msg string, fields ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)

	// Simple structured logging format
	logEntry := fmt.Sprintf("[%s] %s: %s", timestamp, level, msg)

	// Add fields if provided
	if len(fields) > 0 {
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				logEntry += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
			}
		}
	}

	// Output to stdout for now (can be extended to support different outputs)
	fmt.Fprintln(os.Stdout, logEntry)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.log("FATAL", msg, fields...)
	os.Exit(1)
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	// For now, return the same logger
	// In a more sophisticated implementation, this would create a new logger with embedded fields
	return l
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level string) {
	l.level = level
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() string {
	return l.level
}
