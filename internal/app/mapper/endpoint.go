package mapper

import (
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

// BuildExecuteServiceParams builds the parameters needed for executing a service.
// It takes a *api.Request as input and returns a context.Context and a vo.ExecuteEndpoint.
// The context.Context is obtained from the input request.
// The vo.ExecuteEndpoint is created using the Gopen, Endpoint, and Request from the input request.
func BuildExecuteServiceParams(req *api.Request) (context.Context, vo.ExecuteEndpoint) {
	return req.Context(), vo.NewExecuteEndpoint(req.Gopen(), req.Endpoint(), req.Request())
}
