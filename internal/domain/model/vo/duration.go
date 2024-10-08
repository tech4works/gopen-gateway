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
	"github.com/tech4works/checker"
	"strconv"
	"time"
)

type Duration time.Duration

func NewDuration(timeDuration time.Duration) Duration {
	return Duration(timeDuration)
}

func (d Duration) Time() time.Duration {
	return time.Duration(d)
}

func (d Duration) String() string {
	return d.Time().String()
}

func (d Duration) MarshalJSON() ([]byte, error) {
	if checker.IsNil(d) || checker.IsEmpty(d) {
		return nil, nil
	}
	return []byte(strconv.Quote(d.String())), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	str, err := strconv.Unquote(string(b))
	if checker.NonNil(err) {
		return err
	} else if checker.IsEmpty(str) {
		return nil
	}

	duration, err := time.ParseDuration(str)
	if checker.NonNil(err) {
		return err
	}

	*d = Duration(duration)
	return nil
}
