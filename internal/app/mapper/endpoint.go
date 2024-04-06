package mapper

import (
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

// BuildExecuteServiceParams builds the parameters needed for executing a service.
// It takes a *api.Context as input and returns a context.Context and a vo.ExecuteEndpoint.
// The context.Context is obtained from the input request.
// The vo.ExecuteEndpoint is created using the Gopen, Endpoint, and Request from the input request.
func BuildExecuteServiceParams(ctx *api.Context) (context.Context, vo.ExecuteEndpoint) {
	return ctx.Context(), vo.NewExecuteEndpoint(ctx.Gopen(), ctx.Endpoint(), ctx.Request())
}
