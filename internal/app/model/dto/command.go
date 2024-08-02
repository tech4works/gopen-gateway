package dto

import (
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type ExecuteEndpoint struct {
	TraceID  string
	ClientIP string
	Gopen    *vo.Gopen
	Endpoint *vo.Endpoint
	Request  *vo.HTTPRequest
}
