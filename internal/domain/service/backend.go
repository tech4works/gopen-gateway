package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/martini-gateway/internal/domain/model/valueobject"
	"github.com/iancoleman/orderedmap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type BackendRequest struct {
	Host     string
	Endpoint string
	Url      string
	Method   string
	Header   http.Header
	Query    url.Values
	Params   map[string]string
	Body     any
}

type BackendResponse struct {
	StatusCode int
	Header     http.Header
	Body       any
	Group      string
	Hide       bool
}

type BuildBackendRequestInput struct {
	Backend valueobject.Backend
	Host    string
	Header  http.Header
	Query   url.Values
	Params  map[string]string
	Body    any
}

type ExecuteBackendInput struct {
	Backend        valueobject.Backend
	BackendRequest BackendRequest
}

type backend struct {
}

type Backend interface {
	BuildBackendRequest(input BuildBackendRequestInput) (*BackendRequest, error)
	Execute(ctx context.Context, input ExecuteBackendInput) (*BackendResponse, error)
}

func NewBackend() Backend {
	return backend{}
}

func (b backend) BuildBackendRequest(input BuildBackendRequestInput) (*BackendRequest, error) {
	//removemos todos os headers nao mapeados no backend.forward-headers
	if helper.NotContains(input.Backend.ForwardHeaders, "*") {
		for key := range input.Header {
			if helper.NotContains(input.Backend.ForwardHeaders, key) {
				input.Header.Del(key)
			}
		}
	}
	//removemos os queryParams nao mapeados no backend
	for key := range input.Query {
		if helper.NotContains(input.Backend.Query, key) {
			input.Query.Del(key)
		}
	}
	//substituímos os parâmetros no endpoint pelo valor do parâmetro por exemplo find/user/:userID para find/user/2
	endpoint := input.Backend.Endpoint
	for key, value := range input.Params {
		rKey := fmt.Sprint(":", key)
		if helper.Contains(input.Backend.Endpoint, rKey) {
			endpoint = strings.ReplaceAll(input.Backend.Endpoint, rKey, value)
		}
	}
	return &BackendRequest{
		Host:     input.Host,
		Endpoint: input.Backend.Endpoint,
		Url:      fmt.Sprint(input.Host, endpoint),
		Method:   input.Backend.Method,
		Header:   input.Header,
		Query:    input.Query,
		Params:   input.Params,
		Body:     input.Body,
	}, nil
}

func (b backend) Execute(ctx context.Context, input ExecuteBackendInput) (*BackendResponse, error) {
	response, err := b.makeRequest(ctx, input.BackendRequest)
	if helper.IsNotNil(err) {
		return nil, err
	}
	defer b.closeBodyResponse(response)
	return &BackendResponse{
		Header:     response.Header,
		StatusCode: response.StatusCode,
		Body:       b.readResponseBody(response),
		Group:      input.Backend.Group,
		Hide:       input.Backend.HideResponse,
	}, nil
}

func (b backend) makeRequest(ctx context.Context, backendRequest BackendRequest) (*http.Response, error) {
	var bodyToSend io.ReadCloser
	if helper.IsNotNil(backendRequest.Body) {
		bodyBytes, err := helper.ConvertToBytes(backendRequest.Body)
		if helper.IsNotNil(err) {
			return nil, err
		}
		bodyToSend = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	client := &http.Client{}
	request, err := http.NewRequestWithContext(
		ctx,
		backendRequest.Method,
		backendRequest.Url,
		bodyToSend,
	)
	if helper.IsNotNil(err) {
		return nil, err
	}
	request.Header = backendRequest.Header
	request.URL.RawQuery = backendRequest.Query.Encode()
	return client.Do(request)
}

func (b backend) closeBodyResponse(response *http.Response) {
	_ = response.Body.Close()
}

func (b backend) readResponseBody(resp *http.Response) (result any) {
	bodyBytes, _ := io.ReadAll(resp.Body)
	if helper.IsStringMap(bodyBytes) {
		orderedMap := orderedmap.New()
		helper.SimpleConvertToDest(bodyBytes, orderedMap)
		result = orderedMap
	} else if helper.IsStringSlice(bodyBytes) {
		var sliceOrderedMap []*orderedmap.OrderedMap
		helper.SimpleConvertToDest(bodyBytes, &sliceOrderedMap)
		result = sliceOrderedMap
	} else if helper.IsNotEmpty(bodyBytes) {
		bodyString := string(bodyBytes)
		result = &bodyString
	}
	return result
}
