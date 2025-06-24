package vo

import "github.com/tech4works/checker"

type Message struct {
	body            string
	groupID         *string
	deduplicationID *string
	delay           Duration
}

func NewMessage(body, groupID, deduplicationID string, delay Duration) *Message {
	return &Message{
		body:            body,
		groupID:         checker.IfEmptyReturns(&groupID, nil),
		deduplicationID: checker.IfEmptyReturns(&deduplicationID, nil),
		delay:           delay,
	}
}

func (m Message) Body() string {
	return m.body
}

func (m Message) GroupID() *string {
	return m.groupID
}

func (m Message) DeduplicationID() *string {
	return m.deduplicationID
}

func (m Message) Delay() Duration {
	return m.delay
}
