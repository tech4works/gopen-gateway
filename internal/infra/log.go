package infra

import (
	"bytes"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"
)

type logProvider struct {
}

func NewLogProvider() interfaces.LogProvider {
	return logProvider{}
}

func (l logProvider) InitializeLoggerOptions(ctx *gin.Context) {
	// obtemos os valores para imprimir nos logs da requisição atual
	traceId := ctx.GetHeader(consts.XTraceId)
	ip := ctx.GetHeader(consts.XForwardedFor)
	uri := ctx.Request.URL.String()
	method := ctx.Request.Method

	// setamos as opções globais de log
	logger.SetOptions(&logger.Options{
		HideArgCaller:         true,
		CustomAfterPrefixText: l.buildLoggerAfterPrefixText(traceId, ip, uri, method),
	})
}

func (l logProvider) BuildInitialRequestMessage(ctx *gin.Context) string {
	// inicializamos o body
	var bodyInfo any

	// obtemos o tipo do body e size do mesmo
	bodyType := ctx.GetHeader("Content-Type")
	bodySize := ctx.GetHeader("Content-Length")
	if helper.ContainsIgnoreCase(bodyType, "application/json") || helper.ContainsIgnoreCase(bodyType, "plain/text") {
		// caso ele seja json ou texto obtemos o mesmo para imprimir
		bodyBytes, _ := io.ReadAll(ctx.Request.Body)
		// voltamos para requisição, ja que podemos precisar nos handles futuros
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		// convertemos esses bytes para o any inicializado
		helper.SimpleConvertToDest(bodyBytes, &bodyInfo)
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
	return l.replaceAllBreakLineLogger(bodyInfo)
}

func (l logProvider) BuildFinishRequestMessage(writer dto.ResponseWriter, startTime time.Time) string {
	// obtemos quanto tempo demorou a requisição
	latency := time.Now().Sub(startTime)

	// inicializamos o text de mensagem de retorno
	var text strings.Builder

	// obtemos o texto de status code de resposta
	textStatusCode := l.getLoggerTextStatusCode(writer.Status())
	// obtemos o text de latência
	textLatency := latency.String()
	// obtemos o text do body de resposta
	textBody := string(writer.Body.Bytes())

	// montamos o texto
	text.WriteString(textStatusCode)
	text.WriteString(" • ")
	text.WriteString(textLatency)
	text.WriteString(" ")
	text.WriteString(logger.StyleReset)
	if helper.IsNotEmpty(textBody) {
		text.WriteString(" ")
		text.WriteString(textBody)
	}
	// retornamos o text de log
	return text.String()
}

func (l logProvider) buildLoggerAfterPrefixText(traceId, ip, uri, method string) string {
	textTrace := l.getLoggerTextTraceId(traceId)
	textMethod := l.getLoggerTextMethod(method)
	textUri := l.getLoggerTextUri(uri)
	return fmt.Sprint("(", textTrace, " | ", ip, " |", textMethod, "| ", textUri, ")")
}

func (l logProvider) getRequestUri(request *http.Request) (uri string) {
	uri = request.URL.Path
	raw := request.URL.RawQuery
	if helper.IsNotEmpty(raw) {
		uri = fmt.Sprint(uri, "?", raw)
	}
	return uri
}

func (l logProvider) getLoggerTextTraceId(traceId string) string {
	return fmt.Sprint(logger.StyleBold, traceId, logger.StyleReset)
}

func (l logProvider) getLoggerTextMethod(method string) string {
	return fmt.Sprint(l.getMethodColorTextLogger(method), " ", method, " ", logger.StyleReset)
}

func (l logProvider) getLoggerTextUri(uri string) string {
	return fmt.Sprint("\"", uri, "\"")
}

func (l logProvider) getLoggerTextStatusCode(statusCode int) string {
	return fmt.Sprint(l.getStatusCodeColorTextLogger(statusCode), " ", statusCode)
}

func (l logProvider) getStatusCodeColorTextLogger(statusCode int) string {
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

func (l logProvider) getMethodColorTextLogger(method string) string {
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

func (l logProvider) replaceAllBreakLineLogger(a any) string {
	s := helper.SimpleConvertToString(a)
	s = strings.ReplaceAll(s, "\n", " \\n ")
	return s
}
