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

package consts

const (
	ContentType     = "Content-Type"
	ContentEncoding = "Content-Encoding"
	ContentLength   = "Content-Length"
	// XForwardedFor represents the name of the "X-Forwarded-For" HTTP header.
	// It is used to indicate the original IP address of a client connecting to a web server
	// through an HTTP proxy or a load balancer.
	XForwardedFor = "X-Forwarded-For"
	// XGopenCache represents the name of the "X-Gopen-Cache" HTTP header.
	// It is used to indicate whether a cache is being used for the request.
	XGopenCache = "X-Gopen-Cache"
	// XGopenCacheTTL represents the name of the "X-Gopen-Cache-Ttl" HTTP header.
	// It is used to indicate the time-to-live (TTL) value of a cache response,
	// which specifies how long the response should be considered valid and can be reused.
	// The value of X-Gopen-Cache-Ttl header is typically a duration in seconds.
	XGopenCacheTTL = "X-Gopen-Cache-Ttl"
	// XGopenComplete represents the name of the "X-Gopen-Complete" HTTP header. It is used to indicate the completion
	// status of a request.
	XGopenComplete = "X-Gopen-Complete"
	// XGopenSuccess represents the name of the "X-Gopen-Success" HTTP header.
	// It is used to indicate the success status of a request.
	XGopenSuccess = "X-Gopen-Success"
)
