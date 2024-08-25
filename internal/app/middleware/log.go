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

package middleware

import (
	"github.com/tech4works/gopen-gateway/internal/app"
)

type logMiddleware struct {
	log app.HTTPLog
}

type Log interface {
	Do(ctx app.Context)
}

func NewLog(log app.HTTPLog) Log {
	return logMiddleware{
		log: log,
	}
}

func (l logMiddleware) Do(ctx app.Context) {
	l.log.PrintRequest(ctx)

	ctx.Next()

	l.log.PrintResponse(ctx)
}
