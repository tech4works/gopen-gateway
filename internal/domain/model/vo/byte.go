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
	"github.com/tech4works/errors"
	"math"
	"regexp"
	"strconv"
)

var unitMap = map[string]float64{
	"B":  0,
	"KB": 1,
	"MB": 2,
	"GB": 3,
	"TB": 4,
	"PB": 5,
	"EB": 6,
	"ZB": 7,
	"YB": 8,
}

type Bytes int64

func NewBytes(bytesUnit string) Bytes {
	bytes, _ := NewBytesWithErr(bytesUnit)
	return bytes
}

func NewBytesWithErr(bytesUnit string) (Bytes, error) {
	regex := regexp.MustCompile(`^(\d+)\s?(\w+)?$`)
	match := regex.FindAllStringSubmatch(bytesUnit, -1)

	if len(match) == 0 || len(match[0]) < 3 {
		return 0, errors.Newf("Error byte unit mal formated")
	}

	qty, _ := strconv.ParseFloat(match[0][1], 64)
	unit := match[0][2]

	if exp, ok := unitMap[unit]; ok {
		return Bytes(qty * math.Pow(1024, exp)), nil
	}

	return 0, errors.Newf("Error byte unit mal formated")
}

func NewBytesByInt(i int) Bytes {
	return Bytes(i)
}

func (b *Bytes) UnmarshalJSON(v []byte) error {
	str, err := strconv.Unquote(string(v))
	if checker.NonNil(err) {
		return err
	} else if checker.IsEmpty(str) {
		return nil
	}

	value, err := NewBytesWithErr(str)
	if checker.NonNil(err) {
		return err
	}

	*b = value

	return nil
}

func (b *Bytes) MarshalJSON() ([]byte, error) {
	if checker.IsNil(b) || checker.IsEmpty(b) {
		return nil, nil
	}

	return []byte(strconv.Quote(b.String())), nil
}

func (b Bytes) String() string {
	return fmt.Sprintf("%vB", int(b))
}
