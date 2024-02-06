package factory

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/oj"
	"io"
)

type backend struct {
}

type Backend interface {
}

func NewBackend() Backend {
	return backend{}
}

func (b backend) getRequestParams(ctx *gin.Context) map[string]string {
	result := map[string]string{}
	for _, param := range ctx.Params {
		result[param.Key] = param.Value
	}
	return result
}

func (b backend) getRequestBody(ctx *gin.Context) any {
	var bodyResult any
	bytesBody, err := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bytesBody))
	if err == nil {
		if ctx.GetHeader("Content-Type") == "application/json" {
			bodyResult, _ = oj.ParseString(string(bytesBody))
		} else {
			bodyResult = string(bytesBody)
		}
	}
	return bodyResult
}
