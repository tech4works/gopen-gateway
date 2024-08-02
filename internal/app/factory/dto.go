package factory

import (
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
)

func BuildExecuteEndpoint(ctx app.Context) dto.ExecuteEndpoint {
	return dto.ExecuteEndpoint{
		Gopen:    ctx.Gopen(),
		Endpoint: ctx.Endpoint(),
		Request:  ctx.Request(),
	}
}
