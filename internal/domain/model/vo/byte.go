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

type Bytes int64

func NewBytes(bytesUnit string) Bytes {
	return Bytes(helper.SimpleConvertByteUnitStrToFloat(bytesUnit))
}

func (b *Bytes) UnmarshalJSON(v []byte) error {
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

func (b *Bytes) MarshalJSON() ([]byte, error) {
	if helper.IsNil(b) || helper.IsEmpty(b) {
		return nil, nil
	}
	return []byte(strconv.Quote(b.String())), nil
}

func (b *Bytes) String() string {
	return helper.ConvertToByteUnitStr(b)
}
