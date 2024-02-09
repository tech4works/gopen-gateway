package usecase

import (
	"bytes"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/handler"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/enum"
	"io"
	"net/http"
	"strings"
	"time"
)

type log struct {
}

type Log interface {
	PrintLogRequest(request *http.Request)
	PrintLogResponse(request *http.Request, responseWriter handler.ResponseWriter, startTime time.Time)
}

func NewLogger() Log {
	return log{}
}

func (l log) PrintLogRequest(request *http.Request) {
	l.initializeLoggerOptions(request)
	var body any
	bodyBytes, _ := io.ReadAll(request.Body)
	request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	helper.SimpleConvertToDest(string(bodyBytes), &body)
	logger.Info("Start!", l.replaceAllBreakLine(body))
}

func (l log) PrintLogResponse(request *http.Request, responseWriter handler.ResponseWriter, startTime time.Time) {
	responseTime := time.Now()
	latency := responseTime.Sub(startTime)
	l.initializeLoggerOptions(request)
	var text strings.Builder
	textStatusCode := l.getTextStatusCode(responseWriter.Status())
	textLatency := latency.String()
	textBody := string(responseWriter.Body.Bytes())
	text.WriteString(textStatusCode)
	text.WriteString(logger.StyleBold)
	text.WriteString(" • ")
	text.WriteString(logger.StyleReset)
	text.WriteString(textLatency)
	if helper.IsNotEmpty(textBody) {
		text.WriteString(logger.StyleBold)
		text.WriteString(": ")
		text.WriteString(logger.StyleReset)
		text.WriteString(textBody)
	}
	logger.Info("Finish!", l.replaceAllBreakLine(text.String()))
}

func (l log) initializeLoggerOptions(request *http.Request) {
	traceId := request.Header.Get(enum.XTraceId)
	ip := request.Header.Get(enum.XForwardedFor)
	uri := l.getRequestUri(request)
	method := request.Method
	logger.SetOptions(&logger.Options{
		HideArgCaller:         true,
		CustomAfterPrefixText: l.getAfterPrefixText(traceId, ip, uri, method),
	})
}

func (l log) getRequestUri(request *http.Request) (uri string) {
	uri = request.URL.Path
	raw := request.URL.RawQuery
	if helper.IsNotEmpty(raw) {
		uri = fmt.Sprint(uri, "?", raw)
	}
	return uri
}

func (l log) getTextTraceId(traceId string) string {
	return fmt.Sprint(logger.StyleBold, traceId, logger.StyleReset)
}

func (l log) getAfterPrefixText(traceId, ip, uri, method string) string {
	div := fmt.Sprint(" ", l.divText(), " ")
	doubleDot := fmt.Sprint(logger.StyleBold, ":", logger.StyleReset)
	textTrace := l.getTextTraceId(traceId)
	textMethod := l.getTextMethod(method)
	textUri := l.getTextUri(uri)
	return fmt.Sprint(textTrace, div, ip, div, textMethod, div, textUri, doubleDot)
}

func (l log) getTextMethod(method string) string {
	return fmt.Sprint(l.getMethodColorTextLogger(method), method, logger.StyleReset)
}

func (l log) getTextUri(uri string) string {
	return uri
}

func (l log) getTextStatusCode(statusCode int) string {
	return fmt.Sprint(l.getStatusCodeColorTextLogger(statusCode), statusCode, logger.StyleReset)
}

func (l log) divText() string {
	return fmt.Sprint(logger.StyleBold, "|", logger.StyleReset)
}

func (l log) getStatusCodeColorTextLogger(statusCode int) string {
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

func (l log) getMethodColorTextLogger(method string) string {
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

func (l log) replaceAllBreakLine(a any) string {
	s := helper.SimpleConvertToString(a)
	s = strings.ReplaceAll(s, "\n", " \\n ")
	return s
}
