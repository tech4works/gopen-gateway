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

package app

const (
	ContentType                   = "Content-Type"
	ContentEncoding               = "Content-Encoding"
	ContentLength                 = "Content-Length"
	XForwardedFor                 = "X-Forwarded-For"
	XGopenRequestID               = "X-Gopen-Request-Id"
	XGopenHeaderDegraded          = "X-Gopen-Header-Degraded"
	XGopenMetadataDegraded        = "X-Gopen-Metadata-Degraded"
	XGopenQueryDegraded           = "X-Gopen-Query-Degraded"
	XGopenURLPathDegraded         = "X-Gopen-Url-Path-Degraded"
	XGopenBodyDegraded            = "X-Gopen-Body-Degraded"
	XGopenPayloadDegraded         = "X-Gopen-Payload-Degraded"
	XGopenDeduplicationIDDegraded = "X-Gopen-Deduplication-Id-Degraded"
	XGopenGroupIDDegraded         = "X-Gopen-Group-Id-Degraded"
	XGopenAttributeDegraded       = "X-Gopen-Attribute-Degraded"
	XGopenTimeout                 = "X-Gopen-Timeout"
	XGopenCache                   = "X-Gopen-Cache"
	XGopenCacheTTL                = "X-Gopen-Cache-Ttl"
	XGopenDegraded                = "X-Gopen-Degraded"
	XGopenDegradedBackendCount    = "X-Gopen-Degraded-Backend-Count"
	XGopenDegradedBackends        = "X-Gopen-Degraded-Backends"
	XGopenComplete                = "X-Gopen-Complete"
	XGopenSuccess                 = "X-Gopen-Success"
)

func TransportHTTPHeaderKeys() []string {
	return []string{
		ContentType,
		ContentEncoding,
		ContentLength,
		XGopenCache,
		XGopenCacheTTL,
		XGopenDegraded,
		XGopenComplete,
		XGopenSuccess,
		XGopenTimeout,
	}
}
