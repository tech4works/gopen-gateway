package mapper

import (
	"bytes"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/gin-gonic/gin"
)

func BuildResponseWriter(ctx *gin.Context) *dto.ResponseWriter {
	return &dto.ResponseWriter{
		Body:           &bytes.Buffer{},
		ResponseWriter: ctx.Writer,
	}
}
