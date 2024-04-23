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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
	"strings"
	"time"
	"unicode"
)

type logProvider struct {
}

// LogProvider is an interface that defines methods for logging in a software application.
// It provides functionality to initialize logger options, build an initial request message,
// and build a finish request message.
type LogProvider interface {
	// InitializeLoggerOptions initializes the logger options using the provided context.
	InitializeLoggerOptions(ctx *api.Context)
	// BuildInitialRequestMessage builds the initial request message for logging purposes based on the provided context.
	BuildInitialRequestMessage(ctx *api.Context) string
	// BuildFinishRequestMessage constructs a finish request message for logging purposes.
	// It takes a responseVO object that represents the response of an API call, and a startTime
	// object that represents the start time of the API call. It returns a string containing the
	// finish request message.
	BuildFinishRequestMessage(responseVO *vo.Response, startTime time.Time) string
}

// NewLogProvider creates and returns a new instance of LogProvider.
func NewLogProvider() LogProvider {
	return logProvider{}
}

// InitializeLoggerOptions initializes the logger options for the current request.
// It obtains the values to be printed in the request logs such as traceId, IP address, URL, and method.
// Then, it sets the global log options with the values obtained.
func (l logProvider) InitializeLoggerOptions(ctx *api.Context) {
	// obtemos os valores para imprimir nos logs da requisição atual
	traceId := ctx.HeaderValue(consts.XTraceId)
	ip := ctx.HeaderValue(consts.XForwardedFor)
	url := ctx.Url()
	method := ctx.Method()

	// setamos as opções globais de log
	logger.SetOptions(&logger.Options{
		HideArgCaller:         true,
		CustomAfterPrefixText: l.afterPrefixText(traceId, ip, url, method),
	})
}

// BuildInitialRequestMessage builds the initial request message for logging purposes.
// It initializes the `bodyInfo` variable and obtains the body type and size from the request headers.
// If the body type is "application/json" or "application/xml" or "plain/text", it converts the body to a string and assigns it to `bodyInfo`.
// Otherwise, if the body type and size are not empty, it creates a message string with the content type and content length and assigns it to `bodyInfo`.
// Finally, it converts `bodyInfo` to a string, removes any line breaks, and returns the result.
func (l logProvider) BuildInitialRequestMessage(ctx *api.Context) string {
	// inicializamos o body
	var bodyInfo string

	// obtemos o tipo do body e size do mesmo
	bodyType := ctx.HeaderValue("Content-Type")
	bodySize := ctx.HeaderValue("Content-Length")
	if helper.ContainsIgnoreCase(bodyType, "application/json") ||
		helper.ContainsIgnoreCase(bodyType, "application/xml") ||
		helper.ContainsIgnoreCase(bodyType, "plain/text") {
		// convertemos esses bytes para o any inicializado
		bodyInfo = ctx.BodyString()
	} else if helper.IsNotEmpty(bodyType, bodySize) {
		var msg string
		if helper.IsNotEmpty(bodyType) {
			msg += fmt.Sprintf("content-type: %s ", bodyType)
		}
		if helper.IsNotEmpty(bodyType) {
			msg += fmt.Sprintf("content-length: %s ", bodySize)
		}
		// caso não seja o json e text, imprimimos um resumo
		bodyInfo = msg
	}

	// convertemos em string e removemos os breaks lines
	return l.replaceAllBreakLineText(bodyInfo)
}

// BuildFinishRequestMessage builds the finish request message by taking a writer and start time as input.
// It calculates the latency of the request and initializes a string builder to store the message text.
// The method obtains the status code text, latency text, and response body text.
// It then constructs the message by appending the status code, latency, and response body (if not empty) to the string builder.
// Finally, it returns the build log message as a string.
func (l logProvider) BuildFinishRequestMessage(responseVO *vo.Response, startTime time.Time) string {
	// obtemos quanto tempo demorou a requisição
	latency := time.Now().Sub(startTime)

	// inicializamos o text de mensagem de retorno
	var text strings.Builder

	// obtemos o texto de status code de resposta
	textStatusCode := l.statusCodeText(responseVO.StatusCode())
	// obtemos o text de latência
	textLatency := latency.String()
	// obtemos o text do body de resposta
	bodyBytes := responseVO.BytesBody()

	// montamos o texto
	text.WriteString(textStatusCode)
	text.WriteString(" • ")
	text.WriteString(textLatency)
	text.WriteString(" ")
	text.WriteString(logger.StyleReset)
	if helper.IsNotEmpty(bodyBytes) {
		s := string(bodyBytes)
		if helper.IsNotEqualTo(responseVO.ContentType(), enum.ContentTypeText) {
			s = l.replaceAllBreakLineText(s)
		}
		text.WriteString(" ")
		text.WriteString(s)
	}
	// retornamos o text de log
	return text.String()
}

