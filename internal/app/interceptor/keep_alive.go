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

package interceptor

import (
	"fmt"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type keepAliveMiddleware struct {
	metadata vo.Metadata
}

type KeepAlive interface {
	Do(ctx app.Context)
}

// NewKeepAlive creates a middleware that injects Keep-Alive header in responses.
// The timeout value signals HTTP clients (especially OkHttp on Android) when to discard
// idle connections, preventing stale connection errors (ERR_NETWORK).
//
// Standard behavior:
//   - keep-alive enabled: injects "Keep-Alive: timeout=X" where X = idleTimeout - 5s
//     (margin lets client close before server, avoiding RST races)
//   - keep-alive disabled: injects "Connection: close" (net/http handles this,
//     but the middleware is a no-op in this case)
func NewKeepAlive(serverConfig *vo.ServerConfig) KeepAlive {
	if !serverConfig.KeepAlive() {
		return keepAliveMiddleware{metadata: vo.NewEmptyMetadata()}
	}

	// Signal clients to drop connection 5s before server's idle timeout.
	// This prevents the race condition where server closes first and client
	// sends on a dead socket (common on Android/OkHttp).
	timeoutSeconds := int(serverConfig.IdleTimeout().Seconds()) - 5
	if timeoutSeconds < 1 {
		timeoutSeconds = 1
	}

	metadata := vo.NewMetadata(map[string][]string{
		"Keep-Alive": {fmt.Sprintf("timeout=%d", timeoutSeconds)},
	})

	return keepAliveMiddleware{metadata: metadata}
}

func (k keepAliveMiddleware) Do(ctx app.Context) {
	if len(k.metadata.Keys()) > 0 {
		ctx.WriteMetadata(k.metadata)
	}
	ctx.Next()
}
