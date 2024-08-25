/*
 * Copyright 2024 Tech4Works
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vo

import (
	"fmt"
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
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

func (s StatusCode) OK() bool {
	return checker.IsGreaterThanOrEqual(s.Code(), 200) || checker.IsLessThanOrEqual(s.Code(), 299)
}

func (s StatusCode) Failed() bool {
	return checker.IsGreaterThanOrEqual(s.Code(), 400)
}

func (s StatusCode) Code() int {
	return s.code
}

func (s StatusCode) Description() string {
	return s.description
}

func (s StatusCode) MarshalJSON() ([]byte, error) {
	return converter.ToBytesWithErr(s.Code())
}

func (s *StatusCode) UnmarshalJSON(data []byte) error {
	code, err := converter.ToIntWithErr(data)
	if checker.NonNil(err) {
		return err
	}

	s.code = code
	s.description = http.StatusText(code)

	return nil
}

func (s StatusCode) String() string {
	return fmt.Sprintf("%v %s", s.Code(), s.Description())
}
