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

package mapper

import (
	"github.com/tech4works/checker"
)

const (
	ContentType     = "Content-Type"
	ContentEncoding = "Content-Encoding"
	ContentLength   = "Content-Length"
	XForwardedFor   = "X-Forwarded-For"
	XGopenCache     = "X-Gopen-Cache"
	XGopenCacheTTL  = "X-Gopen-Cache-Ttl"
	XGopenComplete  = "X-Gopen-Complete"
	XGopenSuccess   = "X-Gopen-Success"
)

func mandatoryHeaderKeys() []string {
	return []string{ContentType, ContentEncoding, ContentLength, XForwardedFor, XGopenCache, XGopenCacheTTL,
		XGopenComplete, XGopenSuccess}
}

func gopenHeaderKeys() []string {
	return []string{XGopenCache, XGopenCacheTTL, XGopenComplete, XGopenSuccess}
}

func IsHeaderMandatoryKey(key string) bool {
	return checker.Contains(mandatoryHeaderKeys(), key)
}

func IsNotHeaderMandatoryKey(key string) bool {
	return !IsHeaderMandatoryKey(key)
}
