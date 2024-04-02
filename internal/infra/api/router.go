package api

import (
	"bytes"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *Request)

// Handle handles a request by registering it with the specified engine and endpoint information,
// and then passing it to the provided handler functions.
// The engine parameter is the gin Engine to register the request with.
// The gopenVO parameter is the Gopen configuration object.
// The endpointVO parameter is the Endpoint configuration object.
// The handles parameter is a variadic slice of HandlerFunc functions that will be sequentially called to handle the request.
func Handle(engine *gin.Engine, gopenVO vo.Gopen, endpointVO vo.Endpoint, handles ...HandlerFunc) {
	engine.Handle(endpointVO.Method(), endpointVO.Path(), parseHandles(gopenVO, endpointVO, handles)...)
}

// parseHandles takes a Gopen configuration object, an Endpoint configuration object,
// and a slice of HandlerFunc functions and returns a slice of gin.HandlerFunc functions.
// It iterates over the provided HandlerFuncs and calls the handle function to create
// a gin.HandlerFunc for each one, then appends it to the ginHandler slice.
// Finally, it returns the ginHandler slice.
func parseHandles(gopenVO vo.Gopen, endpointVO vo.Endpoint, handles []HandlerFunc) []gin.HandlerFunc {
	var ginHandler []gin.HandlerFunc
	for _, apiHandler := range handles {
		ginHandler = append(ginHandler, handle(gopenVO, endpointVO, apiHandler))
	}
	return ginHandler
}

// handle handles a request by registering it with the specified Gopen and Endpoint objects,
// and then passing it to the provided HandlerFunc.
// The gopenVO parameter is the Gopen configuration object.
// The endpointVO parameter is the Endpoint configuration object.
// The handle parameter is a HandlerFunc function that will be called to handle the request.
func handle(gopenVO vo.Gopen, endpointVO vo.Endpoint, handle HandlerFunc) gin.HandlerFunc {
	return func(gin *gin.Context) {
		req := buildRequestByContext(gin, gopenVO, endpointVO)
		handle(req)
	}
}

// buildRequestByContext builds a Request object based on the gin context, Gopen configuration, and Endpoint configuration.
// It creates a ResponseWriter and assigns it to the gin context's writer.
// It returns the constructed Request object.
func buildRequestByContext(gin *gin.Context, gopenVO vo.Gopen, endpointVO vo.Endpoint) *Request {
	writer := buildResponseWriter(gin)
	gin.Writer = writer
	return &Request{
		framework: gin,
		gopen:     gopenVO,
		endpoint:  endpointVO,
		writer:    writer,
	}
}

// buildResponseWriter creates a new instance of dto.Writer and initializes its fields.
// The gin parameter is the context of the current HTTP request.
// The returned value is a pointer to the created dto.Writer object.
func buildResponseWriter(gin *gin.Context) *dto.Writer {
	return &dto.Writer{
		Body:           &bytes.Buffer{},
		ResponseWriter: gin.Writer,
	}
}
