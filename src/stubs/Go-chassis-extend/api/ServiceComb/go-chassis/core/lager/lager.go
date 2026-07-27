package lager

type logger struct{}

func (l *logger) Infof(format string, args ...interface{})  {}
func (l *logger) Warnf(err error, format string, args ...interface{}) {}
func (l *logger) Debugf(format string, args ...interface{}) {}
func (l *logger) Errorf(err error, format string, args ...interface{}) {}
func (l *logger) Fatalf(err error, format string, args ...interface{}) {}

var Logger = &logger{}