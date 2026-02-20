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

	"github.com/tech4works/checker"
)

func Print(lvl level, tag, prefix string, msg ...any) {
	tagText := BuildTagText(tag)
	levelText := BuildLevelText(lvl)

	if checker.IsNotEmpty(prefix) {
		fmt.Printf("[%s] %s %s %s", tagText, levelText, prefix, fmt.Sprintln(msg...))
	} else {
		fmt.Printf("[%s] %s %s", tagText, levelText, fmt.Sprintln(msg...))
	}
}

func Printf(lvl level, tag, prefix, format string, msg ...any) {
	tagText := BuildTagText(tag)
	levelText := BuildLevelText(lvl)

	if checker.IsNotEmpty(prefix) {
		fmt.Printf("[%s] %s %s %s\n", tagText, levelText, prefix, fmt.Sprintf(format, msg...))
	} else {
		fmt.Printf("[%s] %s %s\n", tagText, levelText, fmt.Sprintf(format, msg...))
	}
}
