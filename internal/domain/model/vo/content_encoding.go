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
)

type ContentEncoding string

func NewContentEncoding(s string) ContentEncoding {
	return ContentEncoding(s)
}

func NewContentEncodingGzip() ContentEncoding {
	return "gzip"
}

func NewContentEncodingDeflate() ContentEncoding {
	return "deflate"
}

func (c ContentEncoding) String() string {
	return string(c)
}

func (c ContentEncoding) IsSupported() bool {
	return c.IsGzip() || c.IsDeflate()
}

func (c ContentEncoding) Valid() bool {
	return checker.IsNotEmpty(c)
}

func (c ContentEncoding) Invalid() bool {
	return !c.Valid()
}

func (c ContentEncoding) IsGzip() bool {
	return checker.EqualsIgnoreCase(c, "gzip")
}

func (c ContentEncoding) IsDeflate() bool {
	return checker.EqualsIgnoreCase(c, "deflate")
}
