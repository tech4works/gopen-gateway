# Build phase
FROM golang:1.23.1 AS builder

# Set the working directory
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy necessary files for the app build
COPY ./json-schema.json ./json-schema.json
COPY ./cmd/main.go ./cmd/main.go
COPY ./internal ./internal

# Change to the directory containing "main.go"
WORKDIR /app/cmd

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main . && rm main.go

# Execution phase (using Alpine)
FROM alpine:latest

# Install CA certificates (for https requests)
RUN apk --no-cache add ca-certificates

# Set the working directory to /app
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/cmd/main ./main

# Copy the .json-schema.json file
COPY --from=builder /app/json-schema.json ./json-schema.json

# Create the runtime folder
RUN mkdir -p ./runtime

# Check the binary executable exists
RUN chmod +x ./main

# Command to run the Go application
CMD ./main
