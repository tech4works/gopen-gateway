package controller

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type static struct {
	gopenVO vo.GOpen
}

type Static interface {
	Ping(ctx *gin.Context)
	Version(ctx *gin.Context)
	Settings(ctx *gin.Context)
}

func NewStatic(gopenVO vo.GOpen) Static {
	return static{
		gopenVO: gopenVO,
	}
}

func (s static) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "%s", "Pong!")
}

func (s static) Version(ctx *gin.Context) {
	if helper.IsNotEmpty(s.gopenVO.Version()) {
		ctx.String(http.StatusOK, "%s", s.gopenVO.Version())
		return
	}
	ctx.Status(http.StatusNotFound)
}

func (s static) Settings(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, mapper.BuildSettingViewDTO(s.gopenVO))
}
