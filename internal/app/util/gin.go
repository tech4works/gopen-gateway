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

func RespondGateway(ctx *gin.Context, responseVO vo.Response) {
	// se ja tiver abortado não fazemos nada
	if ctx.IsAborted() {
		return
	}

	// iteramos o header para responder o mesmo
	setHeaderResponse(ctx, responseVO.Header())

	statusCode := responseVO.StatusCode()
	encode := responseVO.Endpoint().ResponseEncode()
	body := responseVO.Body()

	// verificamos se tem valor o body
	if body.IsNotEmpty() {
		respondCodeWithBody(ctx, encode, statusCode, responseVO.Header(), responseVO.Body())
	} else if helper.IsNotNil(responseVO.Err()) {
		respondCodeWithError(ctx, encode, statusCode, responseVO.Header(), responseVO.Err())
	} else {
		respondCode(ctx, statusCode)
	}
}

func RespondGatewayError(ctx *gin.Context, encode enum.ResponseEncode, code int, err error) {
	// se ja tiver abortado não fazemos nada
	if ctx.IsAborted() {
		return
	}

	// responde o erro com o header padrão
	respondCodeWithError(ctx, encode, code, vo.NewHeaderFailed(), err)
}

func respondCode(ctx *gin.Context, code int) {
	// se ja tiver abortado não fazemos nada
	if ctx.IsAborted() {
		return
	}

	// setamos o código http recebido
	ctx.Status(code)

	// abortamos a requisição
	ctx.Abort()
}

func respondCodeWithBody(ctx *gin.Context, encode enum.ResponseEncode, code int, header vo.Header, body vo.Body) {
	// se ja tiver abortado não fazemos nada
	if ctx.IsAborted() {
		return
	}

	// chamamos o response principal
	respondCodeByEncode(ctx, encode, code, header, body.Value())

	// abortamos a requisição
	ctx.Abort()
}

func respondCodeWithError(ctx *gin.Context, encode enum.ResponseEncode, code int, header vo.Header, err error) {
	// se ja tiver abortado não fazemos nada
	if ctx.IsAborted() {
		return
	}

	// chamamos o response principal
	respondCodeByEncode(ctx, encode, code, header, buildErrorViewDTO(ctx.Request.URL.String(), err))

	// abortamos a requisição
	ctx.Abort()
}

func respondCodeByEncode(ctx *gin.Context, encode enum.ResponseEncode, code int, header vo.Header, body any) {
	// iteramos o header para responder o mesmo
	setHeaderResponse(ctx, header)

	// respondemos o body a partir do encode configurado
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

func GetResponseWriter(ctx *gin.Context) dto.ResponseWriter {
	// obtemos do buff do contexto
	var responseWriter dto.ResponseWriter
	responseWriterByContext, ok := ctx.Get("writer")
	if ok {
		responseWriter = *responseWriterByContext.(*dto.ResponseWriter)
	}
	return responseWriter
}

func setHeaderResponse(ctx *gin.Context, header vo.Header) {
	for key := range header {
		if helper.EqualsIgnoreCase(key, "Content-Length") || helper.EqualsIgnoreCase(key, "Content-Type") ||
			helper.EqualsIgnoreCase(key, "Date") {
			continue
		}
		ctx.Header(key, header.Get(key))
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
