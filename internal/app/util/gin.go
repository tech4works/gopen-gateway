package util

import (
	"bytes"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/gin-gonic/gin"
	"io"
)

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

func GetRequestUri(ctx *gin.Context) (uri string) {
	uri = ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery
	if helper.IsNotEmpty(raw) {
		uri = fmt.Sprint(uri, "?", raw)
	}
	return uri
}

func GetRequestParams(ctx *gin.Context) map[string]string {
	result := map[string]string{}
	for _, param := range ctx.Params {
		result[param.Key] = param.Value
	}
	return result
}

func GetRequestBody(ctx *gin.Context) any {
	bytesBody, _ := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bytesBody))

	// se vazio, retornamos nil
	if helper.IsEmpty(bytesBody) {
		return nil
	}

	// se ele for json, verificamos se body é um map ou slice para manter ordenado
	if helper.ContainsIgnoreCase(ctx.GetHeader("Content-Type"), "application/json") {
		// convertemos os bytes do body em uma interface de objeto
		var dest any
		helper.SimpleConvertToDest(bytesBody, dest)
		return dest
	}
	//todo: futuramente podemos trabalhar com o XML e o FORM-DATA com o modifier e envio

	// no pior das hipóteses retornamos uma string do body
	return string(bytesBody)
}