// afterPrefixText returns a formatted string that represents the portion of the logger message
// that comes after the log prefix. It includes the `trace` ID, IP address, HTTP method, and URI.
//
// Parameters:
// - traceId: the trace ID value
// - ip: the IP address value
// - uri: the URI value
// - method: the HTTP method value
//
// Returns:
// - The formatted string with the `trace` ID, IP address, HTTP method, and URI enclosed in parentheses.
func (l logProvider) afterPrefixText(traceId, ip, uri, method string) string {
	return fmt.Sprint("(", l.traceIdText(traceId), " | ", ip, " |", l.methodText(method), "| ", l.uriText(uri), ")")
}

// traceIdText returns the formatted traceId with bold style and resets the style afterward.
func (l logProvider) traceIdText(traceId string) string {
	return fmt.Sprint(logger.StyleBold, traceId, logger.StyleReset)
}

// methodText returns the text representation of the logger method.
// It includes the color styling, method name, and the reset styling.
// Parameters:
// - method: The method name to be included in the logger text.
// Returns:
// - The formatted logger text for the method.
func (l logProvider) methodText(method string) string {
	return fmt.Sprint(l.methodTextStyle(method), " ", method, " ", logger.StyleReset)
}

// uriText returns the URI enclosed in double quotes.
func (l logProvider) uriText(uri string) string {
	return fmt.Sprint("\"", uri, "\"")
}

// statusCodeText returns the status code text to be logged.
// It calls the statusCodeTextStyle method to get the colorized text and concatenates it with the status code.
// Example: "200" or "200" (in color)
func (l logProvider) statusCodeText(statusCode int) string {
	return fmt.Sprint(l.statusCodeTextStyle(statusCode), " ", statusCode)
}

// statusCodeTextStyle returns the color text for the given status code.
// It follows the below conditions to determine the color:
//
// If the status code is greater than or equal to http.StatusOK and less than http.StatusMultipleChoices,
// it returns the color text as logger.StyleBold and logger.BackgroundGreen.
//
// If the status code is greater than or equal to http.StatusMultipleChoices and less than http.StatusBadRequest,
// it returns the color text as logger.StyleBold and logger.BackgroundCyan.
//
// If the status code is greater than or equal to http.StatusBadRequest and less than http.StatusInternalServerError,
// it returns the color text as logger.StyleBold and logger.BackgroundYellow.
//
// If the status code is greater than or equal to http.StatusInternalServerError,
// it returns the color text as logger.StyleBold and logger.BackgroundRed.
//
// If none of the above conditions are met, it returns the color text as logger.StyleBold.
func (l logProvider) statusCodeTextStyle(statusCode int) string {
	if helper.IsGreaterThanOrEqual(statusCode, http.StatusOK) &&
		helper.IsLessThan(statusCode, http.StatusMultipleChoices) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundGreen)
	} else if helper.IsGreaterThanOrEqual(statusCode, http.StatusMultipleChoices) &&
		helper.IsLessThan(statusCode, http.StatusBadRequest) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundCyan)
	} else if helper.IsGreaterThanOrEqual(statusCode, http.StatusBadRequest) &&
		helper.IsLessThan(statusCode, http.StatusInternalServerError) {
		return fmt.Sprint(logger.StyleBold, logger.BackgroundYellow)
	} else if helper.IsGreaterThanOrEqual(statusCode, http.StatusInternalServerError) {
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
func (l logProvider) methodTextStyle(method string) string {
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

// replaceAllBreakLineText replaces all the line breaks in the given string with empty spaces.
// It takes an argument of any type and converts it to a string using the helper.SimpleConvertToString function.
// Then, it uses strings.Map to iterate over each rune in the string.
// If the rune is a space, it replaces it with -1 to exclude it from the string.
// If the rune is not a space, it keeps the original value.
// Finally, it returns the modified string.
func (l logProvider) replaceAllBreakLineText(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// se for um espaço, substitua por -1 para excluir da string
			return -1
		}
		// se não for um espaço, mantenha o valor original
		return r
	}, s)
}
