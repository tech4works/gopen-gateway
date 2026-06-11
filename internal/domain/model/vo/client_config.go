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

// ClientConfig holds HTTP client-level settings for outbound requests to backends.
// Controls connection pooling, circuit breaker, retry and client timeout.
type ClientConfig struct {
	timeout             Duration
	maxIdleConns        int
	maxIdleConnsPerHost int
	idleConnTimeout     Duration
	circuitBreaker      CircuitBreakerConfig
	retry               RetryConfig
}

// CircuitBreakerConfig holds circuit breaker settings per backend host.
type CircuitBreakerConfig struct {
	failureThreshold int
	successThreshold int
	openTimeout      Duration
	halfOpenMaxReqs  int
}

// RetryConfig holds retry settings for idempotent requests (GET/HEAD).
type RetryConfig struct {
	maxRetries int
	backoff    Duration
}

func NewClientConfig(
	timeout Duration,
	maxIdleConns int,
	maxIdleConnsPerHost int,
	idleConnTimeout Duration,
	cb CircuitBreakerConfig,
	retry RetryConfig,
) *ClientConfig {
	return &ClientConfig{
		timeout:             timeout,
		maxIdleConns:        maxIdleConns,
		maxIdleConnsPerHost: maxIdleConnsPerHost,
		idleConnTimeout:     idleConnTimeout,
		circuitBreaker:      cb,
		retry:               retry,
	}
}

func NewCircuitBreakerConfig(failureThreshold, successThreshold int, openTimeout Duration, halfOpenMaxReqs int) CircuitBreakerConfig {
	return CircuitBreakerConfig{
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		openTimeout:      openTimeout,
		halfOpenMaxReqs:  halfOpenMaxReqs,
	}
}

func NewRetryConfig(maxRetries int, backoff Duration) RetryConfig {
	return RetryConfig{
		maxRetries: maxRetries,
		backoff:    backoff,
	}
}

// Timeout returns the max duration for a single outbound HTTP request (safety net).
// Default: 300s.
func (c *ClientConfig) Timeout() time.Duration {
	if checker.IsGreaterThan(c.timeout, 0) {
		return c.timeout.Time()
	}
	return 300 * time.Second
}

// MaxIdleConns returns total idle connections in the pool.
// Default: 200.
func (c *ClientConfig) MaxIdleConns() int {
	if checker.IsGreaterThan(c.maxIdleConns, 0) {
		return c.maxIdleConns
	}
	return 200
}

// MaxIdleConnsPerHost returns idle connections per backend host.
// Default: 50 (gateways fan out to many backends concurrently).
func (c *ClientConfig) MaxIdleConnsPerHost() int {
	if checker.IsGreaterThan(c.maxIdleConnsPerHost, 0) {
		return c.maxIdleConnsPerHost
	}
	return 50
}

// IdleConnTimeout returns how long idle connections stay in the pool.
// Default: 90s.
func (c *ClientConfig) IdleConnTimeout() time.Duration {
	if checker.IsGreaterThan(c.idleConnTimeout, 0) {
		return c.idleConnTimeout.Time()
	}
	return 90 * time.Second
}

// CBFailureThreshold returns consecutive failures to trip the breaker open.
// Default: 5.
func (c *ClientConfig) CBFailureThreshold() int {
	if checker.IsGreaterThan(c.circuitBreaker.failureThreshold, 0) {
		return c.circuitBreaker.failureThreshold
	}
	return 5
}

// CBSuccessThreshold returns consecutive successes in half-open to close the breaker.
// Default: 2.
func (c *ClientConfig) CBSuccessThreshold() int {
	if checker.IsGreaterThan(c.circuitBreaker.successThreshold, 0) {
		return c.circuitBreaker.successThreshold
	}
	return 2
}

// CBOpenTimeout returns how long the breaker stays open before testing.
// Default: 30s.
func (c *ClientConfig) CBOpenTimeout() time.Duration {
	if checker.IsGreaterThan(c.circuitBreaker.openTimeout, 0) {
		return c.circuitBreaker.openTimeout.Time()
	}
	return 30 * time.Second
}

// CBHalfOpenMaxReqs returns requests allowed in half-open state.
// Default: 2.
func (c *ClientConfig) CBHalfOpenMaxReqs() int {
	if checker.IsGreaterThan(c.circuitBreaker.halfOpenMaxReqs, 0) {
		return c.circuitBreaker.halfOpenMaxReqs
	}
	return 2
}

// RetryMaxRetries returns max retries for GET/HEAD on transient errors.
// Default: 1.
func (c *ClientConfig) RetryMaxRetries() int {
	if checker.IsGreaterThanOrEqual(c.retry.maxRetries, 0) {
		return c.retry.maxRetries
	}
	return 1
}

// RetryBackoff returns delay between retry attempts.
// Default: 100ms.
func (c *ClientConfig) RetryBackoff() time.Duration {
	if checker.IsGreaterThan(c.retry.backoff, 0) {
		return c.retry.backoff.Time()
	}
	return 100 * time.Millisecond
}
