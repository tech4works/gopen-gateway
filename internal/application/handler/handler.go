package handler

import (
	"bytes"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
	Message  string `json:"message,omitempty"`
}

type ResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (r ResponseWriter) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r ResponseWriter) WriteString(s string) (n int, err error) {
	r.Body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}

func RespondCode(ctx *gin.Context, code int) {
	if ctx.IsAborted() {
		return
	}
	ctx.Status(code)
	ctx.Abort()
}

func RespondCodeWithBody(ctx *gin.Context, code int, body any) {
	if ctx.IsAborted() {
		return
	}
	if helper.IsJson(body) {
		ctx.JSON(code, body)
	} else {
		ctx.String(code, "%s", body)
	}
	ctx.Abort()
}

func RespondCodeWithError(ctx *gin.Context, code int, err error) {
	if ctx.IsAborted() {
		return
	}
	ctx.JSON(code, prepareErrorResponseByErr(ctx.Request.RequestURI, err))
	ctx.Abort()
}

func prepareErrorResponseByErr(endpoint string, err error) errorResponse {
	detailsErr := errors.Details(err)
	return errorResponse{
		File:     detailsErr.GetFile(),
		Line:     detailsErr.GetLine(),
		Endpoint: endpoint,
		Message:  detailsErr.GetMessage(),
	}
}
