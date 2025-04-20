package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger represents the logger instance
var Logger *logrus.Logger

// Setup initializes the logger with appropriate settings
func Setup(env string) {
	Logger = logrus.New()

	// Set logger output to stdout
	Logger.SetOutput(os.Stdout)

	// Set formatter based on environment
	if env == "production" {
		// Use JSON formatter for production for easier parsing
		Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		// Use text formatter for development for better readability
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			ForceColors:     true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set logger level based on environment
	if env == "development" {
		Logger.SetLevel(logrus.DebugLevel)
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}
}

// Info logs an info level message
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Debug logs a debug level message
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Error logs an error level message
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Warn logs a warning level message
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Fatal logs a fatal level message and exits
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// WithField adds a field to the log entry
func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}

// WithFields adds multiple fields to the log entry
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}
