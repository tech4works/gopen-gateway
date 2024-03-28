package mapper

import (
	"bytes"
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"io"
)

func BuildExecuteServiceParams(ctx *gin.Context, gopenVO vo.GOpen, endpointVO vo.Endpoint) (context.Context,
	vo.ExecuteEndpoint) {
	url := ctx.Request.URL.String()
	method := ctx.Request.Method
	header := vo.NewHeader(ctx.Request.Header)
	params := vo.NewParams(getRequestParams(ctx))
	query := vo.NewQuery(ctx.Request.URL.Query())
	body := getRequestBody(ctx)

	return ctx.Request.Context(), vo.NewExecuteEndpoint(gopenVO, endpointVO, url, method, header, params, query, body)
}

func getRequestParams(ctx *gin.Context) map[string]string {
	result := map[string]string{}
	for _, param := range ctx.Params {
		result[param.Key] = param.Value
	}
	return result
}

func getRequestBody(ctx *gin.Context) vo.Body {
	bytesBody, _ := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bytesBody))

	// no pior das hipóteses retornamos uma string do body
	return vo.NewBodyByContentType(ctx.GetHeader("Content-Type"), bytesBody)
}
