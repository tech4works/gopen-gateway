package vo

import "github.com/GabrielHCataldo/go-helper/helper"

type ContentEncoding string

func NewContentEncoding(s string) ContentEncoding {
	return ContentEncoding(s)
}

func NewContentEncodingGzip() ContentEncoding {
	return "gzip"
}

func NewContentEncodingDeflate() ContentEncoding {
	return "deflate"
}

func (c ContentEncoding) String() string {
	return string(c)
}

func (c ContentEncoding) IsSupported() bool {
	return c.IsGzip() || c.IsDeflate()
}

func (c ContentEncoding) Valid() bool {
	return helper.IsNotEmpty(c)
}

func (c ContentEncoding) Invalid() bool {
	return !c.Valid()
}

func (c ContentEncoding) IsGzip() bool {
	return helper.EqualsIgnoreCase(c, "gzip")
}

func (c ContentEncoding) IsDeflate() bool {
	return helper.EqualsIgnoreCase(c, "deflate")
}
