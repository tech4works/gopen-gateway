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
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
	"strings"
)

type loggerProvider struct {
}

func NewLoggerProvider() interfaces.LoggerProvider {
	return loggerProvider{}
}

func (l loggerProvider) PrintBackendErrorf(backend *vo.Backend, format string, msg ...any) {
	logger.ErrorOptsf(format, l.backendLoggerOptions(backend), msg...)
}

func (l loggerProvider) PrintBackendResponseInfo(backend *vo.Backend, httpBackendResponse *vo.HttpBackendResponse) {
	var text strings.Builder
	text.WriteString(statusCodeText(httpBackendResponse.StatusCode()))
	text.WriteString(" - ")
	text.WriteString(httpBackendResponse.Latency().String())
	text.WriteString(" ")
	text.WriteString(logger.StyleReset)
	logger.InfoOpts(l.backendLoggerOptions(backend), text.String())
}

func (l loggerProvider) PrintEndpointWarnf(ctx *api.Context, format string, msg ...any) {
	logger.WarnOptsf(format, l.endpointLoggerOptions(ctx), msg...)
}

func (l loggerProvider) PrintEndpointErrorf(ctx *api.Context, format string, msg ...any) {
	logger.ErrorOptsf(format, l.endpointLoggerOptions(ctx), msg...)
}

func (l loggerProvider) PrintEndpointResponseInfo(ctx *api.Context) {
	var text strings.Builder
	text.WriteString(statusCodeText(ctx.HttpResponse().StatusCode()))
	text.WriteString(" - ")
	text.WriteString(ctx.Latency().String())
	text.WriteString(" ")
	text.WriteString(logger.StyleReset)
	logger.InfoOpts(l.endpointLoggerOptions(ctx), text.String())
}

func (l loggerProvider) PrintHttpResponseInfo(ctx *api.Context) {
	var text strings.Builder
	text.WriteString(statusCodeText(ctx.HttpResponse().StatusCode()))
	text.WriteString(" - ")
	text.WriteString(ctx.Latency().String())
	text.WriteString(" ")
	text.WriteString(logger.StyleReset)
	logger.InfoOpts(l.httpResponseLoggerOptions(ctx), text.String())
}

func (l loggerProvider) backendLoggerOptions(backend *vo.Backend) logger.Options {
	path := backend.Path()
	prefix := fmt.Sprint("BKD (", uriText(path), ")")
	return logger.Options{
		CustomAfterPrefixText: prefix,
		HideAllArgs:           true,
	}
}

func (l loggerProvider) endpointLoggerOptions(ctx *api.Context) logger.Options {
	traceId := ctx.TraceId()
	ip := ctx.Header().Get(consts.XForwardedFor)
	path := ctx.Endpoint().Path()
	prefix := fmt.Sprint("END (", traceIdText(traceId), " | ", ip, "| ", uriText(path), ")")
	return logger.Options{
		CustomAfterPrefixText: prefix,
		HideAllArgs:           true,
	}
}

func (l loggerProvider) httpResponseLoggerOptions(ctx *api.Context) logger.Options {
	traceId := ctx.TraceId()
	ip := ctx.Header().Get(consts.XForwardedFor)
	url := ctx.Url()
	method := ctx.Method()
	prefix := fmt.Sprint("API (", traceIdText(traceId), " | ", ip, " |", methodText(method), "| ", uriText(url), ")")
	return logger.Options{
		CustomAfterPrefixText: prefix,
		HideAllArgs:           true,
	}
}

// traceIdText returns the formatted traceId with bold style and resets the style afterward.
func traceIdText(traceId string) string {
	return fmt.Sprint(logger.StyleBold, traceId, logger.StyleReset)
}

// methodText returns the text representation of the logger method.
// It includes the color styling, method name, and the reset styling.
// Parameters:
// - method: The method name to be included in the logger text.
// Returns:
// - The formatted logger text for the method.
func methodText(method string) string {
	return fmt.Sprint(methodTextStyle(method), " ", method, " ", logger.StyleReset)
}

// uriText returns the URI enclosed in double quotes.
func uriText(uri string) string {
	return fmt.Sprint("\"", uri, "\"")
}

// statusCodeText returns the status code text to be logged.
// It calls the statusCodeTextStyle method to get the colorized text and concatenates it with the status code.
// Example: "200" or "200" (in color)
func statusCodeText(statusCode vo.StatusCode) string {
	return fmt.Sprint(statusCodeTextStyle(statusCode), " ", statusCode)
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
func statusCodeTextStyle(statusCode vo.StatusCode) string {
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
func methodTextStyle(method string) string {
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
