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
	"strings"

	"github.com/tech4works/checker"
)

type ContentType string

func NewContentType(s string) ContentType {
	return ContentType(s)
}

func NewContentTypeTextPlain() ContentType {
	return "text/plain"
}

func NewContentTypeJSON() ContentType {
	return "application/json"
}

func NewContentTypeXML() ContentType {
	return "application/xml"
}

func (c ContentType) String() string {
	return string(c)
}

func (c ContentType) IsJSON() bool {
	s := strings.ToLower(c.String())
	return strings.HasPrefix(s, "application/json") || strings.Contains(s, "+json")
}

func (c ContentType) IsNotJSON() bool {
	return !c.IsJSON()
}

func (c ContentType) IsXML() bool {
	s := strings.ToLower(c.String())
	return strings.HasPrefix(s, "application/xml") || strings.Contains(s, "+xml")
}

func (c ContentType) IsNotXML() bool {
	return !c.IsXML()
}

func (c ContentType) IsPlainText() bool {
	s := strings.ToLower(c.String())
	return strings.HasPrefix(s, "text/plain")
}

func (c ContentType) IsNotPlainText() bool {
	return !c.IsPlainText()
}

func (c ContentType) IsSupported() bool {
	return c.IsJSON() || c.IsXML() || c.IsPlainText()
}

func (c ContentType) IsUnsupported() bool {
	return !c.IsSupported()
}

func (c ContentType) Equals(another ContentType) bool {
	return checker.Equals(c.IsPlainText(), another.IsPlainText()) ||
		checker.Equals(c.IsJSON(), another.IsJSON()) ||
		checker.Equals(c.IsXML(), another.IsXML())
}
