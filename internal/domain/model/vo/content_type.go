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
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"strings"
)

type ContentType string

func NewContentType(s string) ContentType {
	return ContentType(s)
}

func NewContentTypeTextPlain() ContentType {
	return "text/plain; charset=UTF-8"
}

func NewContentTypeJson() ContentType {
	return "application/json; charset=UTF-8"
}

func NewContentTypeXml() ContentType {
	return "application/xml; charset=UTF-8"
}

func (c ContentType) String() string {
	return string(c)
}

func (c ContentType) IsJSON() bool {
	return strings.HasPrefix(string(c), "application/json")
}

func (c ContentType) IsNotJSON() bool {
	return !c.IsJSON()
}

func (c ContentType) IsXML() bool {
	return strings.HasPrefix(string(c), "application/xml")
}

func (c ContentType) IsNotXML() bool {
	return !c.IsXML()
}

func (c ContentType) IsText() bool {
	return strings.HasPrefix(string(c), "text/plain")
}

func (c ContentType) IsNotText() bool {
	return !c.IsText()
}

func (c ContentType) IsUnknown() bool {
	return c.IsNotJSON() && c.IsNotXML() && c.IsNotText()
}

func (c ContentType) ToEnum() enum.ContentType {
	if c.IsText() {
		return enum.ContentTypePlainText
	} else if c.IsJSON() {
		return enum.ContentTypeJson
	} else if c.IsXML() {
		return enum.ContentTypeXml
	}
	return ""
}
