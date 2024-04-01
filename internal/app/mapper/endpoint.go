package mapper

import (
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

func BuildExecuteServiceParams(req *api.Request) (context.Context, vo.ExecuteEndpoint) {
	return req.Context(), vo.NewExecuteEndpoint(req.GOpen(), req.Endpoint(), req.Request())
}
