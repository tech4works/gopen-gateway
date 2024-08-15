.SILENT:

# Get the name of the operating system
UNAME_S := $(shell uname -s 2>/dev/null)

# Windows we obtain it through the OS variable
ifeq ($(OS),Windows_NT)
	# Get the port reading bin for Windows
	READPORT := ./bin/readport-windows-amd64.exe
else ifeq ($(UNAME_S),Linux)
	# Get the port reading bin for Linux
	READPORT := ./bin/readport-linux-amd64
else ifeq ($(UNAME_S),Darwin)
	# Get the port reading bin for MacOS
	READPORT := ./bin/readport-darwin-amd64
else
	# Print an unrecognized operating system error
	$(error Unknown operating system: $(UNAME_S))
endif

# Command to build binaries used in the playground
build-readport: build-readport-linux build-readport-darwin build-readport-windows

# Command to build the Linux binary used in the playground
build-readport-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/readport-linux-amd64 ./cmd/readport.go

# Command to build the Darwin binary (MacOS) used in the playground
build-readport-darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/readport-darwin-amd64 ./cmd/readport.go

# Command to build the Windows binary used in the playground
build-readport-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/readport-windows-amd64.exe ./cmd/readport.go

define check_env_param
	echo "Getting env name argument..."

	# Check the ENV passed in the argument
    if [ -z "${ENV}" ]; then echo "No ENV argument provided. Usage: make run ENV=value"; exit 1; fi;
endef

define check_json_exists
	echo "Checking if json ./gopen/$(ENV)/.json exists..."

	# Check if the configuration json exists using the environment name passed
	if [ ! -f "./gopen/$(ENV)/.json" ]; then echo "File ./gopen/$(ENV)/.json does not exist"; exit 1; fi
endef

define get_port_by_json
	echo "Getting port from json ./gopen/$(ENV)/.json..."

	# Defining a temporary variable to store the command output
	$(eval TEMP_PORT := $(shell ./$(READPORT) ./gopen/$(ENV)/.json 2>/dev/null))

	# Check if the command was successful and the TEMP_PORT variable is not empty
	$(eval PORT := $(if $(TEMP_PORT),$(TEMP_PORT),$(error Error to read port on ./gopen/$(ENV)/.json)))
endef

define docker_compose
	# Initialize docker-compose with the environment name and the configured port
	ENV=$(ENV) PORT=$(PORT) TRACER_HOST=$(TRACER_HOST) docker-compose up
endef

# Command to generate a docker image and send to docker-hub
docker-image:
	echo "-----------------------> \033[1mDOCKERFILE\033[0m <-----------------------"
	echo "Building image gabrielhcataldo/gopen-gateway..."
	docker build -t gabrielhcataldo/gopen-gateway:latest .
	echo "Sending to docker-hub..."
	docker push gabrielhcataldo/gopen-gateway:latest

# Command to run API Gateway via docker-compose
run:
	echo "-----------------------> \033[1mMAKEFILE\033[0m <-----------------------"
	$(call check_env_param)
	$(call check_json_exists)
	$(call get_port_by_json)
	echo "Starting docker-compose with ENV=$(ENV) and PORT=$(PORT)..."
	echo ""
	echo "-----------------------> \033[1mDOCKER COMPOSE\033[0m <-----------------------"
	$(call docker_compose)