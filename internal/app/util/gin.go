package util

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
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

	respondCodeByEncode(ctx, encode, code, body)

	ctx.Abort()
}

func RespondCodeWithError(ctx *gin.Context, encode enum.ResponseEncode, code int, err error) {
	if ctx.IsAborted() {
		return
	}

	respondCodeByEncode(ctx, encode, code, buildErrorViewDTO(ctx.Request.URL.String(), err))

	ctx.Abort()
}

func respondCodeByEncode(ctx *gin.Context, encode enum.ResponseEncode, code int, body any) {
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
		ctx.YAML(code, body)
		break
	default:
		if helper.IsJsonType(body) {
			ctx.JSON(code, body)
		} else {
			ctx.String(code, "%s", body)
		}
		break
	}
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
