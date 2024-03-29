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
		e.respondGateway(ctx, endpointVO, responseVO)
	}
}

func (e endpoint) respondGateway(ctx *gin.Context, endpointVO vo.Endpoint, responseVO vo.Response) {
	// se ja tiver abortado não fazemos nada
	if ctx.IsAborted() {
		return
	}

	// iteramos o header para responder o mesmo
	for key := range responseVO.Header() {
		if helper.EqualsIgnoreCase(key, "Content-Length") || helper.EqualsIgnoreCase(key, "Content-Type") ||
			helper.EqualsIgnoreCase(key, "Date") {
			continue
		}
		ctx.Header(key, responseVO.Header().Get(key))
	}

	// verificamos se tem valor o body
	if responseVO.Body().IsNotEmpty() {
		util.RespondCodeWithBody(ctx, endpointVO.ResponseEncode(), responseVO.StatusCode(), responseVO.Body())
	} else if helper.IsNotNil(responseVO.Err()) {
		util.RespondCodeWithError(ctx, endpointVO.ResponseEncode(), responseVO.StatusCode(), responseVO.Err())
	} else {
		util.RespondCode(ctx, responseVO.StatusCode())
	}
}
