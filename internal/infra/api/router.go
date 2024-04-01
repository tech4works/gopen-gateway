package api

import (
	"bytes"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerFunc func(ctx *Request)

func Handle(engine *gin.Engine, gopenVO vo.GOpen, endpointVO vo.Endpoint, handles ...HandlerFunc) {
	engine.Handle(endpointVO.Method(), endpointVO.Path(), parseHandles(gopenVO, endpointVO, handles)...)
}

func parseHandles(gopenVO vo.GOpen, endpointVO vo.Endpoint, handles []HandlerFunc) []gin.HandlerFunc {
	var ginHandler []gin.HandlerFunc
	for _, apiHandler := range handles {
		ginHandler = append(ginHandler, handle(gopenVO, endpointVO, apiHandler))
	}
	return ginHandler
}

func handle(gopenVO vo.GOpen, endpointVO vo.Endpoint, handle HandlerFunc) gin.HandlerFunc {
	return func(gin *gin.Context) {
		req, err := buildRequestByContext(gin, gopenVO, endpointVO)
		if helper.IsNotNil(err) {
			detailsErr := errors.Details(err)
			gin.JSON(http.StatusBadRequest, buildErrorResponse(gin.Request.URL.String(), detailsErr))
			gin.Abort()
			return
		}
		handle(req)
	}
}

func buildRequestByContext(gin *gin.Context, gopenVO vo.GOpen, endpointVO vo.Endpoint) (*Request, error) {
	writer := buildResponseWriter(gin)
	gin.Writer = writer
	request := &Request{
		framework: gin,
		gopen:     gopenVO,
		endpoint:  endpointVO,
		writer:    writer,
	}
	return request, nil
}

func buildResponseWriter(gin *gin.Context) *dto.Writer {
	return &dto.Writer{
		Body:           &bytes.Buffer{},
		ResponseWriter: gin.Writer,
	}
}
