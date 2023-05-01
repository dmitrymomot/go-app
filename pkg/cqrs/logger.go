package cqrs

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/sirupsen/logrus"
)

// Logrus wrapper implements cqrs.Logger interface
type logrusWrapper struct {
	log *logrus.Entry
}

// NewLogrusWrapper returns new instance of logrusWrapper
func NewLogrusWrapper(log *logrus.Entry) Logger {
	return &logrusWrapper{log: log}
}

// Error logs error message
func (l *logrusWrapper) Error(msg string, err error, fields watermill.LogFields) {
	l.log.WithFields(logrus.Fields(fields)).WithError(err).Error(msg)
}

// Info logs info message
func (l *logrusWrapper) Info(msg string, fields watermill.LogFields) {
	l.log.WithFields(logrus.Fields(fields)).Info(msg)
}

// Debug logs debug message
func (l *logrusWrapper) Debug(msg string, fields watermill.LogFields) {
	l.log.WithFields(logrus.Fields(fields)).Debug(msg)
}

// Trace logs trace message
func (l *logrusWrapper) Trace(msg string, fields watermill.LogFields) {
	l.log.WithFields(logrus.Fields(fields)).Trace(msg)
}

// With returns new instance of logrusWrapper with additional fields
func (l *logrusWrapper) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &logrusWrapper{log: l.log.WithFields(logrus.Fields(fields))}
}
