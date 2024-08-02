package log

type noop struct {
}

func NewNoop() noop {
	return noop{}
}

func (n noop) Error(_ string) {
}

func (n noop) Infof(_ string, _ ...interface{}) {
}
