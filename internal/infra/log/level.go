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

package log

import (
	"fmt"
)

type level string

const (
	InfoLevel  level = "INF"
	DebugLevel level = "DBG"
	WarnLevel  level = "WRN"
	ErrorLevel level = "ERR"
)

func (l level) String() string {
	return fmt.Sprint(l.color(), string(l), "\x1b[0m")
}

func (l level) color() string {
	switch l {
	case InfoLevel:
		return "\x1b[34m"
	case DebugLevel:
		return "\x1b[36m"
	case WarnLevel:
		return "\x1b[33m"
	case ErrorLevel:
		return "\x1b[31m"
	default:
		return "\x1b[0m"
	}
}
