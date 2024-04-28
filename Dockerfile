# Início do arquivo Dockerfile
FROM golang:1.22 AS builder

# Definindo o diretório de trabalho
WORKDIR /app

# Copiando os arquivos go mod e sum
COPY go.mod go.sum ./

# Baixando todas as dependências
RUN go mod download

# Copiando o arquivo .env
COPY ./.env ./.env

# Copiando a pasta cmd
COPY ./cmd/main.go ./cmd/main.go

# Copiando a pasta internal
COPY ./internal ./internal

# Copiando a pasta gopen
COPY ./gopen ./gopen

# Mude para o diretório contendo o arquivo "main.go"
WORKDIR /app/cmd

# Construindo o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . && rm main.go

# Fase de execução
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiando do builder
COPY --from=builder /app/cmd .

# Copiando o arquivo .env
COPY --from=builder /app/.env ./.env

# Copiando a pasta gopen
COPY --from=builder /app/gopen ./gopen

# Criamos a pasta runtime no repositorio de trabalho root
RUN mkdir -p ./runtime

# Adicionando uma variável ARG para receber um valor externo
ARG ENV_NAME
ENV ENV_NAME ${ENV_NAME}

# Comando para executar
CMD ./main ${ENV_NAME}