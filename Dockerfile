# Build phase
FROM golang:1.22 AS builder

# Setting the working directory
WORKDIR /app

# Copying the go mod and sum files
COPY go.mod go.sum ./

# Downloading all dependencies
RUN go mod download

# Copy the .json-schema.json file
COPY ./json-schema.json ./json-schema.json

# Copy the cmd folder
COPY ./cmd/main.go ./cmd/main.go

# Copy the internal folder
COPY ./internal ./internal

# Change to the directory containing the "main.go" file
WORKDIR /app/cmd

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . && rm main.go

# Execution phase
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Change to the "root" working directory
WORKDIR /root/

# Copy the main files from the cmd folder
COPY --from=builder /app/cmd .

# Copy the .json-schema.json file
COPY --from=builder /app/json-schema.json ./json-schema.json

# Create the runtime folder in the "root" working repository
RUN mkdir -p ./runtime

# Command to run the Go application
CMD ./main