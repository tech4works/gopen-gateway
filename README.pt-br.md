<img src="logo.png" alt="">

[![Project status](https://img.shields.io/badge/version-v1.0.0_beta-yellow.svg)](https://github.com/GabrielHCataldo/gopen-gateway/releases/tag/v1.0.0-beta)
[![Open Source Helpers](https://www.codetriage.com/gabrielhcataldo/gopen-gateway/badges/users.svg)](https://www.codetriage.com/gabrielhcataldo/gopen-gateway)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/GabrielHCataldo/gopen-gateway)](https://www.tickgit.com/browse?repo=github.com/GabrielHCataldo/gopen-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/GabrielHCataldo/gopen-gateway)](https://goreportcard.com/report/github.com/GabrielHCataldo/gopen-gateway)
[![GoDoc](https://godoc.org/github/GabrielHCataldo/gopen-gateway?status.svg)](https://pkg.go.dev/github.com/GabrielHCataldo/gopen-gateway/helper)

[//]: # ([![build workflow]&#40;https://github.com/GabrielHCataldo/gopen-gateway/actions/workflows/go.yml/badge.svg&#41;]&#40;https://github.com/GabrielHCataldo/gopen-gateway/actions&#41;)

---

![United States](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/us.png "United States")
[Inglês](https://github.com/GabrielHCataldo/gopen-gateway/blob/main/README.md) |
![Spain](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/es.png "Spain")
[Espanhol](https://github.com/GabrielHCataldo/gopen-gateway/blob/main/README.es.md)

O projeto GOPEN foi criado no intuito de ajudar os desenvolvedores a terem uma API Gateway robusta e de fácil manuseio,
com a oportunidade de atuar em melhorias agregando a comunidade, e o mais importante, sem gastar nada. Foi
desenvolvida, pois muitas APIs Gateway do mercado de forma gratuita, não atendem muitas necessidades minímas
para uma aplicação, induzindo-o a fazer o upgrade.

Com essa nova API Gateway você não precisará equilibrar pratos para economizar em sua infraestrutura e arquitetura,
veja abaixo todos os recursos disponíveis:

- Json de configuração simplificado para múltiplos ambientes.
- Configuração rápida de variáveis de ambiente para múltiplos ambientes.
- Versionamento via json de configuração.
- Execução via docker com hot reload opcional.
- Configuração de timeout global e local para cada endpoint.
- Configuração de cache global e local para cada endpoint, com customização da estratégia da chave de armazenamento.
- Armazenamento de cache local ou global utilizando Redis
- Configuração de limitador de tamanho global e local para cada endpoint, limitando o tamanho Header, Body e Multipart
  Memory.
- Configuração de limitador de taxa global e local para cada endpoint, limitando pelo tempo e rajada pelo IP.
- Configuração de segurança de CORS com validações de origens, método http e headers.
- Configuração de global múltiplos middlewares, para serem usados posteriormente no endpoint caso indicado.
- Filtragem personalizada de envio de headers e query para os backends do endpoint.
- Processamento de múltiplos backends, sendo eles beforewares, principais e afterwares para o endpoint.
- Configuração personalizada para abortar processo de execução dos backends pelo código de status retornado.
- Modificadores para todos os conteúdos da requisição e response (Status Code, Path, Header, Params, Query, Body)
  ao nível global (requisição/response do endpoint) e local (atual requisição/response backend) com ações de remoção,
  adição, alteração, substituição e renomeio.
- Obtenha o valor a ser modificado de variáveis de ambiente, da requisição atual, do histórico de respostas do endpoint,
  ou até mesmo do valor passado na configuração.
- Executa os modificadores no contexto que desejar, antes de uma requisição backend ou depois, você decide.
- Faça com que as modificações reflitam em todas as requisições/respostas seguintes, usando a mesma ao nível global.
- Omita a resposta de um backend caso necessite, a mesma não será impressa na resposta do endpoint.
- Omita o body de requisição do seu backend caso precise.
- Agregue suas múltiplas respostas dos backends caso deseje, podendo personalizar o nome do campo a ser alocado a
  resposta do backend.
- Agrupe o body de resposta do seu backend em um campo específico de resposta do endpoint.
- Personalização do tipo de resposta do endpoint podendo ser JSON, TEXT e XML.
- Tenha mais observabilidade com o cadastro automático do trace id no header das requisições seguintes e logs bem
  estruturados.

Documentação
-----------
---
Para entender como funciona, precisamos explicar primeiro a estrutura dos ambientes dinâmicos que GOPEN aceita para sua
configuração em json e arquivo de variáveis de ambiente, então vamos lá!

### Estrutura de pastas

Na estrutura do projeto, em sua raiz precisará ter uma pasta chamada "gopen" e dentro dela precisa ter as pastas
contendo
os nomes dos seus ambientes, você pode dar o nome que quiser, essa pasta precisará ter pelo menos o arquivo ".json"
de configuração da API Gateway, ficará mais o menos assim, por exemplo:

#### Projeto GO:

    gopen-gateway
    | - cmd
    | - internal
    | - gopen
      | - dev
      |   - .json
      |   - .env
      | - prd
      |   - .json
      |   - .env

#### Projeto usando imagem docker:

    nome-do-seu-projeto
    | - docker-compose.yml
    | - gopen
      | - dev
      |   - .json
      |   - .env
      | - prd
      |   - .json
      |   - .env

### Json de configuração

Com base nesse arquivo json de configuração obtido pela env desejada a aplicação terá seus endpoints e suas regras
definidas, veja abaixo um exemplo simples com todos os campos possíveis e seus conceitos e regras:

````json
{
  "$schema": "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json",
  "version": "v1.0.0",
  "port": 8080,
  "hot-reload": true,
  "timeout": "30s",
  "store": {
    "redis": {
      "address": "$REDIS_URL",
      "password": "$REDIS_PASSWORD"
    }
  },
  "cache": {
    "duration": "1m",
    "strategy-headers": [
      "X-Forwarded-For",
      "Device"
    ],
    "allow-cache-control": true
  },
  "limiter": {
    "max-header-size": "1MB",
    "max-body-size": "3MB",
    "max-multipart-memory-size": "10MB",
    "rate": {
      "capacity": 5,
      "every": "1s"
    }
  },
  "security-cors": {
    "allow-origins": [
      "*"
    ],
    "allow-methods": [
      "*"
    ],
    "allow-headers": [
      "*"
    ]
  },
  "middlewares": {
    "save-device": {
      "hosts": [
        "http://192.168.1.2:8051"
      ],
      "path": "/devices",
      "method": "PUT",
      "forward-headers": [
        "*"
      ],
      "modifiers": {
        "header": [
          {
            "context": "RESPONSE",
            "scope": "REQUEST",
            "global": true,
            "action": "SET",
            "key": "X-Device-Id",
            "value": "#response.body.id"
          }
        ]
      }
    }
  },
  "endpoints": [
    {
      "@comment": "Feature: Find user by key",
      "path": "/users/find/:key",
      "cache": {
        "duration": "30s"
      },
      "method": "GET",
      "response-encode": "JSON",
      "beforeware": [
        "save-device"
      ],
      "backends": [
        {
          "hosts": [
            "$USER_SERVICE_URL"
          ],
          "path": "/users/find/:key",
          "method": "GET",
          "forward-headers": [
            "X-Device-Id",
            "X-Forwarded-For",
            "X-Trace-Id"
          ],
          "forward-queries": [
            "*"
          ],
          "modifiers": {
            "statusCode": {},
            "header": [],
            "params": [],
            "query": [],
            "body": []
          },
          "extra-config": {
            "group-response": "",
            "omit-request-body": false,
            "omit-response": false
          }
        }
      ]
    }
  ]
}
````

- #### $schema

Campo obrigatório, para o auxílio na escrita e regras do próprio json de configuração.

- #### version

Campo opcional, usado para retorno do endpoint estático ``/version``.

- ### port

Campo obrigatório, utilizado para indicar a porta a ser ouvida pela API Gateway, valor mínimo 1 e valor máximo 65535.

- ### hot-reload

Campo opcional, o valor padrão é ``false``, caso seja ``true`` é utilizado para o carregamento automático quando
houver alguma alteração no arquivo .json e .env na pasta do ambiente selecionado.

- ### timeout

Campo opcional, o valor padrão é 30 segundos, esse campo é responsável pelo tempo máximo de duração do processamento
de cada requisição, caso seja informado no objeto de endpoint, damos prioridade ao valor informado, caso contrário
seguiremos com o valor informado nesse campo.

````
- Valores aceitos:
    - s para segundos
    - m para minutos
    - h para horas
    - ms para milissegundos
    - us (ou µs) para microssegundos
    - ns para nanossegundos

- Exemplos:
    - 10s
    - 5ms
    - 1h30m
    - 1.5h
````

- ### store

Campo opcional, valor padrão é o armazenamento local em cache, caso seja informado, o campo ``redis`` passa
a ser obrigatório e os outros dois campos que acompanham o mesmo ``address`` e ``password`` também.

Caso utilize o armazenamento global de cache o Redis, é indicado que os valores de endereço e senha sejam preenchidos
utilizando variável de ambiente, como no exemplo acima.

- ### cache


Usabilidade
-----------
---
Use o projeto [playground](https://github.com/GabrielHCataldo/gopen-gateway-playground) para começar a explorar e
utilizar na prática o GOPEN API Gateway!


Como contríbuir?
------------
---


Agradecimentos
------------
---

