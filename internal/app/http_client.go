package app

import (
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type HTTPClient interface {
	MakeRequest(ctx context.Context, endpoint *vo.Endpoint, request *vo.HTTPBackendRequest) *vo.HTTPBackendResponse
}
