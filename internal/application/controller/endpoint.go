package controller

import (
	"bytes"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/handler"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/dto"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/usecase"
	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/oj"
	"io"
	"net/http"
)

type endpoint struct {
	martini         dto.Martini
	endpointUseCase usecase.Endpoint
}

type Endpoint interface {
	Execute(ctx *gin.Context)
}

func NewEndpoint(martini dto.Martini, endpointUseCase usecase.Endpoint) Endpoint {
	return endpoint{
		martini:         martini,
		endpointUseCase: endpointUseCase,
	}
}

func (e endpoint) Execute(ctx *gin.Context) {
	for _, endpointIndex := range e.martini.Endpoints {
		if (helper.Equals(endpointIndex.Endpoint, ctx.Request.URL.Path) ||
			helper.Equals(endpointIndex.Endpoint, ctx.FullPath())) &&
			helper.Equals(endpointIndex.Method, ctx.Request.Method) {
			executeOutput, err := e.endpointUseCase.Execute(ctx, usecase.ExecuteInput{
				Martini:  e.martini,
				Endpoint: endpointIndex,
				Header:   ctx.Request.Header,
				Query:    ctx.Request.URL.Query(),
				Params:   e.getRequestParams(ctx),
				Body:     e.getRequestBody(ctx),
			})
			if errors.Contains(err, usecase.ErrBadGateway) {
				handler.RespondCodeWithError(ctx, http.StatusBadGateway, err)
			} else if helper.IsNotNil(err) {
				handler.RespondCodeWithError(ctx, http.StatusInternalServerError, err)
			} else {
				e.replyGateway(ctx, executeOutput.StatusCode, executeOutput.Header, executeOutput.Body)
			}
			break
		}
	}
	err := errors.New("endpoint is not response nothing:", ctx.Request.RequestURI, "method:", ctx.Request.Method)
	handler.RespondCodeWithError(ctx, http.StatusInternalServerError, err)
}

func (e endpoint) getRequestParams(ctx *gin.Context) map[string]string {
	result := map[string]string{}
	for _, param := range ctx.Params {
		result[param.Key] = param.Value
	}
	return result
}

func (e endpoint) getRequestBody(ctx *gin.Context) (bodyResult any) {
	bytesBody, err := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bytesBody))
	if helper.IsNil(err) && helper.ContainsIgnoreCase(ctx.GetHeader("Content-Type"), "application/json") {
		bodyResult, _ = oj.Parse(bytesBody)
	} else if helper.IsNotEmpty(bytesBody) {
		bodyResult = string(bytesBody)
	}
	return bodyResult
}

func (e endpoint) replyGateway(ctx *gin.Context, statusCode int, header map[string]string, body any) {
	if ctx.IsAborted() {
		return
	}
	for k, v := range header {
		ctx.Header(k, v)
	}
	if helper.IsNotEmpty(body) {
		handler.RespondCodeWithBody(ctx, statusCode, body)
	} else {
		handler.RespondCode(ctx, statusCode)
	}
}
