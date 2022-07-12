package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Logger interface {
	FailOnError(err error, args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
}

type logger struct {
	*logrus.Logger
	activeDate time.Time
	logfile    *os.File
	logPath    LogPath
	config     LoggerConfig
}

type LoggerConfig struct {
	BasicFields logrus.Fields
}

type LogPath string

// NewLogger initializes the logger wrapper
func NewLogger(logPath LogPath, config LoggerConfig) *logger {
	var baseLogger = logrus.New()

	var logger = &logger{
		baseLogger,
		time.Now().AddDate(0, 0, -1),
		nil,
		logPath,
		config,
	}

	logger.Formatter = &logrus.JSONFormatter{}

	return logger
}

// Rotate the log each day
func (l *logger) rotateIfNeeded() {
	currentTime := time.Now()
	if l.activeDate.IsZero() || currentTime.Format("2006-01-02") != l.activeDate.Format("2006-01-02") {
		if l.logfile != nil {
			l.logfile.Close()
			os.Remove(fmt.Sprintf("%slog-%s.log", string(l.logPath), l.activeDate.Format("2006-01-02")))
		}
		l.activeDate = currentTime
		l.logfile, _ = os.OpenFile(fmt.Sprintf("%s/log-%s.log", string(l.logPath), l.activeDate.Format("2006-01-02")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		l.SetOutput(l.logfile)
	}
}

// Defines basic fields to log
func (l *logger) basicFields() *logrus.Entry {
	return l.WithFields(
		l.config.BasicFields,
	)
}

func (l *logger) FailOnError(err error, args ...interface{}) {
	l.rotateIfNeeded()
	if err != nil {
		l.Fatalf(err.Error(), args)
	}
}

func (l *logger) Info(args ...interface{}) {
	l.rotateIfNeeded()
	l.basicFields().Info(args)
}

func (l *logger) Warning(args ...interface{}) {
	l.rotateIfNeeded()
	l.basicFields().Warning(args)
}

func (l *logger) Error(args ...interface{}) {
	l.rotateIfNeeded()
	l.basicFields().Error(args)
}
