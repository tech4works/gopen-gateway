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

package infra

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
	"strings"
)

// httpLoggerProvider is a struct that implements the HttpLoggerProvider interface.
// It provides methods to print HTTP request and response information.
type httpLoggerProvider struct {
}

// HttpLoggerProvider is an interface for logging HTTP request and response information.
type HttpLoggerProvider interface {
	PrintHttpRequestInfo(ctx *api.Context)
	PrintHttpResponseInfo(ctx *api.Context)
}

// NewHttpLoggerProvider creates and returns a new instance of HttpLoggerProvider.
func NewHttpLoggerProvider() HttpLoggerProvider {
	return httpLoggerProvider{}
}

// PrintHttpRequestInfo prints the information of the HTTP request.
// It initializes the logger options with the request data and then prints the log message.
//
// Parameters:
// - ctx: The API context containing the request data.
func (h httpLoggerProvider) PrintHttpRequestInfo(ctx *api.Context) {
	h.initializeLoggerOptions(ctx)

	var bodyInfo string
	bodyType := ctx.Header().Get("Content-Type")
	bodySize := ctx.Header().Get("Content-Length")
	if helper.ContainsIgnoreCase(bodyType, "application/json") ||
		helper.ContainsIgnoreCase(bodyType, "application/xml") ||
		helper.ContainsIgnoreCase(bodyType, "plain/text") {
		bodyInfo = helper.CompactString(ctx.BodyString())
	} else if helper.IsNotEmpty(bodyType, bodySize) {
		msg := fmt.Sprintf("content-type: %s ", bodyType)
		msg += fmt.Sprintf("content-length: %s ", bodySize)
		bodyInfo = msg
	}

	logger.Info("Start!", bodyInfo)
}

// PrintHttpResponseInfo prints the information of the HTTP response.
// It obtains the status code text, latency, and response body data to be logged.
// Then, it prints the log message using the logger.Info function.
// Parameters:
// - ctx: The API context containing the response data.
func (h httpLoggerProvider) PrintHttpResponseInfo(ctx *api.Context) {
	var text strings.Builder
	text.WriteString(h.statusCodeText(ctx.HttpResponse().StatusCode()))
	text.WriteString(" • ")
	text.WriteString(ctx.Latency().String())
	text.WriteString(" ")
	text.WriteString(logger.StyleReset)
	if helper.IsNotNil(ctx.HttpResponse().Body()) {
		text.WriteString(" ")
		text.WriteString(ctx.HttpResponse().Body().CompactString())
	}

	logger.Info("Finish!", text.String())
}

// InitializeLoggerOptions initializes the logger options for the current request.
// It obtains the values to be printed in the request logs such as traceId, IP address, URL, and method.
// Then, it sets the global log options with the values obtained.
func (h httpLoggerProvider) initializeLoggerOptions(ctx *api.Context) {
	traceId := ctx.TraceId()
	ip := ctx.Header().Get(consts.XForwardedFor)
	url := ctx.Url()
	method := ctx.Method()
	logger.SetOptions(&logger.Options{
		CustomAfterPrefixText: h.afterPrefixText(traceId, ip, url, method),
		HideArgCaller:         true,
		HideArgDatetime:       true,
	})
}

// afterPrefixText returns a formatted string that represents the portion of the logger message
// that comes after the log prefix. It includes the `trace` ID, IP address, HTTP method, and URI.
//
// Parameters:
// - traceId: the traceProvider ID value
// - ip: the IP address value
// - uri: the URI value
// - method: the HTTP method value
//
// Returns:
// - The formatted string with the `trace` ID, IP address, HTTP method, and URI enclosed in parentheses.
func (h httpLoggerProvider) afterPrefixText(traceId, ip, uri, method string) string {
	return fmt.Sprint("(", h.traceIdText(traceId), " | ", ip, " |", h.methodText(method), "| ", h.uriText(uri), ")")
}

// traceIdText returns the formatted traceId with bold style and resets the style afterward.
func (h httpLoggerProvider) traceIdText(traceId string) string {
	return fmt.Sprint(logger.StyleBold, traceId, logger.StyleReset)
}

// methodText returns the text representation of the logger method.
// It includes the color styling, method name, and the reset styling.
// Parameters:
// - method: The method name to be included in the logger text.
// Returns:
// - The formatted logger text for the method.
func (h httpLoggerProvider) methodText(method string) string {
	return fmt.Sprint(h.methodTextStyle(method), " ", method, " ", logger.StyleReset)
}

// uriText returns the URI enclosed in double quotes.
func (h httpLoggerProvider) uriText(uri string) string {
	return fmt.Sprint("\"", uri, "\"")
}

// statusCodeText returns the status code text to be logged.
// It calls the statusCodeTextStyle method to get the colorized text and concatenates it with the status code.
// Example: "200" or "200" (in color)
func (h httpLoggerProvider) statusCodeText(statusCode vo.StatusCode) string {
	return fmt.Sprint(h.statusCodeTextStyle(statusCode), " ", statusCode)
}

// statusCodeTextStyle returns the color style to be applied to the status code text.
// It takes a StatusCode value as input and returns the style as a string which can be used
// to format the log message. The style is determined based on the range of the given status code.
// If the status code is within the range of 200 to 299, the style is bold with a green background.
// If the status code is within the range of 300 to 399, the style is bold with a cyan background.
// If the status code is within the range of 400 to 499, the style is bold with a yellow background.
// If the status code is 500 or greater, the style is bold with a red background.
// If the status code does not fall into any of the above ranges, the style is bold.
//
// Parameters:
// - statusCode: The status code value.
//
// Returns:
// - The style as a string which can be used to format the log message.
func (h httpLoggerProvider) statusCodeTextStyle(statusCode vo.StatusCode) string {
	if helper.IsGreaterThanOrEqual(statusCode, 200) && helper.IsLessThan(statusCode, 299) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundGreen)
	} else if helper.IsGreaterThanOrEqual(statusCode, 300) && helper.IsLessThan(statusCode, 400) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundCyan)
	} else if helper.IsGreaterThanOrEqual(statusCode, 400) && helper.IsLessThan(statusCode, 500) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundYellow)
	} else if helper.IsGreaterThanOrEqual(statusCode, 500) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundRed)
	}
	return logger.StyleBold
}

// methodTextStyle returns the color text for the given HTTP method.
// It takes the HTTP method as input and returns the corresponding color text.
// The color text is used to format the log message.
// The color of the text is determined based on the HTTP method as follows:
// - For POST method, the color is bold yellow.
// - For GET method, the color is bold blue.
// - For DELETE method, the color is bold red.
// - For PUT method, the color is bold magenta.
// - For PATCH method, the color is bold cyan.
// - For any other method, the color is bold black.
func (h httpLoggerProvider) methodTextStyle(method string) string {
	switch method {
	case http.MethodPost:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundYellow)
	case http.MethodGet:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundBlue)
	case http.MethodDelete:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundRed)
	case http.MethodPut:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundMagenta)
	case http.MethodPatch:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundCyan)
	default:
		return fmt.Sprint(logger.StyleBold, logger.BackgroundBlack)
	}
}
