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

package http

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/tech4works/checker"
)

type cbState int32

const (
	cbClosed   cbState = 0
	cbOpen     cbState = 1
	cbHalfOpen cbState = 2
)

type circuitBreaker struct {
	state            atomic.Int32
	consecutiveFails atomic.Int64
	consecutiveOK    atomic.Int64
	halfOpenCount    atomic.Int64
	lastFailTime     atomic.Int64

	failureThreshold int64
	successThreshold int64
	openTimeout      time.Duration
	halfOpenMaxReqs  int64

	mu sync.Mutex
}

func newCircuitBreaker(failureThreshold, successThreshold int, openTimeout time.Duration, halfOpenMaxReqs int) *circuitBreaker {
	cb := &circuitBreaker{
		failureThreshold: int64(failureThreshold),
		successThreshold: int64(successThreshold),
		openTimeout:      openTimeout,
		halfOpenMaxReqs:  int64(halfOpenMaxReqs),
	}
	cb.state.Store(int32(cbClosed))
	return cb
}

func (cb *circuitBreaker) Allow() bool {
	state := cbState(cb.state.Load())

	switch state {
	case cbClosed:
		return true
	case cbOpen:
		lastFail := time.Unix(0, cb.lastFailTime.Load())
		if checker.IsGreaterThanOrEqual(time.Since(lastFail), cb.openTimeout) {
			cb.mu.Lock()
			if checker.Equals(cbState(cb.state.Load()), cbOpen) {
				cb.state.Store(int32(cbHalfOpen))
				cb.halfOpenCount.Store(0)
				cb.consecutiveOK.Store(0)
			}
			cb.mu.Unlock()
			return cb.allowHalfOpen()
		}
		return false
	case cbHalfOpen:
		return cb.allowHalfOpen()
	}
	return true
}

func (cb *circuitBreaker) allowHalfOpen() bool {
	count := cb.halfOpenCount.Add(1)
	return checker.IsLessThanOrEqual(count, cb.halfOpenMaxReqs)
}

func (cb *circuitBreaker) RecordSuccess() {
	state := cbState(cb.state.Load())
	switch state {
	case cbClosed:
		cb.consecutiveFails.Store(0)
	case cbHalfOpen:
		successes := cb.consecutiveOK.Add(1)
		if checker.IsGreaterThanOrEqual(successes, cb.successThreshold) {
			cb.mu.Lock()
			if checker.Equals(cbState(cb.state.Load()), cbHalfOpen) {
				cb.state.Store(int32(cbClosed))
				cb.consecutiveFails.Store(0)
				cb.consecutiveOK.Store(0)
			}
			cb.mu.Unlock()
		}
	}
}

func (cb *circuitBreaker) RecordFailure() {
	cb.lastFailTime.Store(time.Now().UnixNano())
	state := cbState(cb.state.Load())
	switch state {
	case cbClosed:
		fails := cb.consecutiveFails.Add(1)
		if checker.IsGreaterThanOrEqual(fails, cb.failureThreshold) {
			cb.mu.Lock()
			if checker.Equals(cbState(cb.state.Load()), cbClosed) {
				cb.state.Store(int32(cbOpen))
			}
			cb.mu.Unlock()
		}
	case cbHalfOpen:
		cb.mu.Lock()
		if checker.Equals(cbState(cb.state.Load()), cbHalfOpen) {
			cb.state.Store(int32(cbOpen))
			cb.halfOpenCount.Store(0)
			cb.consecutiveOK.Store(0)
		}
		cb.mu.Unlock()
	}
}

func (cb *circuitBreaker) State() cbState {
	return cbState(cb.state.Load())
}

func (s cbState) String() string {
	switch s {
	case cbClosed:
		return "CLOSED"
	case cbOpen:
		return "OPEN"
	case cbHalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}
