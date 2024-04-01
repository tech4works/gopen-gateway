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

// Ping is a method that handles the "Ping" request.
// It responds with the string "Pong!" and a status code of 200 (OK).
func (s static) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "%s", "Pong!")
}

// Version is a method that handles the "Version" request.
// It checks if the version is not empty and responds with the version string and a status code of 200 (OK).
// If the version is empty, it responds with a status code of 404 (Not Found).
func (s static) Version(ctx *gin.Context) {
	if helper.IsNotEmpty(s.gopenVO.Version()) {
		ctx.String(http.StatusOK, "%s", s.gopenVO.Version())
		return
	}
	ctx.Status(http.StatusNotFound)
}

// Settings is a method that handles the "Settings" request.
// It retrieves the necessary data from the gopenVO object and constructs a dto.SettingView object.
// The SettingView object contains information such as version, version date, founder, code helpers,
// number of endpoints, number of middlewares, number of backends, number of modifiers, and the gopenVO object itself.
// The method then responds with the constructed SettingView object in JSON format and a status code of 200 (OK).
func (s static) Settings(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, mapper.BuildSettingViewDTO(s.gopenVO))
}
