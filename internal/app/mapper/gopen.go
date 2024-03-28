package mapper

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"os"
)

func BuildConfigViewDTO(gopenDTO dto.GOpen) dto.ConfigView {
	return dto.ConfigView{
		Version:     os.Getenv("VERSION"),
		VersionDate: os.Getenv("VERSION_DATE"),
		Founder:     os.Getenv("FOUNDER"),
		CodeHelpers: os.Getenv("CODE_HELPERS"),
		Endpoints:   gopenDTO.CountEndpoints(),
		Middlewares: gopenDTO.CountMiddlewares(),
		Backends:    gopenDTO.CountBackends(),
		Modifiers:   gopenDTO.CountModifiers(),
		Config:      gopenDTO,
	}
}
