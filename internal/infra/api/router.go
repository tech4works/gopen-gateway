package api

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"sync"
)

// HandlerFunc is a type for defining functions that can be used as HTTP route handlers.
// It takes a *Context parameter which contains information about the current request
// and allows the handler to read and manipulate the request and response.
// The handler function must not return any values.
type HandlerFunc func(ctx *Context)

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
		// verificamos se esse contexto ja foi construído
		ctx, ok := gin.Get("context")
		if !ok {
			// construímos o contexto da requisição através dos objetos de valores e o gin todo: ve se isso tem impacto
			ctx = buildContext(gin, gopenVO, endpointVO)
			// setamos o contexto criado da requisição
			gin.Set("context", ctx)
		}
		// chamamos a função persistida
		handle(ctx.(*Context))
	}
}

// buildContext builds a Context object based on the gin context, Gopen configuration, and Endpoint configuration.
// It creates a ResponseWriter and assigns it to the gin context's writer.
// It returns the constructed Context object.
func buildContext(gin *gin.Context, gopenVO vo.Gopen, endpointVO vo.Endpoint) *Context {
	// o contexto da requisição é criado
	return &Context{
		mutex:     &sync.RWMutex{},
		framework: gin,
		gopen:     gopenVO,
		endpoint:  endpointVO,
		request:   vo.NewRequest(gin),
		response:  vo.NewResponse(endpointVO),
	}
}
