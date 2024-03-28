package controller

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type static struct {
	gopenDTO dto.GOpen
}

type Static interface {
	Ping(ctx *gin.Context)
	Version(ctx *gin.Context)
	Config(ctx *gin.Context)
}

func NewStatic(gopenDTO dto.GOpen) Static {
	return static{
		gopenDTO: gopenDTO,
	}
}

func (s static) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "%s", "Pong!")
}

func (s static) Version(ctx *gin.Context) {
	if helper.IsNotEmpty(s.gopenDTO.Version) {
		ctx.String(http.StatusOK, "%s", s.gopenDTO.Version)
		return
	}
	ctx.Status(http.StatusNotFound)
}

func (s static) Config(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, mapper.BuildConfigViewDTO(s.gopenDTO))
}
