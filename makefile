run:
	echo "Getting env name argument..."

	$(eval ENV_NAME := $(filter-out $@,$(MAKECMDGOALS)))
	@if [ -z "$(ENV_NAME)" ]; then \
		echo "Error: No argument provided. Usage: make run {env name}"; \
		exit 1; \
	fi

	echo "Checking if json ./gopen/$(ENV_NAME)/.json exists..."

	@if [ ! -f "./gopen/$(ENV_NAME)/.json" ]; then \
		echo "Error: File ./gopen/$(ENV_NAME)/.json does not exist"; \
		exit 1; \
	fi

	echo "Getting port from json ./gopen/$(ENV_NAME)/.json..."

	$(eval PORT := $(shell go run ./cmd/readport.go ./gopen/$(ENV_NAME)/.json))

	echo "Starting docker with ENV_NAME=$(ENV_NAME) and PORT=$(PORT)..."

	ENV_NAME=$(ENV_NAME) PORT=$(PORT)  docker-compose up
%:
	@: