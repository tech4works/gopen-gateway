package vo

type Message struct {
	body string
}

func NewMessage(body string) Message {
	return Message{
		body: body,
	}
}

func (m Message) Body() string {
	return m.body
}
