/*
Package log provides a simple, concurrent logger for logging messages
and errors in a structured format. The logger supports various log levels,
allowing users to filter logs based on severity. It is designed to handle
concurrent log requests safely, making it suitable for use in multi-threaded
applications.

This package utilizes the Zap logging library to offer efficient and
high-performance logging. By customizing the logger's configuration, users
can control the log output format and selectively include additional
contextual information, such as the module name.

The log package is essential for monitoring application behavior,
troubleshooting issues, and maintaining an organized logging strategy
in complex systems.
*/
package log

import (
	"fmt"

	"github.com/heyrovsky/disturbdb/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogMessage represents a log message with a description and a long-form message.
type LogMessage struct {
	Desc string // Description of the log message
	Msg  string // Detailed content of the message
}

// ErrLogMessage represents a log message for errors, including the associated error.
type ErrLogMessage struct {
	Desc string // Description of the error log
	Err  error  // The error that occurred
}

// Logger represents a custom logger with a zap logger instance and a module name.
type Logger struct {
	zapLogger *zap.Logger   // The underlying zap logger instance
	level     zapcore.Level // The current logging level
	Module    string        // Name of the module for logging context
}

// NewLogger initializes a new Logger with zap, configuring the logger settings.
func NewLogger(module string) *Logger {
	// Create a new production configuration for the zap logger
	config := zap.NewProductionConfig()
	config.Level.SetLevel(setLogLevel()) // Set log level based on configuration

	// Customize the encoder configuration to remove caller information
	config.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "message",                     // Key for the log message
		LevelKey:       "level",                       // Key for the log level
		TimeKey:        "time",                        // Key for the timestamp
		NameKey:        "logger",                      // Key for the logger name
		CallerKey:      "",                            // Remove caller key to not include caller info (using module system)
		StacktraceKey:  "stacktrace",                  // Key for stack trace
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // Encoder for time format
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // Encoder for log level
		EncodeDuration: zapcore.StringDurationEncoder, // Encoder for duration
		EncodeCaller:   zapcore.FullCallerEncoder,     // Encoder for caller info (not used)
	}

	// Build the zap logger instance
	zapLogger, err := config.Build()
	if err != nil {
		fmt.Println("Failed to initialize zap logger:", err) // Log error if initialization fails
		return nil                                           // Return nil if logger initialization fails
	}

	return &Logger{
		zapLogger: zapLogger,     // Set the zap logger instance
		Module:    module,        // Set the module name
		level:     setLogLevel(), // Set the logging level
	}
}

// setLogLevel determines the appropriate log level from the configuration.
func setLogLevel() zapcore.Level {
	switch config.LOG_LEVEL {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "CRITICAL":
		return zapcore.DPanicLevel // zap does not have a direct "CRITICAL", using DPanicLevel
	case "PANIC":
		return zapcore.PanicLevel
	default:
		fmt.Println("No LOG_LEVEL loaded, defaulting to INFO") // Log message for default log level
		return zapcore.InfoLevel                               // Default log level is INFO
	}
}

// DEBUG logs a debug message, including the module name.
func (l *Logger) DEBUG(msg LogMessage) {
	if l.level <= zapcore.DebugLevel {
		l.zapLogger.Debug(msg.Desc, zap.String("message", msg.Msg), zap.String("module", l.Module))
	}
}

// INFO logs an info message, including the module name.
func (l *Logger) INFO(msg LogMessage) {
	if l.level <= zapcore.InfoLevel {
		l.zapLogger.Info(msg.Desc, zap.String("message", msg.Msg), zap.String("module", l.Module))
	}
}

// WARN logs a warning message, including the module name.
func (l *Logger) WARN(msg LogMessage) {
	if l.level <= zapcore.WarnLevel {
		l.zapLogger.Warn(msg.Desc, zap.String("message", msg.Msg), zap.String("module", l.Module))
	}
}

// ERROR logs an error message, including the error and module name.
func (l *Logger) ERROR(msg ErrLogMessage) {
	if l.level <= zapcore.ErrorLevel {
		l.zapLogger.Error(msg.Desc, zap.Error(msg.Err), zap.String("module", l.Module))
	}
}

// CRITICAL logs a critical message, including the error and module name.
func (l *Logger) CRITICAL(msg ErrLogMessage) {
	if l.level <= zapcore.DPanicLevel {
		l.zapLogger.DPanic(msg.Desc, zap.Error(msg.Err), zap.String("module", l.Module))
	}
}

// PANIC logs a panic message, including the error and module name.
func (l *Logger) PANIC(msg ErrLogMessage) {
	if l.level <= zapcore.PanicLevel {
		l.zapLogger.Panic(msg.Desc, zap.Error(msg.Err), zap.String("module", l.Module))
	}
}

// Close cleans up any resources used by the logger, ensuring all log entries are flushed.
func (l *Logger) Close() {
	_ = l.zapLogger.Sync() // Flush any buffered log entries
}
