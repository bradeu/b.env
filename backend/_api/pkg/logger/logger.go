package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Log levels
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[int]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

type CustomLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	logFile     *os.File
}

var Logger *CustomLogger

// InitLogger initializes the logger with file and console output
func InitLogger(logFilePath string) error {
	// Create logs directory if it doesn't exist
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Create multi-writer for both file and console
	multiWriter := io.MultiWriter(file, os.Stdout)

	// Initialize logger with different levels
	Logger = &CustomLogger{
		debugLogger: log.New(multiWriter, "", 0),
		infoLogger:  log.New(multiWriter, "", 0),
		warnLogger:  log.New(multiWriter, "", 0),
		errorLogger: log.New(multiWriter, "", 0),
		fatalLogger: log.New(multiWriter, "", 0),
		logFile:     file,
	}

	return nil
}

// Close closes the log file
func (l *CustomLogger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}

// formatLog formats the log message with timestamp, level, and caller information
func formatLog(level int, format string, args ...interface{}) string {
	// Get caller information
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)

	// Format timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// Format message
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}

	return fmt.Sprintf("[%s] [%s] [%s:%d] %s",
		timestamp,
		levelNames[level],
		file,
		line,
		msg,
	)
}

// Debug logs a debug message
func (l *CustomLogger) Debug(format string, args ...interface{}) {
	l.debugLogger.Println(formatLog(DEBUG, format, args...))
}

// Info logs an info message
func (l *CustomLogger) Info(format string, args ...interface{}) {
	l.infoLogger.Println(formatLog(INFO, format, args...))
}

// Warn logs a warning message
func (l *CustomLogger) Warn(format string, args ...interface{}) {
	l.warnLogger.Println(formatLog(WARN, format, args...))
}

// Error logs an error message
func (l *CustomLogger) Error(format string, args ...interface{}) {
	l.errorLogger.Println(formatLog(ERROR, format, args...))
}

// Fatal logs a fatal message and exits the program
func (l *CustomLogger) Fatal(format string, args ...interface{}) {
	l.fatalLogger.Println(formatLog(FATAL, format, args...))
	os.Exit(1)
}

// Convenience functions for the global logger instance
func Debug(format string, args ...interface{}) {
	Logger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	Logger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	Logger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	Logger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	Logger.Fatal(format, args...)
}
