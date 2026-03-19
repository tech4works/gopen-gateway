package vo

type EndpointExecution struct {
	allExecuted  bool
	allOK        bool
	degradations []BackendDegradation
}

func NewEmptyEndpointExecution() EndpointExecution {
	return EndpointExecution{}
}

func NewEndpointExecution(allExecuted, allOK bool, degradations []BackendDegradation) EndpointExecution {
	return EndpointExecution{
		allExecuted:  allExecuted,
		allOK:        allOK,
		degradations: degradations,
	}
}

func (e EndpointExecution) AllExecuted() bool {
	return e.allExecuted
}

func (e EndpointExecution) AllOK() bool {
	return e.allOK
}

func (e EndpointExecution) Degradations() []BackendDegradation {
	return e.degradations
}
