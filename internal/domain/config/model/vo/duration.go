package vo

/*
 * Copyright 2024 Gabriel Cataldo
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

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"strconv"
	"time"
)

type Duration time.Duration

// Time converts the Duration object to a time.Duration object.
func (d *Duration) Time() time.Duration {
	return time.Duration(*d)
}

// String returns the string representation of the Duration object's time.Duration value.
func (d *Duration) String() string {
	return d.Time().String()
}

// UnmarshalJSON decodes a JSON value into the Duration object.
// It removes the quotes around the value using the Unquote function.
// Then it uses ParseDuration to convert the string into a time.Duration.
// The resulting duration value is assigned to the receiver of the method.
// If the input value is empty, it returns nil.
// If there is an error during the conversion, it returns the error.
func (d *Duration) UnmarshalJSON(b []byte) error {
	// Remova as aspas em torno do valor com a função Unquote.
	str, err := strconv.Unquote(string(b))
	if helper.IsNotNil(err) {
		return err
	} else if helper.IsEmpty(str) {
		return nil
	}

	// Use ParseDuration para converter a string em um time.Duration.
	duration, err := time.ParseDuration(str)
	if helper.IsNotNil(err) {
		return err
	}

	// Atribua o valor de duração ao seu receptor.
	*d = Duration(duration)

	return nil
}

// MarshalJSON converts the Duration object to its JSON representation.
// If the Duration object is nil or empty, it returns nil, nil.
// Otherwise, it returns the byte slice representation of the Duration object's string value
// enclosed in quotes using the strconv.Quote function.
// The resulting byte slice and nil error are returned.
//
// See UnmarshalJSON for the reverse operation.
func (d *Duration) MarshalJSON() ([]byte, error) {
	if helper.IsNil(d) || helper.IsEmpty(d) {
		return nil, nil
	}
	return []byte(strconv.Quote(d.String())), nil
}
