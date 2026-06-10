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
	"time"

	"github.com/tech4works/checker"
)

// ServerConfig holds HTTP server-level settings that control connection management.
// Defaults follow industry-standard API gateway practices (NGINX, Envoy, Traefik).
type ServerConfig struct {
	readTimeout       Duration
	writeTimeout      Duration
	readHeaderTimeout Duration
	idleTimeout       Duration
	keepAlive         bool
}

func NewServerConfig(readTimeout, writeTimeout, readHeaderTimeout, idleTimeout Duration, keepAlive bool) *ServerConfig {
	return &ServerConfig{
		readTimeout:       readTimeout,
		writeTimeout:      writeTimeout,
		readHeaderTimeout: readHeaderTimeout,
		idleTimeout:       idleTimeout,
		keepAlive:         keepAlive,
	}
}

// ReadTimeout returns the max duration for reading the entire request (headers + body).
// Default: 60s (aligned with NGINX client_body_timeout + client_header_timeout).
func (s *ServerConfig) ReadTimeout() time.Duration {
	if checker.IsGreaterThan(s.readTimeout, 0) {
		return s.readTimeout.Time()
	}
	return 60 * time.Second
}

// WriteTimeout returns the max duration for writing the response.
// Default: 60s (aligned with NGINX proxy_send_timeout).
func (s *ServerConfig) WriteTimeout() time.Duration {
	if checker.IsGreaterThan(s.writeTimeout, 0) {
		return s.writeTimeout.Time()
	}
	return 60 * time.Second
}

// ReadHeaderTimeout returns the max duration for reading request headers.
// Default: 10s (protects against slowloris attacks).
func (s *ServerConfig) ReadHeaderTimeout() time.Duration {
	if checker.IsGreaterThan(s.readHeaderTimeout, 0) {
		return s.readHeaderTimeout.Time()
	}
	return 10 * time.Second
}

// IdleTimeout returns how long to keep idle keep-alive connections open.
// Default: 65s (slightly above typical client timeout of 60s to let client close first).
// The gateway injects Keep-Alive header with timeout = IdleTimeout - 5s to signal
// clients to close before the server does (avoids stale connection races).
func (s *ServerConfig) IdleTimeout() time.Duration {
	if checker.IsGreaterThan(s.idleTimeout, 0) {
		return s.idleTimeout.Time()
	}
	return 65 * time.Second
}

// KeepAlive returns whether HTTP keep-alive connections are enabled.
// Default: true. When false, the server sends Connection: close on every response.
func (s *ServerConfig) KeepAlive() bool {
	return s.keepAlive
}
