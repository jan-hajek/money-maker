package log

import "github.com/sirupsen/logrus"

type Log interface {
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	WithField(key string, value interface{}) *logrus.Entry
}
