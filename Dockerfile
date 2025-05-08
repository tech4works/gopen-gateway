# Build phase (Go binary creation)
FROM golang:1.23.1 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy necessary files for the app build
COPY ./json-schema.json ./json-schema.json
COPY ./cmd/main.go ./cmd/main.go
COPY ./internal ./internal

# Change to the directory containing main.go
WORKDIR /app/cmd

# Build the Go app for Linux/amd64 architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Execution phase (using Alpine to run the binary)
FROM alpine:latest

# Install SSL certificates (required for https requests)
RUN apk --no-cache add ca-certificates

# Set the working directory for execution
WORKDIR /app

# Copy the compiled binary from the build phase
COPY --from=builder /app/cmd/main ./main

# Copy the .json-schema.json file
COPY --from=builder /app/json-schema.json ./json-schema.json

# Create the runtime folder
RUN mkdir -p ./runtime

# Check the binary executable exists
RUN chmod +x ./main

# Command to run the Go application
CMD ["/app/main"]
