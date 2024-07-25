package boot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/joho/godotenv"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"regexp"
	"strings"
	"time"
)

const runtimeFolder = "./runtime"
const jsonRuntimeUri = runtimeFolder + "/.json"
const jsonSchemaUri = "./json-schema.json"

func Start(env string) {
	PrintInfo("Loading Gopen envs...")
	err := loadEnvs(env)
	if helper.IsNotNil(err) {
		PrintWarn(err)
	}

	PrintInfo("Loading Gopen json...")
	gopenDTO, err := loadJson(env)
	if helper.IsNotNil(err) {
		panic(err)
	}

	PrintInfo("Configuring cache store...")
	store := NewMemoryStore()
	if helper.IsNotNil(gopenDTO.Store) {
		store = NewRedisStore(gopenDTO.Store.Redis.Address, gopenDTO.Store.Redis.Password)
	}
	defer store.Close()

	err = writeRuntimeJson(gopenDTO)
	if helper.IsNotNil(err) {
		PrintWarn(err)
	}

	PrintInfo("Building application...")
	gopenAPP := app.NewGopen(gopenDTO, store)

	if gopenDTO.HotReload {
		PrintInfo("Configuring watcher...")
		watcher, err := NewWatcher(env, restart(env, gopenAPP))
		if helper.IsNotNil(err) {
			PrintWarn("Error configure watcher:", err)
		} else {
			defer watcher.Close()
		}
	}

	PrintInfo("Starting application...")
	gopenAPP.ListerAndServer()
}

func Stop() {
	fmt.Println()
	err := removeRuntimeJson()
	if helper.IsNotNil(err) {
		PrintWarn("Error to remove runtime json!")
	}
	PrintTitle("STOPPED")
}

func restart(env string, oldGopenAPP app.Gopen) func() {
	return func() {
		defer func() {
			if r := recover(); helper.IsNotNil(r) {
				errorDetails := errors.Details(r.(error))
				PrintError("Error restart server:", errorDetails.GetCause())

				recovery(oldGopenAPP)
			}
		}()

		fmt.Println()
		fmt.Println()
		PrintTitle("RESTART")

		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()

		PrintInfo("Shutting down current server...")
		err := oldGopenAPP.Shutdown(ctx)
		if helper.IsNotNil(err) {
			PrintWarnf("Error shutdown current server: %s!", err)
			return
		}

		go Start(env)
	}
}

func recovery(oldGopenAPP app.Gopen) {
	fmt.Println()
	PrintTitle("RECOVERY")

	go oldGopenAPP.ListerAndServer()
}

func loadEnvs(env string) (err error) {
	gopenEnvUri := buildEnvUri(env)

	if err = godotenv.Overload(gopenEnvUri); helper.IsNotNil(err) {
		err = errors.New("Error load Gopen envs from uri:", gopenEnvUri, "err:", err)
	}

	return err
}

func loadJson(env string) (*dto.Gopen, error) {
	gopenJsonUri := buildJsonUri(env)

	gopenJsonBytes, err := os.ReadFile(gopenJsonUri)
	if helper.IsNotNil(err) {
		return nil, errors.New("Error read Gopen config from file json:", gopenJsonUri, "err:", err)
	}
	gopenJsonBytes = fillEnvValues(gopenJsonBytes)

	if err = validateJsonBySchema(jsonSchemaUri, gopenJsonBytes); helper.IsNotNil(err) {
		return nil, err
	}

	var gopenDTO dto.Gopen
	err = helper.ConvertToDest(gopenJsonBytes, &gopenDTO)
	if helper.IsNotNil(err) {
		return nil, err
	}

	return &gopenDTO, nil
}

func fillEnvValues(gopenJsonBytes []byte) []byte {
	// todo: aceitar campos não string receber variável de ambiente também
	//  foi pensado que talvez utilizar campos string e any para isso, convertendo para o tipo desejado apenas
	//  quando objeto de valor for montado
	gopenJsonStr := helper.SimpleConvertToString(gopenJsonBytes)

	regex := regexp.MustCompile(`\$\w+`)
	words := regex.FindAllString(gopenJsonStr, -1)

	count := 0
	for _, word := range words {
		envKey := strings.ReplaceAll(word, "$", "")
		envValue := os.Getenv(envKey)
		if helper.IsNotEmpty(envValue) {
			gopenJsonStr = strings.ReplaceAll(gopenJsonStr, word, envValue)
			count++
		}
	}

	return helper.SimpleConvertToBytes(gopenJsonStr)
}

func validateJsonBySchema(jsonSchemaUri string, jsonBytes []byte) error {
	schemaLoader := gojsonschema.NewReferenceLoader(jsonSchemaUri)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if helper.IsNotNil(err) {
		return errors.New("Error validate schema:", err)
	} else if !result.Valid() {
		errorMsg := fmt.Sprintf("Map poorly formatted!\n")
		for _, desc := range result.Errors() {
			errorMsg += fmt.Sprintf("- %s\n", desc)
		}
		err = errors.New(errorMsg)
	}

	return err
}

func writeRuntimeJson(gopenDTO *dto.Gopen) error {
	if _, err := os.Stat(runtimeFolder); os.IsNotExist(err) {
		err = os.MkdirAll(runtimeFolder, 0755)
		if helper.IsNotNil(err) {
			return err
		}
	}

	gopenJsonBytes, err := json.MarshalIndent(gopenDTO, "", "\t")
	if helper.IsNil(err) {
		err = os.WriteFile(jsonRuntimeUri, gopenJsonBytes, 0644)
	}

	return err
}

func removeRuntimeJson() error {
	err := os.Remove(runtimeFolder)
	if errors.IsNot(err, os.ErrNotExist) {
		err = nil
	}
	return err
}

func buildEnvUri(env string) string {
	return fmt.Sprintf("./gopen/%s/.env", env)
}

func buildJsonUri(env string) string {
	return fmt.Sprintf("./gopen/%s/.json", env)
}
