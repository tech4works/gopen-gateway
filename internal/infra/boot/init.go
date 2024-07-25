package boot

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"os"
)

func Init() string {
	if helper.IsLessThanOrEqual(os.Args, 1) {
		panic("Please enter ENV as second argument! ex: dev, prd")
	}

	err := loadDefaultEnvs()
	if helper.IsNotNil(err) {
		panic(err)
	}

	jaegerHost := os.Getenv("JAEGER_HOST")
	if helper.IsNotEmpty(jaegerHost) {
		PrintInfo("Booting Jaeger...")
		jaeger, closer, err := InitJaeger(jaegerHost)
		if helper.IsNotNil(err) {
			PrintWarn(err)
		} else {
			defer closer.Close()
			opentracing.SetGlobalTracer(jaeger)
		}
	}

	return os.Args[1]
}

func loadDefaultEnvs() (err error) {
	if err = godotenv.Load("./.env"); helper.IsNotNil(err) {
		err = errors.New("Error load Gopen envs default:", err)
	}
	return err
}
