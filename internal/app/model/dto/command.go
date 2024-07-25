package dto

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type ExecuteEndpoint struct {
	Gopen    *vo.Gopen
	Endpoint *vo.Endpoint
	Request  *vo.HTTPRequest
}
