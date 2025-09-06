package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	// Set log level based on environment
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	// Set formatter based on environment
	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// Set output to stdout
	Log.SetOutput(os.Stdout)
}

// GetLogger returns the configured logger instance
func GetLogger() *logrus.Logger {
	return Log
}

// WithField creates a new logger entry with a field
func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

// WithFields creates a new logger entry with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}

// Debug logs a message at debug level
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Info logs a message at info level
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Warn logs a message at warn level
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Error logs a message at error level
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Fatal logs a message at fatal level and exits
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// Panic logs a message at panic level and panics
func Panic(args ...interface{}) {
	Log.Panic(args...)
}

// Debugf logs a formatted message at debug level
func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

// Infof logs a formatted message at info level
func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

// Warnf logs a formatted message at warn level
func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Errorf logs a formatted message at error level
func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

// Fatalf logs a formatted message at fatal level and exits
func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

// Panicf logs a formatted message at panic level and panics
func Panicf(format string, args ...interface{}) {
	Log.Panicf(format, args...)
}
