package factory

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"os"
)

func BuildSettingView(gopen dto.Gopen) dto.SettingView {
	copied := gopen
	copied.Store = nil

	return dto.SettingView{
		Version:      os.Getenv("VERSION"),
		VersionDate:  os.Getenv("VERSION_DATE"),
		Founder:      os.Getenv("FOUNDER"),
		Contributors: helper.SimpleConvertToInt(os.Getenv("CONTRIBUTORS")),
		Endpoints:    countEndpoints(gopen),
		Middlewares:  countMiddlewares(gopen),
		Backends:     countBackends(gopen),
		Setting:      copied,
	}
}

func countEndpoints(gopen dto.Gopen) int {
	return len(gopen.Endpoints)
}

func countMiddlewares(gopen dto.Gopen) int {
	if helper.IsNotNil(gopen.Middlewares) {
		return len(gopen.Middlewares)
	}
	return 0
}

func countBackends(gopen dto.Gopen) (count int) {
	for _, endpoint := range gopen.Endpoints {
		count += len(endpoint.Beforewares) + len(endpoint.Backends) + len(endpoint.Afterwares)
	}
	return count
}
