package util

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"time"
)

func RespondCode(ctx *gin.Context, code int) {
	if ctx.IsAborted() {
		return
	}
	ctx.Status(code)
	ctx.Abort()
}

func RespondCodeWithBody(ctx *gin.Context, encode enum.ResponseEncode, code int, body vo.Body) {
	if ctx.IsAborted() {
		return
	}

	switch encode {
	case enum.ResponseEncodeText:
		ctx.String(code, "%s", body)
		break
	case enum.ResponseEncodeJson:
		ctx.JSON(code, body)
		break
	case enum.ResponseEncodeXml:
		ctx.XML(code, body)
		break
	case enum.ResponseEncodeYaml:
		v, _ := yaml.Marshal(body.Interface())
		ctx.YAML(code, string(v))
		break
	default:
		if helper.IsJsonType(body) {
			ctx.JSON(code, body)
		} else {
			ctx.String(code, "%s", body)
		}
		break
	}

	ctx.Abort()
}

func RespondCodeWithError(ctx *gin.Context, code int, err error) {
	if ctx.IsAborted() {
		return
	}
	ctx.JSON(code, buildErrorViewDTO(ctx.Request.URL.String(), err))
	ctx.Abort()
}

func buildErrorViewDTO(requestUrl string, err error) dto.ErrorView {
	errDetails := errors.Details(err)
	return dto.ErrorView{
		File:      errDetails.GetFile(),
		Line:      errDetails.GetLine(),
		Endpoint:  requestUrl,
		Message:   errDetails.GetMessage(),
		Timestamp: time.Now(),
	}
}
