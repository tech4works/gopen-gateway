package middleware

import (
	"bytes"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/gin-gonic/gin"
)

type writer struct {
}

type Writer interface {
	Do(ctx *gin.Context)
}

func NewWriter() Writer {
	return writer{}
}

func (w writer) Do(ctx *gin.Context) {
	// inicializamos o writer de resposta
	responseWriter := &dto.ResponseWriter{
		Body:           &bytes.Buffer{},
		ResponseWriter: ctx.Writer,
	}
	ctx.Writer = responseWriter

	// gravamos no contexto
	ctx.Set("writer", responseWriter)

	// damos próximo
	ctx.Next()
}
