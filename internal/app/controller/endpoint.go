package controller

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/util"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type endpoint struct {
	gopen           vo.GOpen
	endpointService service.Endpoint
}

type Endpoint interface {
	Execute(endpointVO vo.Endpoint) gin.HandlerFunc
}

func NewEndpoint(gopenVO vo.GOpen, endpointService service.Endpoint) Endpoint {
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
		util.RespondGateway(ctx, responseVO)
	}
}
