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

package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"strconv"
)

// Bytes represents a type for storing byte values.
type Bytes int64

// NewBytes creates a new Bytes value based on the given byte unit string.
// It converts the byte unit string to a float value using the helper.SimpleConvertByteUnitStrToFloat function.
func NewBytes(bytesUnit string) Bytes {
	return Bytes(helper.SimpleConvertByteUnitStrToFloat(bytesUnit))
}

// UnmarshalJSON unmarshals a JSON-encoded byte value into a Bytes value.
// It takes a byte slice v as input, converts it to a string s,
// and then checks if the string is not empty.
// If it's not empty, it converts the string to a float value using the helper.ConvertByteUnitStrToFloat function,
// and assigns the converted value to the pointer receiver b.
// Finally, it returns nil.
func (b *Bytes) UnmarshalJSON(v []byte) error {
	// Remova as aspas em torno do valor com a função Unquote.
	str, err := strconv.Unquote(string(v))
	if helper.IsNotNil(err) {
		return err
	} else if helper.IsEmpty(str) {
		return nil
	}

	value, err := helper.ConvertByteUnitStrToFloat(str)
	if helper.IsNotNil(err) {
		return err
	}
	*b = Bytes(value)

	return nil
}

// MarshalJSON marshals a Bytes value into a JSON-encoded byte slice.
// It checks if the pointer receiver b is nil or empty. If it is,
// it returns nil and nil error.
// Otherwise, it converts the Bytes value to a string using the b.String() method,
// converts the string to a byte slice using []byte(),
// and returns the byte slice and nil error.
func (b *Bytes) MarshalJSON() ([]byte, error) {
	if helper.IsNil(b) || helper.IsEmpty(b) {
		return nil, nil
	}
	return []byte(strconv.Quote(b.String())), nil
}

// String returns a string representation of the Bytes value.
// It uses the helper.ConvertToByteUnitStr function to format the value as a byte unit string.
func (b *Bytes) String() string {
	return helper.ConvertToByteUnitStr(b)
}
