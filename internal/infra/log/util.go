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

package log

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

func BuildLevelText(lvl level) string {
	return fmt.Sprint(StyleBold, lvl.String(), StyleReset)
}

func BuildTagText(tag string) string {
	return fmt.Sprint(StyleBold, tag, StyleReset)
}

func BuildTraceIDText(traceID string) string {
	return traceID
}

func BuildTintText(method string) string {
	return fmt.Sprint(tintTextStyle(method), " ", method, " ", StyleReset)
}

func BuildUriText(uri string) string {
	return strconv.Quote(uri)
}

func BuildStatusCodeText(statusCode vo.StatusCode) string {
	return fmt.Sprint(statusCodeTextStyle(statusCode.Code()), " ", statusCode.Code(), " ", StyleReset)
}

func statusCodeTextStyle(code int) string {
	if checker.IsGreaterThanOrEqual(code, 200) && checker.IsLessThan(code, 299) {
		return BackgroundGreen
	} else if checker.IsGreaterThanOrEqual(code, 300) && checker.IsLessThan(code, 400) {
		return BackgroundCyan
	} else if checker.IsGreaterThanOrEqual(code, 400) && checker.IsLessThan(code, 500) {
		return BackgroundYellow
	} else if checker.IsGreaterThanOrEqual(code, 500) {
		return BackgroundRed
	}
	return StyleBold
}

func tintTextStyle(s string) string {
	switch s {
	case http.MethodPost:
		return fmt.Sprint(BackgroundYellow)
	case http.MethodGet, enum.BackendBrokerAwsSns.String(), enum.BackendBrokerAwsSqs.String():
		return fmt.Sprint(BackgroundBlue)
	case http.MethodDelete:
		return fmt.Sprint(BackgroundRed)
	case http.MethodPut:
		return fmt.Sprint(BackgroundMagenta)
	case http.MethodPatch:
		return fmt.Sprint(BackgroundCyan)
	default:
		return fmt.Sprint(BackgroundBlack)
	}
}
