package log

import apmv2 "go.elastic.co/apm/v2"

type apm struct {
}

func NewAPM() apmv2.Logger {
	return apm{}
}

func (a apm) Debugf(format string, args ...interface{}) {
	Printf(DebugLevel, "APM", "", format, args...)
}

func (a apm) Errorf(format string, args ...interface{}) {
	Printf(ErrorLevel, "APM", "", format, args...)
}

func (a apm) Warningf(format string, args ...interface{}) {
	Printf(WarnLevel, "APM", "", format, args...)
}
