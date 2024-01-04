package factory

import (
	"bytes"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/oj"
	"io"
	"net/http"
)

type backend struct {
}

type Backend interface {
	CreateBackendRequest(ctx *gin.Context, backend dto.Backend, requests *[]dto.BackendRequest, isAuthBackend bool)
	CreateBackendResponse(responseHttp *http.Response, backend dto.Backend, responses *[]dto.BackendResponse) error
}

func NewBackend() Backend {
	return backend{}
}

func (b backend) CreateBackendRequest(
	ctx *gin.Context,
	backend dto.Backend,
	requests *[]dto.BackendRequest,
	isAuthBackend bool,
) {
	ctx.Request.Header.Set("X-Forwarded-For", ctx.ClientIP())
	if isAuthBackend {
		*requests = append(*requests, dto.BackendRequest{
			Endpoint: backend.Endpoint,
			Header:   ctx.Request.Header.Clone(),
		})
	} else {
		*requests = append(*requests, dto.BackendRequest{
			Endpoint: backend.Endpoint,
			Header:   ctx.Request.Header.Clone(),
			Query:    ctx.Request.URL.Query(),
			Params:   b.getRequestParams(ctx),
			Body:     b.getRequestBody(ctx),
		})
	}
}

func (b backend) CreateBackendResponse(responseHttp *http.Response, backend dto.Backend, responses *[]dto.BackendResponse,
) error {
	bodyParsed, err := b.parseHttpResponse(responseHttp)
	if err != nil {
		return err
	}
	resp := dto.BackendResponse{
		Header:     responseHttp.Header,
		StatusCode: responseHttp.StatusCode,
		Body:       bodyParsed,
		Group:      backend.Group,
		Remove:     backend.RemoveResponse,
	}
	*responses = append(*responses, resp)
	return nil
}

func (b backend) parseHttpResponse(resp *http.Response) (any, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(bodyBytes) > 0 {
		responseJSON, err := oj.ParseString(string(bodyBytes))
		if err != nil {
			return string(bodyBytes), nil
		}
		return responseJSON, nil
	}
	return nil, nil
}

func (b backend) getRequestParams(ctx *gin.Context) map[string]string {
	result := map[string]string{}
	for _, param := range ctx.Params {
		result[param.Key] = param.Value
	}
	return result
}

func (b backend) getRequestBody(ctx *gin.Context) any {
	var bodyResult any
	bytesBody, err := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bytesBody))
	if err == nil {
		if ctx.GetHeader("Content-Type") == "application/json" {
			bodyResult, _ = oj.ParseString(string(bytesBody))
		} else {
			bodyResult = string(bytesBody)
		}
	}
	return bodyResult
}
