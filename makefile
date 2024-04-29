.SILENT:

UNAME_S := $(shell uname -s 2>/dev/null)

ifeq ($(OS),Windows_NT)
	READPORT := ./bin/readport-windows-amd64.exe
else ifeq ($(UNAME_S),Linux)
	READPORT := ./bin/readport-linux-amd64
else ifeq ($(UNAME_S),Darwin)
	READPORT := ./bin/readport-darwin-amd64
else
$(error Unknown operating system: $(UNAME_S))
endif

build-readport: build-readport-linux build-readport-darwin build-readport-windows

build-readport-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/readport-linux-amd64 ./cmd/readport.go

build-readport-darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/readport-darwin-amd64 ./cmd/readport.go

build-readport-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/readport-windows-amd64.exe ./cmd/readport.go

run:
	echo "Getting env name argument..."

	# Obtemos o nome do ambiente passado no primeiro argumento
	$(eval ENV_NAME := $(filter-out $@,$(MAKECMDGOALS)))
	@if [ -z "$(ENV_NAME)" ]; then \
		echo "Error: No argument provided. Usage: make run {env name}"; \
		exit 1; \
	fi

	echo "Checking if json ./gopen/$(ENV_NAME)/.json exists..."

	# Verificamos se o json de configuração existe pelo nome do ambiente passado
	@if [ ! -f "./gopen/$(ENV_NAME)/.json" ]; then \
		echo "Error: File ./gopen/$(ENV_NAME)/.json does not exist"; \
		exit 1; \
	fi

	echo "Checking you operation system..."
	echo "Operating System: $(UNAME_S)"
	echo "Getting port from json ./gopen/$(ENV_NAME)/.json..."

	$(eval PORT := $(shell ./$(READPORT) ./gopen/$(ENV_NAME)/.json))

	echo "Starting docker with ENV_NAME=$(ENV_NAME) and PORT=$(PORT)..."

	ENV_NAME=$(ENV_NAME) PORT=$(PORT)  docker-compose up
%:
	@: