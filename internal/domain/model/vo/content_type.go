package vo

import "strings"

type ContentType string

func NewContentType(s string) ContentType {
	return ContentType(s)
}

func NewContentTypeTextPlain() ContentType {
	return "text/plain; charset=UTF-8"
}

func NewContentTypeJson() ContentType {
	return "application/json; charset=UTF-8"
}

func NewContentTypeXml() ContentType {
	return "application/xml; charset=UTF-8"
}

func (c ContentType) String() string {
	return string(c)
}

func (c ContentType) IsJson() bool {
	return strings.HasPrefix(string(c), "application/json")
}

func (c ContentType) IsNotJson() bool {
	return !c.IsJson()
}

func (c ContentType) IsXml() bool {
	return strings.HasPrefix(string(c), "application/xml")
}

func (c ContentType) IsText() bool {
	return strings.HasPrefix(string(c), "text/plain")
}
