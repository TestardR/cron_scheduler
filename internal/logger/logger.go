package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -package=mock -source=logger.go -destination=$MOCK_FOLDER/logger.go Logger

// Logger is the interface to interact with logger.
type Logger interface {
	Info(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}

// New function iinitializes logger service.
func New(appName string) Logger {
	var log = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Level:     logrus.InfoLevel,
	}

	entry := log.WithFields(logrus.Fields{
		"appname": appName,
	})

	return entry
}
