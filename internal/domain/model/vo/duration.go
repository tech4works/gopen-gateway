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
	"time"
)

type Duration time.Duration

func NewDuration(timeDuration time.Duration) Duration {
	return Duration(timeDuration)
}

func (d *Duration) Time() time.Duration {
	return time.Duration(*d)
}

func (d *Duration) String() string {
	return d.Time().String()
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	str, err := strconv.Unquote(string(b))
	if helper.IsNotNil(err) {
		return err
	} else if helper.IsEmpty(str) {
		return nil
	}

	duration, err := time.ParseDuration(str)
	if helper.IsNotNil(err) {
		return err
	}

	*d = Duration(duration)
	return nil
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	if helper.IsNil(d) || helper.IsEmpty(d) {
		return nil, nil
	}
	return []byte(strconv.Quote(d.String())), nil
}
