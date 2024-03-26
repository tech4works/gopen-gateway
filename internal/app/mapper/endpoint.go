package mapper

import (
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/util"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
)

func BuildExecuteServiceParams(ctx *gin.Context, gopenVO vo.Gopen, endpointVO vo.Endpoint) (context.Context,
	vo.ExecuteEndpoint) {
	url := util.GetRequestUri(ctx)
	method := ctx.Request.Method
	header := vo.NewHeader(ctx.Request.Header)
	params := vo.NewParams(util.GetRequestParams(ctx))
	query := vo.NewQuery(ctx.Request.URL.Query())
	body := util.GetRequestBody(ctx)

	return ctx.Request.Context(), vo.NewExecuteEndpoint(gopenVO, endpointVO, url, method, header, params, query, body)
}
