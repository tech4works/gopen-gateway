package publisher

type Response struct {
	OK   bool  `json:"ok,omitempty"`
	Body *Body `json:"body,omitempty"`
}

type Body struct {
	Path             string `json:"path,omitempty"`
	Provider         string `json:"provider,omitempty"`
	MessageID        string `json:"messageId,omitempty"`
	SequentialNumber string `json:"sequentialNumber,omitempty"`
}
