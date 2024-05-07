package vo

import (
	"bytes"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/model/enum"
)

// CacheBody represents the caching value of an HTTP httpResponse body.
type CacheBody struct {
	// ContentType represents the format of the content.
	ContentType enum.ContentType `json:"content-type,omitempty"`
	// Value represents the caching content of an HTTP httpResponse body.
	// It is a pointer to the CacheBodyValue type, which is an alias for bytes.Buffer.
	// The value is nullable and is omitted in JSON if it is empty.
	Value *CacheBodyValue `json:"value,omitempty"`
}

// CacheBodyValue is an alias for bytes.Buffer type used to represent the caching value
// of an HTTP httpResponse body. It contains methods to convert the value to different
// representations, such as string and JSON.
type CacheBodyValue bytes.Buffer

// newCacheBody creates a new instance of CacheBody based on the provided body.
// If the body is nil, it returns nil.
// Otherwise, it sets the ContentType field of CacheBody based on the ContentType method of body.
// It sets the Value field of CacheBody by calling newCacheBodyValue with the Value method of body as an argument.
// It returns a pointer to the constructed CacheBody instance.
func newCacheBody(body *Body) *CacheBody {
	if helper.IsNil(body) {
		return nil
	}
	return &CacheBody{
		ContentType: body.ContentType(),
		Value:       newCacheBodyValue(body.Value()),
	}
}

func newCacheBodyValue(buffer *bytes.Buffer) *CacheBodyValue {
	if helper.IsNil(buffer) || helper.IsEmpty(buffer.Bytes()) {
		return nil
	}
	return (*CacheBodyValue)(buffer)
}

// String returns the string representation of the CacheBodyValue instance.
// It calls the String method of the underlying bytes.Buffer type to get the string representation.
func (c *CacheBodyValue) String() string {
	return (*bytes.Buffer)(c).String()
}

// Bytes returns the byte slice representation of the CacheBodyValue instance.
// It calls the Bytes method of the underlying bytes.Buffer type to get the byte slice representation.
func (c *CacheBodyValue) Bytes() []byte {
	return (*bytes.Buffer)(c).Bytes()
}

// MarshalJSON returns the JSON encoding of the CacheBodyValue instance.
// The JSON encoding is obtained by calling the String method of the
// underlying bytes.Buffer type to retrieve the string representation,
// and then encoding it using json.Marshal.
// It returns a byte slice representing the JSON encoding and an error,
// if any occurred during the encoding process.
func (c *CacheBodyValue) MarshalJSON() ([]byte, error) {
	return c.Bytes(), nil
}

// UnmarshalJSON decodes the JSON data into a string and writes
// the string to the underlying bytes.Buffer type.
// It returns an error if there is an issue with decoding or writing
// the string to the buffer.
func (c *CacheBodyValue) UnmarshalJSON(data []byte) error {
	_, err := (*bytes.Buffer)(c).Write(data)
	return err
}
