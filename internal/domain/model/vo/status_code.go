package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"net/http"
)

type StatusCode struct {
	code        int
	description string
}

func NewStatusCode(code int) StatusCode {
	return StatusCode{
		code:        code,
		description: http.StatusText(code),
	}
}

func (s *StatusCode) OK() bool {
	return helper.IsGreaterThanOrEqual(s.Code(), 200) || helper.IsLessThanOrEqual(s.Code(), 299)
}

func (s *StatusCode) Failed() bool {
	return helper.IsGreaterThanOrEqual(s.Code(), 400)
}

func (s *StatusCode) Code() int {
	return s.code
}

func (s *StatusCode) Description() string {
	return s.description
}

func (s *StatusCode) MarshalJSON() ([]byte, error) {
	return helper.ConvertToBytes(s.Code())
}

func (s *StatusCode) UnmarshalJSON(data []byte) error {
	code, err := helper.ConvertToInt(data)
	if helper.IsNotNil(err) {
		return err
	}
	s.code = code
	s.description = http.StatusText(code)
	return nil
}

func (s *StatusCode) String() string {
	return fmt.Sprintf("%v %s", s.Code(), s.Description())
}
