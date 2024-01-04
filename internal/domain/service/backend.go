package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"slices"
	"strings"
)

type backend struct {
}

type Backend interface {
	Execute(ctx *gin.Context, backend dto.Backend, requests *[]dto.BackendRequest) (*http.Response, error)
}

func NewBackend() Backend {
	return backend{}
}

func (b backend) Execute(ctx *gin.Context, backend dto.Backend, requests *[]dto.BackendRequest) (
	*http.Response, error) {
	request, err := b.prepareBackendRequest(ctx, backend, requests)
	if err != nil {
		return nil, err
	}
	return b.makeRequest(ctx, backend, *request)
}

func (b backend) prepareBackendRequest(ctx *gin.Context, backend dto.Backend, requests *[]dto.BackendRequest) (
	*dto.BackendRequest, error) {
	request := (*requests)[len(*requests)-1]
	//removemos todos os headers nao mapeados no backend.forward-headers
	if !slices.Contains(backend.ForwardHeaders, "*") {
		for key := range request.Header {
			if !slices.Contains(backend.ForwardHeaders, key) {
				request.Header.Del(key)
			}
		}
	}
	//removemos os queryParams nao mapeados no backend
	for key := range request.Query {
		if !slices.Contains(backend.Query, key) {
			request.Query.Del(key)
		}
	}
	//substituímos os parâmetros no endpoint pelo valor do parâmetro por exemplo find/user/:userID para find/user/2
	requestURL := request.Endpoint
	for key, value := range request.Params {
		if strings.Contains(request.Endpoint, ":"+key) {
			requestURL = strings.ReplaceAll(request.Endpoint, ":"+key, value)
		}
	}
	request.Endpoint = requestURL
	//preparamos o body para envio
	if request.Body != nil {
		var bodyBytes []byte
		if ctx.GetHeader("Content-Type") == "application/json" {
			bodyBytes, _ = json.Marshal(request.Body)
		} else {
			bodyBytes = []byte(fmt.Sprintf("%v", request.Body))
		}
		request.BodyToSend = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	(*requests)[len(*requests)-1] = request
	return &request, nil
}

func (b backend) makeRequest(ctx *gin.Context, backend dto.Backend, currentRequest dto.BackendRequest) (
	*http.Response, error) {
	//para cada host uma request: todo -> talvez desenvolver uma opção de balancer através do host
	host := backend.Host[rand.Intn(len(backend.Host))]
	client := &http.Client{}
	request, err := http.NewRequestWithContext(
		ctx.Request.Context(),
		backend.Method,
		host+currentRequest.Endpoint,
		currentRequest.BodyToSend,
	)
	if err != nil {
		return nil, err
	}
	request.Header = currentRequest.Header
	request.URL.RawQuery = currentRequest.Query.Encode()
	return client.Do(request)
}
