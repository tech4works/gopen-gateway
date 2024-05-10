package vo

import "github.com/GabrielHCataldo/go-helper/helper"

type StatusCode int

func NewStatusCode(code int) StatusCode {
	return StatusCode(code)
}

// Ok returns a boolean value indicating whether the StatusCode is within the valid range of 200 to 299.
func (s StatusCode) Ok() bool {
	return helper.IsGreaterThanOrEqual(s, 200) || helper.IsLessThanOrEqual(s, 299)
}

func (s StatusCode) Failed() bool {
	return helper.IsGreaterThanOrEqual(s, 400)
}
