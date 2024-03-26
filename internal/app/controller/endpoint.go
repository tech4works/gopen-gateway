package controller

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/util"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type endpoint struct {
	gopen           vo.Gopen
	endpointService service.Endpoint
}

type Endpoint interface {
	Execute(endpointVO vo.Endpoint) gin.HandlerFunc
}

func NewEndpoint(gopenVO vo.Gopen, endpointService service.Endpoint) Endpoint {
	return endpoint{
		gopen:           gopenVO,
		endpointService: endpointService,
	}
}

func (e endpoint) Execute(endpointVO vo.Endpoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// executamos o serviço de dominío para processar o endpoint
		responseVO := e.endpointService.Execute(mapper.BuildExecuteServiceParams(ctx, e.gopen, endpointVO))
		// respondemos o gateway
		e.respondGateway(ctx, responseVO)
	}
}

func (e endpoint) respondGateway(ctx *gin.Context, responseVO vo.Response) {
	// se ja tiver abortado não fazemos nada
	if ctx.IsAborted() {
		return
	}
	for k := range responseVO.Header() {
		ctx.Header(k, responseVO.Header().)
	}
	if helper.IsNotEmpty(responseVO.body) {
		util.RespondCodeWithBody(ctx, responseVO.statusCode, responseVO.body)
	} else {
		util.RespondCode(ctx, responseVO.statusCode)
	}
}
