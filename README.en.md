<img src="assets/logo.png" alt="">

[![Project status](https://img.shields.io/badge/version-v1.0.0-orange.svg)](https://github.com/GabrielHCataldo/gopen-gateway/releases/tag/v1.0.0)
[![Texto de exemplo para imagem](https://img.shields.io/badge/%F0%9F%8F%90-playground-9900cc.svg)](https://github.com/GabrielHCataldo/gopen-gateway)
[![Open Source Helpers](https://www.codetriage.com/gabrielhcataldo/gopen-gateway/badges/users.svg)](https://www.codetriage.com/gabrielhcataldo/gopen-gateway)
[![TODOs](https://img.shields.io/badge/TODOs-13-red.svg?style=flat)](https://github.com/GabrielHCataldo/gopen-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/GabrielHCataldo/gopen-gateway)](https://goreportcard.com/report/github.com/GabrielHCataldo/gopen-gateway)
[![GoDoc](https://godoc.org/github/GabrielHCataldo/gopen-gateway?status.svg)](https://pkg.go.dev/github.com/GabrielHCataldo/gopen-gateway/helper)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/GabrielHCataldo/gopen-gateway)](https://github.com/GabrielHCataldo/gopen-gateway/blob/main/go.mod)
[![HitCount](https://hits.dwyl.com/GabrielHCataldo/gopen-gateway.svg)](http://hits.dwyl.com/GabrielHCataldo/gopen-gateway)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway.svg?type=small)](https://app.fossa.com/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway?ref=badge_small)

![Brazil](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/br.png)
[Portuguese](https://github.com/GabrielHCataldo/gopen-gateway/blob/main/README.md) |
![Spain](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/es.png)
[Spanish](https://github.com/GabrielHCataldo/gopen-gateway/blob/main/README.es.md)

The GOPEN project was created with the aim of helping developers to have a robust and easy-to-handle API Gateway, with
the opportunity to act on improvements, contributing to the community, and most importantly, without incurring any
costs. It was developed because many freely available API Gateways in the market do not meet the many minimum
requirements for an application, leading you to upgrade.

With this new API Gateway, you will not need to balance plates to save on your infrastructure and architecture, and you
will still be able to optimise your development. Below are the features available:

- Simplified JSON configuration for multiple environments.
- Quick configuration of environment variables for multiple environments.
- Versioning via JSON configuration.
- Docker execution with optional hot reload.
- Configuration of global and local timeout for each endpoint.
- Configuration of global and local cache for each endpoint, with customisation of the storage key strategy, and
  conditions based on response status codes and HTTP method for reading and saving the same.
- Local or global cache storage using Redis.
- Configuration of size limit, global and local for each endpoint, limiting the size of the Header, Body and Multipart
  Memory.
- Configuration of rate limit, global and local for each endpoint, rate-limited by time and burst by IP.
- Configuration of CORS security with validations of origins, HTTP method and headers.
- Global configuration of multiple middlewares, to be used later at the endpoint if indicated.
- Customised filtering of sending headers and query to the endpoint backends.
- Processing of multiple backends, being them beforewares, mains and afterwares for the endpoint.
- Custom configuration to abort the process of backend execution by the status code returned.
- Modifiers for all request and response contents (Status Code, Path, Header, Params, Query, Body) at global (
  request/response) and local (backend request/backend response) levels with removal, addition, modification,
  replacement, and renaming actions.
- Get the value to be modified from environment variables, the current request, the history of endpoint responses, or
  even from the value passed in the configuration.
- Run the modifiers in the context you prefer, before a backend request or afterwards, it's your decision.
- Make the modifications reflect in all subsequent requests/responses, by using the same at a global level.
- Omit a backend response if needed, it will not be used in the endpoint response.
- Omit the request body from your backend if needed.
- Aggregate multiple responses from your backends if desired, customizing the name of the field to be allocated to the
  backend response.
- Group the body response from your backend in a specific field of the endpoint response.
- Customise the type of endpoint response, which could be JSON, TEXT and XML.
- Increase your observability with automatic registration of the trace id in the header of the following requests and
  structured logs.

# Documentation

To understand how it works, we first need to explain the structure of the dynamic environments that GOPEN accepts for
its JSON configuration and environment variable files, so let's get started!

## Folder Structure

In the project structure, you will need to have a root folder called "gopen". Inside this folder, you should have
subfolders containing the names of your environments. You can name these subfolders as you wish. Each of these
subfolders should contain at least one ".json" configuration file for the API Gateway. For example, the structure might
look like this:

### GO Project

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

### Docker Image Project

    nome-do-seu-projeto
    | - docker-compose.yml
    | - gopen
      | - dev
      |   - .json
      |   - .env
      | - prd
      |   - .json
      |   - .env

## JSON Configuration

Based on this JSON configuration file obtained by the desired environment, the application will have its endpoints and
rules defined. Below is a simple example with all possible fields and their concepts and rules:

````json
{
  "$schema": "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json",
  "@comment": "My First API Gateway",
  "version": "v1.0.0",
  "port": 8080,
  "hot-reload": true,
  "store": {
    "redis": {
      "address": "$REDIS_URL",
      "password": "$REDIS_PASSWORD"
    }
  },
  "timeout": "30s",
  "cache": {
    "duration": "1m",
    "strategy-headers": [
      "X-Forwarded-For",
      "Device"
    ],
    "only-if-status-codes": [
      200,
      201,
      202,
      203,
      204
    ],
    "only-if-methods": [
      "GET"
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
    "allow-origins": [],
    "allow-methods": [],
    "allow-headers": []
  },
  "middlewares": {
    "save-device": {
      "@comment": "Middleware service to validate and save Device",
      "name": "device-manager",
      "hosts": [
        "http://192.168.1.2:8051"
      ],
      "path": "/devices",
      "method": "PUT",
      "forward-headers": [],
      "modifiers": {
        "@comment": "Primeiro modificador",
        "header": [
          {
            "@comment": "Adds the X-Device-Id field with the value from the middleware's response body id",
            "context": "RESPONSE",
            "scope": "REQUEST",
            "propagate": true,
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
      "@comment": "Endpoint responsible for returning the user found by the key, which can be either email or phone.",
      "path": "/users/find/:key",
      "method": "GET",
      "timeout": "10s",
      "cache": {
        "enabled": true,
        "duration": "30s",
        "strategy-headers": [],
        "only-if-status-codes": [],
        "allow-cache-control": false
      },
      "limiter": {
        "max-header-size": "1MB",
        "max-body-size": "1MB",
        "max-multipart-memory-size": "1MB",
        "rate": {
          "capacity": 10,
          "every": "1s"
        }
      },
      "response-encode": "JSON",
      "aggregate-responses": false,
      "abort-if-status-codes": [],
      "beforewares": [
        "save-device"
      ],
      "afterwares": [],
      "backends": [
        {
          "@comment": "Backend service responsible for retrieving the user by the key.",
          "name": "user",
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
          "forward-queries": [],
          "modifiers": {
            "@comment": "My first modifiers.",
            "status-code": 0,
            "header": [],
            "param": [],
            "query": [],
            "body": []
          },
          "extra-config": {
            "group-response": false,
            "omit-request-body": false,
            "omit-response": false
          }
        }
      ]
    }
  ]
}
````

TODO:

### $schema

Mandatory field, to assist in writing and rules of the configuration JSON itself, the only accepted value is:

```text
https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json`
```

### @comment

Optional field, of type string, free for notes related to your project.

### version

Optional field, used for version control and also used in the return of the static endpoint [/version](#version-1).

### port

Mandatory field, used to indicate the port to be listened to by the API Gateway, minimum value `1` and maximum
value `65535`.

### hot-reload

Optional field, the default value is `false`, if `true` it is used for automatic loading when there is a change in the
.json and .env file in the selected environment folder.

### store

Optional field, default value is the local cache storage, if informed, the `redis` field becomes mandatory and
the `address` field as well.

If you use Redis for global cache storage, it is recommended that address and password values be filled in using an
environment variable, as in the example above.

### timeout

Optional field, the default value is `30 seconds`, this field is responsible for the maximum processing time for each
request.

If the request exceeds this time, the API Gateway will abort all ongoing transactions and return the status
code `504 (Gateway Timeout)`. Learn more about this return by clicking [here](#504-gateway-timeout).

IMPORTANT: If it is informed in the endpoint object, we give priority to the informed value of the endpoint, otherwise
we will continue with the informed or default value of this field, at the root of the JSON configuration.

```text
- Accepted values:
    - s for seconds
    - m for minutes
    - h for hours
    - ms for milliseconds 
    - us (or µs) for microseconds
    - ns for nanoseconds

- Examples:
    - 10s
    - 5ms
    - 1h30m
    - 1.5h
```

### cache

Optional field, if informed, the `duration` field becomes mandatory!

If the object is informed in the [endpoint.cache](#cache) structure, we give priority to the values informed there,
otherwise we will follow the values informed in this field.

The cache value is only recorded once every X duration informed in the `every` field.

The `only-if-status-codes` and `only-if-methods` fields are used to check if that enabled endpoint can read and write
the cache based on the HTTP method and response status code, learn more about them below.

If the answer is not "fresh", in other words, was answered by the cache, the header `X-Gopen-Cache` will have the
value `true` otherwise the value will be `false`.

#### cache.duration

Indicates the time that the cache will last, it is of type `time.Duration`.

```text
- Accepted values:
    - s for seconds
    - m for minutes
    - h for hours
    - ms for milliseconds 
    - us (or µs) for microseconds
    - ns for nanoseconds

- Examples:
    - 10s
    - 5ms
    - 1h30m
    - 1.5h
```

#### cache.strategy-headers

Optional field, the default cache key strategy is by the request's url and HTTP method making it a global cache per
endpoint, if the headers to be used in the strategy are informed, they are aggregated in the default key values, for
example, there in the example it was indicated to use the `X-Forwarded-For` and `Device` fields, the final value of the
key would be:

      GET:/users/find/479976139:177.130.228.66:95D4AF55-733D-46D7-86B9-7EF7D6634EBC

The description of the logic behind this key is:

      method:url:X-Forwarded-For:Device

Without the filled strategy, the default logic is:

      method:url

So the default value for this endpoint without the filled strategy is:

      GET:/users/find/479976139

In this example we make the cache that was global for the specific endpoint, to be per client! Remembering that this is
a simple example, you can have the strategy you want based on the header of your application.

#### cache.only-if-methods

Optional field, the default value is a list with only the HTTP `GET` method, if informed empty, any HTTP method will be
accepted.

This field is responsible for deciding whether to read and write the cache of the endpoint (which is enabled to have
cache) by the same HTTP method.

#### cache.only-if-status-codes

Optional field, the default value is a list of recognized success HTTP status codes, if it is informed empty, any HTTP
response status code will be accepted.

This field is responsible for deciding whether to write the cache of the endpoint (which is enabled to have cache) by
the HTTP response status code of the same.

#### cache.allow-cache-control

Optional field, the default value is `false`, if it is informed as `true` the API Gateway will consider
the `Cache-Control` header following the rules below from the value informed in the request or in the response from the
backends:

`no-cache`: this value is only considered in the request header, if informed we disregard the cache read and continue
with the normal process to obtain the "fresh" response.

`no-store`: this value is only considered in the response written by your backends, if informed we do not store the
cache.

### limiter

Optional field, object responsible for the API Gateway limitation rules, whether by size or rate, the default values
vary from field to field, see:

#### limiter.max-header-size

Optional field, it is of type `byteUnit`, default value is `1MB`, is responsible for limiting the size of the request
header.

If the size of the header exceeds the informed value, the API Gateway will abort the request with status
code `431 (Request header fields too large)`. Learn more about this return by
clicking [here](#431-request-header-fields-too-large).

```text
- Accepted values:
  - B for Byte
  - KB for KiloByte
  - MB for Megabyte
  - GB for Gigabyte
  - TB for Terabyte
  - PB for Petabyte
  - EB for Exabyte
  - ZB for Zettabyte
  - YB for Yottabyte

- Examples:
  - 1B
  - 50KB
  - 5MB
  - 1.5GB
```

TODO:

#### limiter.max-body-size

Campo opcional, ele é do tipo `byteUnit`, valor padrão é `3MB`, campo é responsável por limitar o tamanho do corpo
da requisição.

Caso o tamanho do corpo ultrapasse o valor informado, a API Gateway irá abortar a requisição com o código de status
`413 (Request entity too large)`.

```text
- Accepted values:
    - B for Byte
    - KB for KiloByte
    - MB for Megabyte
    - GB for Gigabyte
    - TB for Terabyte
    - PB for Petabyte
    - EB for Exabyte
    - ZB for Zettabyte
    - YB for Yottabyte

- Examples:
    - 1B
    - 50KB
    - 5MB
    - 1.5GB
```

#### limiter.max-multipart-memory-size

Campo opcional, ele é do tipo `byteUnit`, valor padrão é `5MB`, esse campo é responsável por limitar o tamanho do
corpo multipart/form da requisição, geralmente utilizado para envio de arquivos, imagens, etc.

Caso o tamanho do corpo ultrapasse o valor informado, a API Gateway irá abortar a requisição com o código de status
`413 (Request entity too large)`. Veja mais sobre esse retorno clicando [aqui](#413-request-entity-too-large).

```text
- Accepted values:
    - B for Byte
    - KB for KiloByte
    - MB for Megabyte
    - GB for Gigabyte
    - TB for Terabyte
    - PB for Petabyte
    - EB for Exabyte
    - ZB for Zettabyte
    - YB for Yottabyte

- Examples:
    - 1B
    - 50KB
    - 5MB
    - 1.5GB
```

#### limiter.rate

Campo opcional, caso seja informado, o campo `capacity` torna-se obrigatório, esse objeto é responsável por limitar
a taxa de requisição pelo IP, esse limite é imposto obtendo a capacidade máxima pelo campo `capacity` por X duração,
informado no campo `every`.

Caso essa capacidade seja ultrapassada, a API Gateway por segurança abortará a requisição, retornando
`429 (Too many requests)`. Veja mais sobre esse retorno clicando [aqui](#429-too-many-requests).

#### limiter.rate.capacity

Campo opcional, caso o objeto rate seja informado, ele passa a ser obrigatório, o valor padrão é `5`, e o mínimo
que poderá ser informado é `1`, indica a capacidade máxima de requisições.

#### limiter.rate.every

Campo opcional, o valor padrão é `1 segundo`, indica o valor da duração da verificação da capacidade máxima de
requisições.

### security-cors

Campo opcional, usado para segurança do CORS da API Gateway, todos os campos por padrão são vazios, não restringindo
os valores de origin, methods e headers.

Caso queira restringir, e a requisição não corresponda com as configurações impostas, a API Gateway por segurança
irá abortar a requisição retornando `403 (Forbidden)`.

#### security-cors.allow-origins

Campo opcional, do tipo lista de string, os itens da lista precisam indicar quais IPs de origem a API Gateway
permite receber nas requisições.

#### security-cors.allow-methods

Campo opcional, do tipo lista de string, os itens da lista precisam indicar quais métodos HTTP a API Gateway
permite receber nas requisições.

#### security-cors.allow-headers

Campo opcional, do tipo lista de string, os itens da lista precisam indicar quais campos de cabeçalho HTTP a API Gateway
permite receber nas requisições.

### middlewares

Campo opcional, é responsável pela configuração de seus middlewares de aplicação, é um mapa com chaves
em string mencionando o nome do seu middleware, esse nome poderá ser utilizado em seus [endpoints](#endpoints)
como `beforewares` e `afterwares`.

O valor da chave é um objeto de [backend](#backendname), porém, com uma observação, esse objeto terá
sua resposta caso não [abortada](#resposta-abortada), omitida automáticamente pelo endpoint, apenas seu cabeçalho de
resposta poderá ser retornado, caso tenha, já que respostas de middlewares não são exibidas para o cliente final HTTP,
porém, sua resposta será armazenada ao longo da requisição HTTP feita no endpoint, podendo ter seus valores de
requisição e resposta obtidos e manipulados.

Por exemplo, um `beforeware` quando mencionado no endpoint, ele será utilizado como middleware de pré-requisições, isto
é, ele será chamado antes dos backends principais do endpoint, então podemos, por exemplo, ter um middleware
de manipulação de device, como no JSON de configuração acima, aonde ele irá chamar esse backend de middleware
configurado no endpoint como `beforewares`, validando e salvando o dispositivo a partir de informações do header da
requisição, caso o backend responda um código de status de falha, no exemplo, o gateway abortará todos os backends
seguintes retornando o que o backend de device respondeu, caso tenha retornado um código de status de sucesso, ele irá
modificar o header de todas as requisições seguintes (`propagate:true`), adicionando o campo `X-Device-Id`, com o valor
do id do body de resposta do próprio backend, podendo ser utilizado nos outros backends seguintes do endpoint.

Para saber mais sobre os `modifiers` [veja](#backendmodifiers).

Para entender melhor essa ferramenta poderosa, na prática, veja os exemplos de middlewares usados como `beforewares`
e `afterwares` feitos no projeto de [playground](https://github.com/GabrielHCataldo/gopen-gateway-playground).

### endpoints

Campo obrigatório, é uma lista de objeto, representa cada endpoint da API Gateway que será registrado para ouvir e
servir as requisições HTTP.

Veja como funciona o fluxo básico de um endpoint na imagem abaixo:

<img src="assets/basic-endpoint-flowchart-pt-br.png" alt="fluxo básico de endpoint">

Abaixo iremos listar e explicar cada campo desse objeto tão importante:

### endpoint.@comment

Campo opcional, do tipo string, campo livre para anotações relacionadas ao seu endpoint.

### endpoint.path

Campo obrigatório, do tipo string, responsável pelo caminho URI do endpoint, exemplo `"/users/:id"`.

Caso queira ter parâmetros dinâmicos nesse endpoint, apenas use o padrão `":nome do parâmetro"` por exemplo
`"/users/:id/status/:status"`, a API Gateway irá entender que teremos 2 parâmetros dinâmicos desse endpoint,
esses valores podem ser repassados para os backends subjacentes.

Exemplo usando o parâmetro dinâmico para os backends subjacentes:

- Endpoint
    - path: `"/users/:id/status/:status"`
    - resultado: `"/users/1/status/removed"`
- Backend 1
    - path: `"/users/:id"`
    - resultado: `"/users/1"`
- Backend 2
    - path: `"/users/:id/status/:status"`
    - resultado: `"/users/1/status/removed"`

No exemplo acima vemos que o parâmetro pode ser utilizado como quiser como path nas requisições de backend do endpoint
em questão.

### endpoint.method

Campo obrigatório, do tipo string, responsável por definir qual método HTTP o endpoint será registrado.

#### endpoint.timeout

É semelhante ao campo [timout](#timeout), porém, será aplicado apenas para o endpoint
em questão.

Caso omitido, será herdado o valor do campo [timeout](#timeout).

### endpoint.cache

Campo opcional, do tipo objeto, por padrão ele virá vazio apenas com o campo `enabled` preenchido com o valor `false`.

Caso informado, o campo `enabled` se torna obrigatório, os outros campos, caso omitidos, irá herdar da configuração
[cache](#cache) na raiz caso exista e se preenchida.

Se por acaso, tenha omitido o campo `duration` tanto na atual configuração como na configuração [cache](#cache) na raiz,
o campo `enabled` é ignorado considerando-o sempre como `false` pois não foi informado a duração do cache em ambas
configurações.

#### endpoint.cache.enabled

Campo obrigatório, do tipo booleano, indica se você deseja que tenha cache em seu endpoint, `true` para habilitado,
`false` para não habilitado.

Caso esteja `true` mas não informado o campo `duration` na configuração atual e nem na [raiz](#cache), esse campo
será ignorado considerando-o sempre como `false`.

#### endpoint.cache.ignore-query

Campo opcional, do tipo booleano, caso não informado o valor padrão é `false`.

Caso o valor seja `true` a API Gateway irá ignorar os parâmetros de busca da URL ao
criar a chave de armazenamento, caso contrário ela considerára os parâmetros de busca da URL
ordenando alfabéticamente as chaves e valores.

#### endpoint.cache.duration

É semelhante ao campo [cache.duration](#cacheduration), porém, será aplicado apenas para o endpoint
em questão.

Caso omitido, será herdado o valor do campo [cache.duration](#cacheduration).

Caso seja omitido nas duas configurações, o campo `enabled` será ignorado considerando-o sempre como `false`.

#### endpoint.cache.strategy-headers

É semelhante ao campo [cache.strategy-headers](#cachestrategy-headers), porém, será aplicado apenas para o endpoint
em questão.

Caso omitido, será herdado o valor do campo [cache.strategy-headers](#cachestrategy-headers).

Caso seja informado vazio, o valor do não será herdado, porém, será aplicado o valor [padrão](#cachestrategy-headers)
para o endpoint em questão.

#### endpoint.cache.only-if-status-codes

É semelhante ao campo [cache.only-if-status-codes](#cacheonly-if-status-codes), porém, será aplicado apenas para o
endpoint em questão.

Caso omitido, será herdado o valor do campo [cache.only-if-status-codes](#cacheonly-if-status-codes).

Caso seja informado vazio, o valor do não será herdado, porém, será aplicado o
valor [padrão](#cacheonly-if-status-codes) para o endpoint em questão.

#### endpoint.cache.allow-cache-control

É semelhante ao campo [cache.allow-cache-control](#cacheallow-cache-control), porém, será aplicado apenas para o
endpoint em questão.

Caso omitido, será herdado o valor do campo [cache.allow-cache-control](#cacheallow-cache-control).

#### endpoint.limiter

É semelhante ao campo [limiter](#limiter), porém, será aplicado apenas para o endpoint
em questão.

Caso omitido, será herdado o valor do campo [limiter](#limiter).

#### endpoint.response-encode

Campo opcional, do tipo string, o valor padrão é vazio, indicando que a resposta do endpoint será codificada seguindo
a [lógica de resposta](#lógica-de-resposta) da API Gateway, sem forçar a codificação indicada.

```text
- Valores aceitos:
  - JSON 
  - XML
  - TEXT
```

#### endpoint.aggregate-responses

Campo opcional, do tipo booleano, o valor padrão é `false`, indicando que a resposta do endpoint não será agregada.

Caso informado com o valor `true` e tiver mais de uma resposta dos backends informados no endpoint ele irá agregar as
respostas dos backends, veja mais sobre as regras de resposta da API Gateway clicando [aqui](#lógica-de-resposta).

#### endpoint.abort-if-status-codes

Campo opcional, do tipo lista de inteiros, o valor padrão é nulo, indicando que qualquer backend executado no endpoint
que tenha respondido o código de status HTTP maior ou igual a `400 (Bad request)` será abortado.

Caso informado, e um backend retorna o código de status HTTP indicado na configuração, o endpoint será abortado, isso
significa que os outros backends configurados após o mesmo, não serão executados, e o endpoint irá retornar a resposta
do mesmo ao cliente final.

Caso queira que nenhum código de status HTTP seja abortado no endpoint, apenas informe o campo vazio.

Veja como o endpoint será respondido após um backend ser abortado clicando [aqui](#lógica-de-resposta).

#### endpoint.beforewares

Campo opcional, do tipo lista de string, o valor padrão é vazio, indicando que o endpoint não tem nenhum middleware
de pré-requisições.

Caso informado, o endpoint irá executar as requisições, posição por posição, começando no início da lista. Caso o valor
em string da posição a ser executada estiver configurada no campo [middlewares](#middlewares) corretamente, será
executado
o backend configurado no mesmo. Caso contrário irá ignorar a posição apenas imprimindo um log de atenção.

#### endpoint.afterwares

Campo opcional, do tipo lista de string, o valor padrão é vazio, indicando que o endpoint não tem nenhum middleware
de pós-requisições.

Caso informado, o endpoint irá executar as requisições, posição por posição, começando no início da lista. Caso o valor
em string da posição a ser executada estiver configurada no campo [middlewares](#middlewares) corretamente, será
executado o backend configurado no mesmo. Caso contrário irá ignorar a posição apenas imprimindo um log de atenção.

#### endpoint.backends

Campo obrigatório, do tipo lista de objeto, responsável pela execução principal do endpoint, o próprio nome já diz tudo,
é uma lista que indica todos os serviços necessários para que o endpoint retorne a resposta esperada.

Veja como funciona o fluxo básico de um backend na imagem abaixo:

<img src="assets/basic-backend-flowchart-pt-br.png" alt="fluxo básico de um backend">

Abaixo iremos listar e explicar cada campo desse objeto tão importante:

### backend.@comment

Campo opcional, do tipo string, campo livre para anotações relacionadas ao seu backend.

### backend.name

Campo opcional, do tipo string, é responsável pelo nome do seu serviço backend, é utilizado para dar nome ao campo de
resposta agregada do mesmo, caso o campo [backend.extra-config.group-response](#backendextra-configgroup-response)
esteja como `true`.

### backend.hosts

Campo obrigatório, do tipo lista de string, é responsável pelos hosts do seu serviço que a API Gateway irá chamar
juntamente com o campo [backend.path](#backendpath).

De certa forma podemos ter um load balancer "burro", pois o backend irá sortear nessa lista qual host irá ser chamado,
com isso podemos informar múltiplas vezes o mesmo host para balancear as chamadas, veja:

````
50% cada
[
  "https://instance-01", 
  "https://instance-02"
]
````

````
instance-01: 15%
instance-02: 75%
[
  "https://instance-01", 
  "https://instance-02",
  "https://instance-02",
  "https://instance-02"
]
````

````
instance-01: 33.3%
instance-02: 66.7%
[
  "https://instance-01", 
  "https://instance-02",
  "https://instance-02"
]
````

### backend.path

Campo obrigatório, do tipo string, o valor indica a URL do caminho do serviço backend.

Utilizamos um dos [backend.hosts](#backendhosts) informados e juntamos com o path fornecido, por exemplo, no campo
hosts temos o valor

```text
[
  "https://instance-01", 
  "https://instance-02"
]
```

E nesse campo path temos o valor

```text
/users/status
```

O backend irá construir a seguinte URL de requisição

```text
https://instance-02/users/status
```

Veja como o host é balanceado clicando [aqui](#backendhosts).

### backend.method

Campo obrigatório, do tipo string, o valor indica qual método HTTP o serviço backend espera.

### backend.forward-queries

Campo opcional, do tipo lista de string, o valor padrão é vazio, indicando que qualquer parâmetro de busca será
repassado para o serviço backend.

Caso informado, apenas os campos indicados serão repassados para o serviço backend, por exemplo, recebemos uma
requisição com a seguinte URL

````
/users?id=23&email=gabrielcataldo@gmail.com&phone=47991271234
````

Nesse exemplo temos o campo `forward-queries` com os seguintes valores

````
[
  "email",
  "phone"
]
````

A URL de requisição ao backend foi

````
/users?email=gabrielcataldo@gmail.com&phone=47991271234
````

Vimos que o parâmetro de busca `id` não foi repassado para o serviço backend, pois ele não foi mencionado na lista.

### backend.forward-headers

Campo opcional, do tipo lista de string, o valor padrão é vazio, indicando que qualquer cabeçalho recebido será
repassado para o serviço backend.

Caso informado, apenas os campos indicados serão repassados para o serviço backend, por exemplo, recebemos uma
requisição com o seguinte cabeçalho:

```json
{
  "Device": "95D4AF55-733D-46D7-86B9-7EF7D6634EBC",
  "User-Agent": "IOS",
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

Nesse exemplo temos o campo `forward-headers` com os seguintes valores:

```json
[
  "User-Agent",
  "Authorization"
]
```

O cabeçalho de requisição ao backend foi:

```json
{
  "User-Agent": "IOS",
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

Vimos que o campo `Device` do cabeçalho recebido não foi repassado para o serviço backend, pois ele não foi mencionado
na lista.

### backend.extra-config

Campo opcional, do tipo objeto, indica configuração extras do serviço backend, veja abaixo sobre os campos e suas
responsabilidades:

#### backend.extra-config.omit-request-body

Campo opcional, do tipo booleano, o valor padrão é `false`, indicando que o corpo da requisição será repassado ao
backend caso tenha.

Caso informado `true` o corpo da requisição não será repassado ao backend.

#### backend.extra-config.group-response

Campo opcional, do tipo booleano, o valor padrão é `false`, indicando que o corpo da resposta do backend não precisará
ser agrupada em um campo JSON para a resposta ao cliente final.

Caso informado com o valor `true` o body de resposta caso tenha, será agrupado em um campo JSON da resposta final,
o nome do campo será o [nome](#backendname) do serviço backend caso preenchido, se não temos um padrão de nomenclatura
que é `backend-posição na lista` que seria por exemplo `backend-0`.

Para entender a importância desse campo, veja mais sobre a [lógica de resposta](#lógica-de-resposta) da API Gateway.

#### backend.extra-config.omit-response

Campo opcional, do tipo booleano, o valor padrão é `false`, indicando que a resposta do backend em questão não será
omitida para o cliente final.

Caso informado com o valor `true` toda a resposta do backend em questão será omitida, tenha cuidado, pois se tiver
apenas
esse backend, e o mesmo for omitido, a API Gateway responderá por padrão o código de status HTTP `204 (No Content)`.

Para entender a importância desse campo, veja mais sobre a [lógica de resposta](#lógica-de-resposta) da API Gateway.

### backend.modifiers

Campo opcional, do tipo objeto, o valor padrão é vazio, indicando não haver nenhum processo de modificação nesse
backend em questão.

Abaixo iremos listar e explicar cada campo desse objeto tão importante:

### modifiers.@comment

Campo opcional, do tipo string, campo livre para anotações relacionadas aos seus modificadores.

### modifiers.status-code

Campo opcional, do tipo inteiro, valor padrão é `0`, indicando não haver nada a ser modificado no código de status HTTP
de resposta do backend.

Caso informado, o código de status HTTP de resposta do backend será modificado pelo valor inserido, isso pode ter ou não
influência na resposta final do endpoint, veja a [lógica-de-resposta](#lógica-de-resposta) da API Gateway para saber
mais.

### modifiers.header

Campo opcional, do tipo lista de objeto, valor padrão é vazio, responsável pelas modificações de cabeçalho da requisição
e resposta do backend.

Veja abaixo os campos desse objeto e suas responsabilidade:

#### header.@comment

Campo opcional, do tipo string, campo livre para anotações relacionadas ao seu modificador.

#### header.context

Campo obrigatório, do tipo string, é responsável por indicar qual contexto a modificação deve atuar.

Valores aceitos:

`REQUEST`: Para atuar na pré-requisição do backend.

`RESPONSE`: Para atuar pós-requisição do backend.

Importante lembrar que caso o valor for `REQUEST` poderá utilizar no campo [header.scope](#headerscope) apenas o valor
`REQUEST`.

#### header.scope

Campo opcional, do tipo string, o valor padrão será baseado no campo [header.context](#headercontext) informado, o valor
indica qual escopo devemos alterar, se o escopo de requisição ou de resposta.

Valores aceitos:

`REQUEST`: Para modificar o escopo de requisição, esse tipo de escopo pode ter uma atuação global propagando essa
mudança
nas requisições backends seguintes, basta utilizar o campo [header.propagate](#headerpropagate) como `true`.

`RESPONSE`: Para modificar o escopo de resposta do backend.

#### header.action

Campo obrigatório, do tipo string, responsável pela ação a ser tomada na modificação do cabeçalho.

Valores aceitos:

`ADD`: Adiciona a chave informada no campo [header.key](#headerkey) caso não exista, e agrega o valor informado no
campo [header.value](#headervalue).

`APD`: Acrescenta o valor informado no campo [header.value](#headervalue) caso a chave informada no
campo  [header.key](#headerkey) exista.

`SET`: Define o valor da chave informada no campo [header.key](#headerkey) pelo valor passado no
campo [header.value](#headervalue).

`RPL`: Substitui o valor da chave informada no campo [header.key](#headerkey) pelo valor passado no
campo [header.value](#headervalue) caso exista.

`REN`: Renomeia a chave informada no campo [header.key](#headerkey) pelo valor passado no
campo [header.value](#headervalue) caso exista.

`DEL`: Remove a chave informada no campo [header.key](#headerkey) caso exista.

#### header.key

Campo obrigatório, do tipo string, utilizado para indicar qual chave do cabeçalho deve ser modificada.

#### header.value

Campo obrigatório, do tipo string, utilizado como valor a ser usado para modificar a chave indicada no
campo [header.key](#headerkey).

Temos possibilidades de utilização de [valores dinâmicos](#valores-dinâmicos-para-modificação),
e de [variáveis de ambiente](#variáveis-de-ambiente) para esse campo.

OBS: se torna opcional apenas se [query.action](#queryaction) tiver o valor `DEL`.

#### header.propagate

Campo opcional, do tipo booleano, o valor padrão é `false` indicando que o modificador não deve propagar essa mudança
em questão para os backends seguintes.

Caso informado como `true` essa modificação será propagada para os seguintes backends.

IMPORTANTE: Esse campo só é aceito se o [escopo](#headerscope) tiver o valor `REQUEST`.

### modifiers.param

Campo opcional, do tipo lista de objeto, valor padrão é vazio, responsável pelas modificações de parâmetros da
requisição para o backend.

Veja abaixo os campos desse objeto e suas responsabilidade:

#### param.@comment

Campo opcional, do tipo string, campo livre para anotações relacionadas ao seu modificador.

#### param.context

Campo obrigatório, do tipo string, é responsável por indicar qual contexto a modificação deve atuar.

Valores aceitos:

`REQUEST`: Para atuar na pré-requisição do backend.

`RESPONSE` Para atuar pós-requisição do backend.

#### param.action

Campo obrigatório, do tipo string, responsável pela ação a ser tomada na modificação dos parâmetros da requisição.

Valores aceitos:

`SET`: Define o valor da chave informada no campo [param.key](#paramkey) pelo valor passado no
campo [param.value](#paramvalue).

`RPL`: Substitui o valor da chave informada no campo [header.key](#headerkey) pelo valor passado no
campo [header.value](#headervalue) caso exista.

`REN`: Renomeia a chave informada no campo [param.key](#paramkey) pelo valor passado no
campo [param.value](#paramvalue) caso exista.

`DEL`: Remove a chave informada no campo [param.key](#paramkey) caso exista.

#### param.key

Campo obrigatório, do tipo string, utilizado para indicar qual chave de parâmetro deve ser modificada.

#### param.value

Campo obrigatório, do tipo string, utilizado como valor a ser usado para modificar a chave indicada no
campo [param.key](#paramkey).

Temos possibilidades de utilização de [valores dinâmicos](#valores-dinâmicos-para-modificação),
e de [variáveis de ambiente](#variáveis-de-ambiente) para esse campo.

OBS: se torna opcional apenas se [query.action](#queryaction) tiver o valor `DEL`.

#### param.propagate

Campo opcional, do tipo booleano, o valor padrão é `false` indicando que o modificador não deve propagar essa mudança
em questão para os backends seguintes.

Caso informado como `true` essa modificação será propagada para os seguintes backends.

### modifiers.query

Campo opcional, do tipo lista de objeto, valor padrão é vazio, responsável pelas modificações de parâmetros de busca da
requisição para o backend.

Veja abaixo os campos desse objeto e suas responsabilidade:

#### query.@comment

Campo opcional, do tipo string, campo livre para anotações relacionadas ao seu modificador.

#### query.context

Campo obrigatório, do tipo string, é responsável por indicar qual contexto a modificação deve atuar.

Valores aceitos:

`REQUEST`: Para atuar na pré-requisição do backend.

`RESPONSE`: Para atuar pós-requisição do backend.

#### query.action

Campo obrigatório, do tipo string, responsável pela ação a ser tomada na modificação dos parâmetros de busca da
requisição.

Valores aceitos:

`ADD`: Adiciona a chave informada no campo [query.key](#querykey) caso não exista, e agrega o valor informado no
campo [query.value](#queryvalue).

`APD`: Acrescenta o valor informado no campo [query.value](#queryvalue) caso a chave informada no
campo [query.key](#querykey) exista.

`SET`: Define o valor da chave informada no campo [query.key](#querykey) pelo valor passado no
campo [query.value](#queryvalue).

`RPL`: Substitui o valor da chave informada no campo [query.key](#querykey) pelo valor passado no
campo [query.value](#queryvalue) caso exista.

`REN`: Renomeia a chave informada no campo [query.key](#querykey) pelo valor passado no
campo [query.value](#queryvalue) caso exista.

`DEL`: Remove a chave informada no campo [query.key](#querykey) caso exista.

#### query.key

Campo obrigatório, do tipo string, utilizado para indicar qual chave de parâmetro de busca deve ser modificada.

#### query.value

Campo obrigatório, do tipo string, utilizado como valor a ser usado para modificar a chave indicada no
campo [query.key](#querykey).

Temos possibilidades de utilização de [valores dinâmicos](#valores-dinâmicos-para-modificação),
e de [variáveis de ambiente](#variáveis-de-ambiente) para esse campo.

OBS: se torna opcional apenas se [query.action](#queryaction) tiver o valor `DEL`.

#### query.propagate

Campo opcional, do tipo booleano, o valor padrão é `false` indicando que o modificador não deve propagar essa mudança
em questão para os backends seguintes.

Caso informado como `true` essa modificação será propagada para os seguintes backends.

### modifiers.body

Campo opcional, do tipo lista de objeto, valor padrão é vazio, responsável pelas modificações de body de
requisição ou resposta do backend.

Veja abaixo os campos desse objeto e suas responsabilidade:

#### body.@comment

Campo opcional, do tipo string, campo livre para anotações relacionadas ao seu modificador.

#### body.context

Campo obrigatório, do tipo string, é responsável por indicar qual contexto a modificação deve atuar.

Valores aceitos:

`REQUEST`: Para atuar na pré-requisição do backend.

`RESPONSE`: Para atuar pós-requisição do backend.

Importante lembrar que caso o valor for `REQUEST` poderá utilizar no campo [body.scope](#bodyscope) apenas o valor
`REQUEST`.

#### body.scope

Campo opcional, do tipo string, o valor padrão será baseado no campo [body.context](#bodycontext) informado, o valor
indica qual escopo devemos alterar, se o escopo de requisição ou de resposta.

Valores aceitos:

`REQUEST`: Para modificar o escopo de requisição, esse tipo de escopo pode ter uma atuação global propagando essa
mudança
nas requisições de backend seguintes, basta utilizar o campo [body.propagate](#bodypropagate) como `true`.

`RESPONSE`: Para modificar o escopo de resposta do backend.

#### body.action

Campo obrigatório, do tipo string, responsável pela ação a ser tomada na modificação do body.

Valores aceitos se o body for JSON:

`ADD`: Adiciona a chave informada no campo [body.key](#bodykey) caso não exista, e agrega o valor informado no
campo [body.value](#bodyvalue).

`APD`: Acrescenta o valor informado no campo [body.value](#bodyvalue) caso a chave informada no
campo [body.key](#bodykey) exista.

`SET`: Defini o valor da chave informada no campo [body.key](#bodykey) pelo valor passado no
campo [body.value](#bodyvalue).

`RPL`: Substitui o valor da chave informada no campo [body.key](#bodykey) pelo valor passado no
campo [body.value](#bodyvalue) caso exista.

`REN`: Renomeia a chave informada no campo [body.key](#bodykey) pelo valor passado no
campo [body.value](#bodyvalue) caso exista.

`DEL`: Remove a chave informada no campo [body.key](#bodykey) caso exista.

Valores aceitos se o body for TEXTO:

`ADD`: Agrega o valor informado no campo [body.value](#bodyvalue) ao texto.

`APD`: Acrescenta o valor informado no campo [body.value](#bodyvalue) caso body não for vazio.

`RPL`: Irá substituir todos os valores semelhantes à chave informada no campo [body.key](#bodykey) pelo valor passado no
campo [body.value](#bodyvalue).

`DEL`: Remove todos os valores semelhantes à chave informada no campo [body.key](#bodykey).

#### body.key

Campo obrigatório, do tipo string, utilizado para indicar qual chave do cabeçalho deve ser modificada.

OBS: se torna opcional se seu body for do tipo TEXTO e [body.action](#bodyaction) tiver o valor `ADD`.

#### body.value

Campo obrigatório, do tipo string, utilizado como valor a ser usado para modificar a chave indicada no
campo [body.key](#bodykey).

Temos possibilidades de utilização de [valores dinâmicos](#valores-dinâmicos-para-modificação),
e de [variáveis de ambiente](#variáveis-de-ambiente) para esse campo.

OBS: se torna opcional apenas se [body.action](#bodyaction) tiver o valor `DEL`.

#### body.propagate

Campo opcional, do tipo booleano, o valor padrão é `false` indicando que o modificador não deve propagar essa mudança
em questão para os backends seguintes.

Caso informado como `true` essa modificação será propagada para os seguintes backends.

IMPORTANTE: Esse campo só é aceito se o [escopo](#bodyscope) tiver o valor `REQUEST`.

## JSON de tempo de execução

O GOPEN API Gateway quando iniciado, gera um arquivo JSON, baseado no [JSON de configuração](#json-de-configuração),
localizado na pasta `runtime` na raiz da sua aréa de trabalho.

Esse JSON, indica qual foi o entendimento da aplicação ao ler o [JSON de configuração](#json-de-configuração), todas
as [#variáveis de configuração](#variáveis-de-ambiente) já terão seus valores substituídos, caso exista.

Esse json também pode ser lido utilizando a rota estática [/settings](#settings).

## Rotas estáticas

O GOPEN API Gateway tem alguns endpoints estáticos, isto é, indepêndente de qualquer configuração feita, teremos
atualmente três endpoints cadastrados nas rotas do mesmo, veja abaixo cada um e suas responsabilidades:

### `/ping`

Endpoint para saber se a API Gateway está viva o path, retorna `404 (Not found)` se tiver off, e
`200 (OK)` se tiver no ar.

### `/version`

Endpoint que retorna a versão obtida na config [version](#version), retorna `404 (Not Found)` se não tiver sido
informado no [json de configuração](#json-de-configuração), caso contrário retorna o `200 (OK)` com o valor no body
como texto.

### `/settings`

Endpoint retorna algumas informações sobre o projeto, como versão, data da versão, quantidade de contribuintes e
um resumo de quantos endpoints, middlewares, backends e modifiers configurados no momento e o json de configuração
que está rodando ativamente.

```json
{
  "version": "v1.0.0",
  "version-date": "03/27/2024",
  "founder": "Gabriel Cataldo",
  "contributors": 1,
  "endpoints": 1,
  "middlewares": 0,
  "backends": 1,
  "modifiers": 0,
  "setting": {
    "port": 8080,
    "endpoints": [
      {
        "path": "/users/find/:key",
        "method": "GET",
        "backends": [
          {
            "hosts": [
              "http://192.168.1.8:8090"
            ],
            "path": "/users/find/:key",
            "method": "GET"
          }
        ]
      }
    ]
  }
}
```

## Variáveis de ambiente

As variáveis de ambiente podem ser fácilmente instânciadas utilizando o arquivo .env, na pasta indicada pelo ambiente
dinâmico de inicialização como mencionado no tópico [ESTRUTURA DE PASTAS](#estrutura-de-pastas).

Caso preferir inserir os valores utilizando docker-compose também funcionará corretamente, ponto é que a API
Gateway irá ler o valor gravado na máquina, independente de como foi inserido nela.

Os valores podem ser utilizados na configuração do JSON da API Gateway, basta utilizar a sintaxe `$NOME` como
um valor string, veja no exemplo abaixo.

Um trecho de um JSON de configuração, temo os seguintes valores:

```json
{
  "version": "$VERSION",
  "port": 8080,
  "hot-reload": true,
  "store": {
    "redis": {
      "address": "$REDIS_URL",
      "password": "$REDIS_PASSWORD"
    }
  },
  "timeout": "$TIMEOUT"
}
```

E na nossa máquina temos as seguintes variáveis de ambiente:

```dotenv
VERSION=1.0.0

REDIS_URL=redis-18885.c259.us-east-1-4.ec2.cloud.redislabs.com:18985
REDIS_PASSWORD=12345

TIMEOUT=5m
```

A API Gateway gera um arquivo de [JSON de tempo de execução](#json-de-tempo-de-execução) ao rodar a aplicação, veja o
resultado do mesmo após iniciar a aplicação:

```json
{
  "version": "1.0.0",
  "port": 8080,
  "hot-reload": true,
  "store": {
    "redis": {
      "address": "redis-18885.c259.us-east-1-4.ec2.cloud.redislabs.com:18985",
      "password": "12345"
    }
  },
  "timeout": "5m"
}
```

Vimos que todos os valores com a sintaxe `$NOME` foram substituídos pelos seus devidos valores, caso um valor
tenha sido mencionado por essa sintaxe, porém não existe nas variáveis de ambiente, o mesmo valor informado
será mantido.

### Observação

Nessa primeira versão podemos utilizar essa sintaxe apenas em campos do tipo string, caso queira
[contríbuir](#como-contríbuir) com uma solução que aceitará essa sintaxe independente do tipo, ficaríamos gratos!

## Valores dinâmicos para modificação

Podemos utilizar valores de requisição e resposta do tempo de execução do endpoint, conforme o mesmo foi configurado.
Esses valores podem ser obtidos por uma sintaxe específica, temos as seguintes possibilidades de obtenção desses
valores, veja:

### Requisição

Quando menciona a sintaxe `#request...` você estará obtendo os valores da requisição recebida, caso mencione o campo
`#request.history...` você estará obtendo os valores
do [histórico temporário de requisições](#histórico-temporário-de-requisições).

`#request.header...`

Esse trecho da sintaxe irá obter do cabeçalho da requisição o valor indicado, por exemplo,
`#request.header.X-Forwarded-For.0` irá obter o primeiro valor do campo `X-Forwarded-For` do cabeçalho da requisição
caso exista, substituindo a sintaxe pelo valor, o resultado foi `127.0.0.1`.

`#request.params...`

Esse trecho da sintaxe irá obter dos parâmetros da requisição o valor indicado, por exemplo,
`#request.params.id` irá obter o valor do campo `id` dos parâmetros da requisição caso exista,
substituindo a sintaxe pelo valor, o resultado foi `72761`.

`#request.query...`

Esse trecho da sintaxe irá obter dos parâmetros de busca da requisição o valor indicado, por exemplo,
`#request.query.email.0` irá obter o primeiro valor do campo `email` dos parâmetros de busca da requisição caso exista,
substituindo a sintaxe pelo valor, o resultado foi `gabrielcataldo.adm@gmail.com`.

`#request.body...`

Esse trecho da sintaxe irá obter do body da requisição o valor indicado, por exemplo,
`#request.body.deviceId` irá obter o valor do campo `deviceId` do body da requisição caso exista,
substituindo a sintaxe pelo valor, o resultado foi `991238`.

#### Histórico temporário de requisições

Caso queira obter os valores de requisição especificamente de um backend, basta adicionar o campo `history` na sintaxe,
e a posição do seu backend na história, lembrando que os middlewares `beforewares` e `afterwares` conta como um backend
no histórico, outro ponto é que só terá as informações de requisição neste campo se o backend já tenha sido processado,
ou o backend está em processamento. Veja como ficaria um exemplo abaixo:

No exemplo, eu tenho apenas um backend e o mesmo está sendo processado, então posso está utilizando a sintaxe:

`#request.history.0.params.phone`

Nesse outro exemplo de sintaxe temos três backends configurados e dois já foram processados, então podemos utilizar a
seguinte sintaxe no processo do terceiro backend:

`#request.history.1.body.id`

Nesses exemplos citados vemos que podemos obter o valor da requisição que esta sendo feita ou já foi feita por um
backend configurado, e que estão armazenados em um tipo de histórico temporário.

### Resposta

Quando menciona a sintaxe `#response...` você estará obtendo os valores da resposta da API Gateway preparada com base
na [lógica de resposta](#lógica-de-resposta), caso mencione o campo `#response.history...` você estará obtendo os
valores do [histórico temporário de respostas](#histórico-temporário-de-respostas).

`#response.statusCode`

Essa sintaxe irá obter o código de status HTTP da resposta atual do endpoint. Caso queira obter da
resposta de um backend específico use o [histórico temporário de respostas](#histórico-temporário-de-respostas).

`#response.header...`

Esse trecho da sintaxe irá obter do cabeçalho de resposta do endpoint o valor indicado, por exemplo,
`#response.header.X-Trace-Id` irá obter os valores em lista do campo `X-Trace-Id` do cabeçalho de resposta caso exista,
substituindo a sintaxe pelo valor, o resultado foi `["4ae6c92d16089e521626", "5ae7c93d16089e521628"]`.

`#response.body...`

Esse trecho da sintaxe irá obter do body de resposta do endpoint o valor indicado, por exemplo,
`#response.body.role` irá obter o valor do campo `role` do body da requisição caso exista,
substituindo a sintaxe pelo valor, o resultado foi `ADMIN`.

#### Histórico temporário de respostas

Caso queira obter os valores de resposta especificamente de um backend, basta adicionar o campo `history` na sintaxe,
e a posição do seu backend na história, lembrando que os middlewares `beforewares` e `afterwares` conta como um backend
no histórico, outro ponto é que só terá as informações de resposta neste campo se o backend já tenha sido processado.
Veja como ficaria um exemplo abaixo:

No exemplo, eu tenho apenas um backend e o mesmo foi processado, então posso está utilizando a sintaxe:

`#response.history.0.header.X-Value.0`

Nesse outro exemplo de sintaxe temos três backends configurados e dois já foram processados, então podemos utilizar a
seguinte sintaxe no processo do terceiro backend:

`#request.history.1.body.users.0`

Nesses exemplos citados vemos que podemos obter o valor da resposta de um backend que já foi processado,
e que estão armazenados em um tipo de histórico temporário.

### Importante

Você pode utilizar com base nesses campos,
a [sintaxe de JSON path](https://github.com/tidwall/gjson/blob/master/README.md#path-syntax) que se enquadra em seus
valores, apenas se lembre que, os objetos header, query são mapas de lista de string, e o params é um mapa de string.

### Veja também

- Configuração de modificadores de backend clicando [aqui](#backendmodifiers).

- Aprenda na prática como utilizar os valores dinâmicos para modificação usando o
  projeto [playground](https://github.com/GabrielHCataldo/gopen-gateway-playground) que já vem com alguns exemplos de
  modificadores com valores dinâmicos.

## Observabilidade

A API Gateway ao receber uma requisição irá acrescentar o IP do cliente no cabeçalho `X-Forwarded-For` das requisições
subjacentes, e também adiciona, caso não exista, um valor único gerado ao cabeçalho chamado `X-Trace-Id` para melhor
observabilidade da requisição recebida tanto nos logs da API Gateway como nos microserviços subjacentes.

## Lógica de resposta

Quando utilizamos uma API Gateway nos perguntamos, como será retornado ao meu cliente
a resposta desse endpoint configurado?

Para facilitar o entendimento criamos esse tópico para resumir a lógica de resposta da nossa API Gateway,
então vamos começar.

### Como funciona?

A API Gateway foi desenvolvida com uma inteligência e flexibilidade ao responder um endpoint, ela se baseia em dois
pontos importantes, primeiro, na quantidade de respostas de serviços backends que foram processados, e segundo, nos
campos de customização da resposta configurados nos objetos [endpoint](#endpointcomment) e [backend](#backendname).
Vamos ver alguns exemplos abaixo para melhor entendimento.

#### Único backend

Nesse exemplo trabalharemos apenas com um único backend, veja como a API Gateway se comportará ao responder
a esse cenário

Json de configuração

```json
{
  "$schema": "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json",
  "port": 8080,
  "endpoints": [
    {
      "path": "/users/find/:key",
      "method": "GET",
      "backends": [
        {
          "hosts": [
            "$USER_SERVICE_URL"
          ],
          "path": "/users/find/:key",
          "method": "GET"
        }
      ]
    }
  ]
}
```

Ao processar esse endpoint a resposta da API Gateway foi

```text
HTTP/1.1 200 OK
```

Cabeçalho ([Veja sobre os cabeçalhos de resposta aqui](#cabeçalho-de-resposta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Tue, 23 Apr 2024 11:37:26 GMT
Content-Length: 620
```

Corpo

```json
{
  "id": "6499b8826493f85e45eb3794",
  "name": "Gabriel Cataldo",
  "birthDate": "1999-01-21T00:00:00Z",
  "gender": "MALE",
  "currentPage": "HomePage",
  "createdAt": "2023-06-26T16:10:42.265Z",
  "updatedAt": "2024-03-10T20:19:03.452Z"
}
```

Vimos que no exemplo a API Gateway serviu como um proxy redirecionando a requisição para o serviço backend configurado e
espelhando seu body de resposta, e agregando seus valores no cabeçalho de resposta tenha.

Nesse mesmo exemplo vamos forçar um cenário de infelicidade na resposta do backend, veja:

```text
HTTP/1.1 404 Not Found
```

Cabeçalho ([Veja sobre os cabeçalhos de resposta aqui](#cabeçalho-de-resposta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: false
Date: Tue, 23 Apr 2024 21:56:33 GMT
Content-Length: 235
```

Corpo

```json
{
  "file": "datastore/user.go",
  "line": 227,
  "endpoint": "/users/find/gabrielcataldo.adma@gmail.com",
  "message": "user not found"
}
```

Neste caso a API Gateway também espelhou a resposta da única chamada de backend do endpoint.

#### Utilizando middlewares

Nesse exemplo, vamos utilizar os middlewares de [beforewares](#endpointbeforewares) e [afterwares](#endpointafterwares),
como esses backends são omitidos ao cliente final se tiverem sucesso, vamos simular uma chamada com o device bloqueado
para que o [beforeware](#endpointbeforewares) retorne um erro, e depois um [afterware](#endpointafterwares) que
responderá também um erro, pois não existe, vamos lá!

Json de configuração

```json
{
  "$schema": "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json",
  "port": 8080,
  "middlewares": {
    "save-device": {
      "hosts": [
        "$DEVICE_SERVICE_URL"
      ],
      "path": "/devices",
      "method": "PUT"
    },
    "increment-attempts": {
      "hosts": [
        "$SECURITY_SERVICE_URL"
      ],
      "path": "/attempts",
      "method": "POST"
    }
  },
  "endpoints": [
    {
      "path": "/users/find/:key",
      "method": "GET",
      "beforewares": [
        "save-device"
      ],
      "afterwares": [
        "increment-attempts"
      ],
      "backends": [
        {
          "hosts": [
            "$USER_SERVICE_URL"
          ],
          "path": "/users/find/:key",
          "method": "GET"
        }
      ]
    }
  ]
}
```

Ao processar esse endpoint de exemplo simulando o erro na chamada de [beforeware](#endpointbeforewares) a resposta da
API Gateway foi

```text
HTTP/1.1 403 Forbidden
```

Cabeçalho ([Veja sobre os cabeçalhos de resposta aqui](#cabeçalho-de-resposta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: true
Date: Tue, 23 Apr 2024 23:02:09 GMT
Content-Length: 154
```

Corpo

```json
{
  "file": "service/device.go",
  "line": 49,
  "endpoint": "/devices",
  "message": "unprocessed entity: device already exists and is not active"
}
```

Vimos que a resposta foi o espelho do retorno do beforeware `save-device`, pois como o mesmo retornou
falha `403 (Forbidden)`, o endpoint abortou, não chamando os backends seguintes, lembrando que você
pode configurar os códigos de status HTTP que vão ser abortados pelo seu endpoint, basta preencher o
campo [endpoint.abort-if-status-codes](#endpointabort-if-status-codes).

No seguinte exemplo iremos forçar um erro no afterware `increment-attempts` a da API Gateway resposta foi

```text
HTTP/1.1 404 Not Found
```

Cabeçalho ([Veja sobre os cabeçalhos de resposta aqui](#cabeçalho-de-resposta))

```text
Content-Type: text/plain
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: false
Date: Tue, 23 Apr 2024 23:16:57 GMT
Content-Length: 18
```

Corpo

```text
404 page not found
```

Vimos que a resposta também foi o espelho do retorno do afterware `increment-attempts`, por mais que seja a última
chamada de um serviço backend do endpoint, pois caiu na regra de resposta abortada, então, todas as outras respostas
dos outros backends foram ignoradas e apenas foi retornado a resposta do backend abortado.

Veja mais sobre a [resposta abortada](#resposta-abortada).

#### Múltiplos backends

Nesse exemplo iremos trabalhar com três [backends](#endpointbackends) principais no endpoint, então, vamos lá!

Json de configuração

```json
{
  "$schema": "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json",
  "port": 8080,
  "endpoints": [
    {
      "path": "/users/find/:key",
      "method": "GET",
      "backends": [
        {
          "name": "user",
          "hosts": [
            "$USER_SERVICE_URL"
          ],
          "path": "/users/find/:key",
          "method": "GET"
        },
        {
          "hosts": [
            "$DEVICE_SERVICE_URL"
          ],
          "path": "/devices",
          "method": "PUT"
        },
        {
          "name": "version",
          "hosts": [
            "$USER_SERVICE_URL"
          ],
          "path": "/version",
          "method": "GET"
        }
      ]
    }
  ]
}
```

No exemplo iremos executar os três backend com sucesso, a API Gateway respondeu

```text
HTTP/1.1 200 OK
```

Cabeçalho ([Veja sobre os cabeçalhos de resposta aqui](#cabeçalho-de-resposta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Tue, 23 Apr 2024 23:49:12 GMT
Content-Length: 755
```

Corpo

```json
[
  {
    "ok": true,
    "code": 200,
    "id": "6499b8826493f85e45eb3794",
    "name": "Gabriel Cataldo",
    "birthDate": "1999-01-21T00:00:00Z",
    "gender": "MALE",
    "currentPage": "HomePage",
    "createdAt": "2023-06-26T16:10:42.265Z",
    "updatedAt": "2024-03-10T20:19:03.452Z"
  },
  {
    "ok": true,
    "code": 200,
    "id": "661535275d6fc736d831c754",
    "usersId": [
      "6499b8826493f85e45eb3793"
    ],
    "status": "ACTIVE",
    "createdAt": "2024-04-09T12:31:35.907Z",
    "updatedAt": "2024-04-23T23:49:12.759Z"
  },
  {
    "ok": true,
    "code": 200,
    "version": "v1.0.0"
  }
]
```

Temos alguns pontos nesse exemplo que vale ressaltar, primeiro com o formato, a API Gateway entendeu que seu endpoint
tem múltiplas respostas e não foi utilizado o campo [endpoint.aggregate-responses](#endpointaggregate-responses)
com o valor `true`, então ela lista as respostas como JSON acrescentando os seguintes campos:

`ok`: Indica se a resposta do backend em questão teve o código de status HTTP entre `200` e `299`.

`code`: É preenchido com código de status HTTP respondido pelo seu backend.

Esses campos são apenas acrescentado se houver múltiplas respostas e o
campo [endpoint.aggregate-responses](#endpointaggregate-responses) não for informado com o valor `true`.

Segundo ponto a destacar é no trecho `"version": "v1.0.0"` do último backend, o mesmo respondeu apenas um texto no body
de resposta que foi `v1.0.0`, porém para esse cenário como foi mencionado, a API Gateway força o parse desse valor para
um JSON, adicionando um novo campo com o nome obtido na configuração [backend.name](#backendname) e com o valor do
mesmo.

Terceiro ponto é sobre o código de status HTTP, o mesmo é retornado pela maior frequência, isto é, se temos três
retornos `200 OK` como no exemplo a API Gateway também retornará esse código. Se tivermos um retorno igualitário o
último código de status HTTP retornado será considerado, veja os cenários possíveis dessa lógica

```json
[
  {
    "ok": true,
    "code": 204
  },
  {
    "ok": true,
    "code": 200
  },
  {
    "ok": true,
    "code": 201
  }
]
```

a API Gateway responderá `201 Created`.

```json
[
  {
    "ok": true,
    "code": 100
  },
  {
    "ok": true,
    "code": 100
  },
  {
    "ok": true,
    "code": 201
  }
]

```

a API Gateway responderá `100 Continue`.

Quarto ponto a ser destacado, é que como o endpoint tem múltiplas respostas, consequentemente temos múltiplos cabeçalhos
de resposta, a API Gateway irá agregar todos os campos e valores para o cabeçalho da resposta final, veja mais sobre o
comportamento do cabeçalho de resposta clicando [aqui](#cabeçalho-de-resposta).

Último ponto a ser destacado, é que caso um desses retornos de backend entre no cenário em que o endpoint aborta a
resposta, ele não seguirá nenhuma diretriz mostrada no tópico em questão e sim
[lógica de resposta abortada](#resposta-abortada).

#### Múltiplos backends agregados

Nesse exemplo iremos utilizar uma configuração parecida com JSON de configuração do exemplo acima, porém com
campo [endpoint](#endpointaggregate-responses) com o valor `true`.

Json de configuração

```json
{
  "$schema": "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json",
  "port": 8080,
  "endpoints": [
    {
      "path": "/users/find/:key",
      "method": "GET",
      "aggregate-responses": true,
      "backends": [
        {
          "name": "user",
          "hosts": [
            "$USER_SERVICE_URL"
          ],
          "path": "/users/find/:key",
          "method": "GET"
        },
        {
          "hosts": [
            "$DEVICE_SERVICE_URL"
          ],
          "path": "/devices",
          "method": "PUT"
        },
        {
          "name": "version",
          "hosts": [
            "$USER_SERVICE_URL"
          ],
          "path": "/version",
          "method": "GET"
        }
      ]
    }
  ]
}
```

Ao processarmos o endpoint a resposta da API Gateway foi

```text
HTTP/1.1 200 OK
```

Cabeçalho ([Veja sobre os cabeçalhos de resposta aqui](#cabeçalho-de-resposta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Wed, 24 Apr 2024 10:57:31 GMT
Content-Length: 665
```

Corpo

```json
{
  "id": [
    "6499b8826493f85e45eb3794",
    "661535275d6fc736d831c754"
  ],
  "name": "Gabriel Cataldo",
  "gender": "MALE",
  "currentPage": "HomePage",
  "lastSeenAt": "2024-02-19T11:43:27.324Z",
  "createdAt": [
    "2024-04-09T12:31:35.907Z",
    "2023-06-26T16:10:42.265Z"
  ],
  "updatedAt": [
    "2024-04-24T11:04:32.184Z",
    "2024-03-10T20:19:03.452Z"
  ],
  "usersId": [
    "6499b8826493f85e45eb3793"
  ],
  "status": "ACTIVE",
  "version": "v1.0.0"
}
```

Vimos a única diferença de resposta do tópico [Múltiplos backends](#múltiplos-backends) é que ele agregou os valores
de todas as respostas em um só JSON, e os campos que se repetiram foram agregados os valores em lista.

As demais regras como código de status HTTP, a conversão forçada para JSON, entre outras, seguem as mesmas regras
mencionadas no tópico [Múltiplos backends](#múltiplos-backends).

No exemplo podemos deixar a resposta agregada um pouco mais organizada, com isso vamos alterar o trecho do nosso
segundo backend adicionando o campo [backend.extra-config.group-response](#backendextra-configgroup-response) com o
valor `true` e dar um nome a ele, veja o trecho do JSON de configuração modificado

```json
{
  "name": "device",
  "hosts": [
    "$DEVICE_SERVICE_URL"
  ],
  "path": "/devices",
  "method": "PUT",
  "extra-config": {
    "group-response": true
  }
}
```

Ao processar novamente o endpoint obtivemos a seguinte resposta

```text
HTTP/1.1 200 OK
```

Cabeçalho ([Veja sobre os cabeçalhos de resposta aqui](#cabeçalho-de-resposta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Wed, 24 Apr 2024 11:23:07 GMT
Content-Length: 697
```

Corpo

```json
{
  "id": "6499b8826493f85e45eb3793",
  "name": "Gabriel Cataldo",
  "birthDate": "1999-01-21T00:00:00Z",
  "gender": "MALE",
  "currentPage": "HomePage",
  "lastSeenAt": "2024-02-19T11:43:27.324Z",
  "createdAt": "2023-06-26T16:10:42.265Z",
  "updatedAt": "2024-03-10T20:19:03.452Z",
  "device": {
    "id": "661535275d6fc736d831c754",
    "usersId": [
      "6499b8826493f85e45eb3793"
    ],
    "status": "ACTIVE",
    "createdAt": "2024-04-09T12:31:35.907Z",
    "updatedAt": "2024-04-24T11:23:07.832Z"
  },
  "version": "v1.0.0"
}
```

Com essa configuração vimos que nossa resposta agregada ficou mais organizada, e como é importante entender sobre
o [json de configuração](#json-de-configuração) e seus campos, para que o GOPEN API Gateway atenda melhor suas
necessidades.

### Resposta abortada

Para uma resposta ser abortada pela API Gateway, um dos backends configurados do endpoint tanto middlewares
como os principais, ao serem processados, na sua resposta, o código de status HTTP precisa seguir valores
no campo [endpoint.abort-if-status-codes](#endpointabort-if-status-codes) do próprio endpoint.

É importante que você entenda esses dois pontos:

- Ao abortar a resposta do backend, a API Gateway irá espelhar apenas a resposta do mesmo, código de status, cabeçalho e
  corpo, sendo assim, as outras respostas já processadas serão ignoradas.
- Mesmo o backend sendo abortado, os [modificadores](#backendmodifiers) serão processados.

Indicamos utilizar essa configuração apenas quando algo fugiu do esperado, como, por exemplo, uma resposta de erro
`500 (Internal server error)`.

### Cabeçalho de resposta

Na resposta, a API Gateway com exceção dos campos `Content-Length`, `Content-Type` e `Date` agrega todos valores de
cabeçalho respondidos pelos backends configurados no endpoint, indepêndente da quantidade de backends, isso inclui
também os [middlewares](#middlewares).

#### Campos de cabeçalho padrão

Também são adicionados até quatro campos no cabeçalho veja abaixo sobre os mesmos:

`X-Gopen-Cache`: Caso a resposta do endpoint não seja "fresca", isto é, foi utilizado a resposta armazenada em cache,
é retornado o valor `true`, caso contrário retorna o valor `false`.

`X-Gopen-Cache-Ttl`: Caso a resposta do endpoint tenha sido feita utilizando o armazenamento em cache, ele retorna a
duração do tempo de vida restante desse cache, caso contrário o campo não é retornado.

`X-Gopen-Complete`: Caso todos os backends tenham sido processados pelo endpoint é retornado o valor `true`, caso
contrário é retornado o valor `false`.

`X-Gopen-Success`: Caso todos os backends tenham retornado sucesso, isto é, o código de status HTTP de resposta entre
`200` a `299`, ele retorna o valor `true`, caso contrário o valor `false`.

Lembrando que se a resposta de um backend for [abortada](#resposta-abortada), apenas o header do mesmo é agregado e
considerado as regras dos campos acima.

Agora vamos ver alguns exemplos de cabeçalho de retorno:

#### Campos únicos de cabeçalho

Cabeçalho de resposta do backend 1:

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521626
X-MS: api-user
Date: Wed, 24 Apr 2024 11:23:07 GMT
Content-Length: 102
```

Cabeçalho de resposta do endpoint

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521626
X-MS: api-user
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Wed, 24 Apr 2024 11:23:08 GMT
Content-Length: 102
```

Vimos que no exemplo foram adicionados os [campos padrões](#campos-de-cabeçalho-padrão), e agregado os valores do
cabeçalho de resposta, que foram `X-Value-Id` e `X-MS`.

#### Campos duplicados de cabeçalho

Cabeçalho de resposta do backend 1:

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521626
X-MS: api-user
Date: Wed, 24 Apr 2024 11:23:07 GMT
Content-Length: 102
```

Cabeçalho de resposta do backend 2:

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521638
X-MS: api-device
X-MS-Success: true
Date: Wed, 24 Apr 2024 11:23:08 GMT
Content-Length: 402
```

Cabeçalho de resposta do endpoint

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521626, 4ae6c92d16089e521638
X-MS: api-user, api-device
X-MS-Success: true
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Wed, 24 Apr 2024 11:23:09 GMT
Content-Length: 504
```

Vimos que no exemplo também foram adicionados os [campos padrões](#campos-de-cabeçalho-padrão), e agregado os valores do
cabeçalho de resposta, que foram `X-Value-Id`, `X-MS` e `X-MS-Success`, vale ressaltar que os campos que se repetiram
foram agrupados e separados por vírgula.

### Respostas padrões

Toda API Gateway tem suas respostas padrão para cada cenário de erro, então iremos listar abaixo cada
cenário e sua respectiva resposta HTTP:

#### 413 (Request Entity Too Large)

Esse cenário acontece quando o tamanho do corpo de requisição é maior do que o permitido para o endpoint, utilizando a
configuração [limiter.max-body-size](#limitermax-header-size) para corpo normal
e [limiter.max-multipart-memory-size](#limitermax-multipart-memory-size) para envio do tipo `form-data`. Você pode
customizar essa configuração para um endpoint específico utilizando o campo [endpoint.limiter](#endpointlimiter).

Cabeçalho

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 11:56:06 GMT
Content-Length: 170
```

Corpo

```json
{
  "file": "infra/size_limiter.go",
  "line": 92,
  "endpoint": "/users",
  "message": "payload too large error: permitted limit is 1.0B",
  "timestamp": "2024-04-26T08:56:06.628636-03:00"
}
```

#### 429 (Too many requests)

Esse cenário acontece quando o limite de requisições são atingidas por um determinado IP, esse limite é definido na
configuração [limiter.rate](#limiterrate). Você pode customizar essa configuração para um endpoint
específico utilizando o campo [endpoint.limiter](#endpointlimiter).

Cabeçalho

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 12:12:53 GMT
Content-Length: 177
```

Corpo

```json
{
  "file": "infra/rate_limiter.go",
  "line": 100,
  "endpoint": "/users",
  "message": "too many requests error: permitted limit is 1 every 1s",
  "timestamp": "2024-04-26T09:12:53.501804-03:00"
}
```

#### 431 (Request Header Fields Too Large)

Esse cenário acontece quando o tamanho do header é maior do que o permitido para o endpoint, utilizando a
configuração [limiter.max-header-size](#limitermax-header-size). Você pode customizar essa configuração para um endpoint
específico utilizando o campo [endpoint.limiter](#endpointlimiter).

Cabeçalho

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 11:39:53 GMT
Content-Length: 186
```

Corpo

```json
{
  "file": "infra/size_limiter.go",
  "line": 80,
  "endpoint": "/multiple/backends/:key",
  "message": "header too large error: permitted limit is 1.0B",
  "timestamp": "2024-04-26T08:39:53.944055-03:00"
}
```

#### 500 (Internal server error)

Esse cenário é específico quando algum erro inesperado ocorreu com a API Gateway, caso isso aconteça relate
o problema [aqui](https://github.com/GabrielHCataldo/gopen-gateway/issues) mostrando a resposta e o log impresso no
terminal de execução.

Cabeçalho

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 12:38:16 GMT
Content-Length: 183
```

Corpo

```json
{
  "file": "middleware/panic_recovery.go",
  "line": 27,
  "endpoint": "/users",
  "message": "gateway panic error occurred! detail: runtime error: invalid memory address or nil pointer dereference",
  "timestamp": "2024-04-26T09:42:23.938997-03:00"
}
```

#### 502 (Bad Gateway)

Esse cenário acontece quando ao tentar se comunicar com o backend, e ocorre alguma falha de comunicação com o mesmo.

Cabeçalho

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Thu, 25 Apr 2024 01:07:36 GMT
Content-Length: 277
```

Corpo

```json
{
  "file": "infra/rest.go",
  "line": 69,
  "endpoint": "/users/find/:key",
  "message": "bad gateway error: Get \"http://192.168.1.8:8090/users/find/gabrielcataldo.adm@gmail.com\": dial tcp 192.168.1.8:8090: connect: connection refused",
  "timestamp": "2024-04-24T22:07:36.558851-03:00"
}
```

#### 504 (Gateway Timeout)

Esse cenário acontece quando o endpoint excede o limite do tempo configurado no campo [timeout](#timeout). Você pode
customizar essa configuração para um endpoint específico utilizando o campo [endpoint.timeout](#endpointtimeout).

Cabeçalho

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 13:29:55 GMT
Content-Length: 150
```

Corpo

```json
{
  "file": "middleware/timeout.go",
  "line": 81,
  "endpoint": "/users/version",
  "message": "gateway timeout: 5m",
  "timestamp": "2024-04-26T10:29:55.908526-03:00"
}
```

# Usabilidade

Use o projeto [playground](https://github.com/GabrielHCataldo/gopen-gateway-playground) para começar a explorar e
utilizar na prática o GOPEN API Gateway!

# Como contríbuir?

Ficamos felizes quando vemos a comunidade se apoiar, e projetos como esse, está de braços abertos para receber
suas ideias, veja abaixo como participar.

## Download

Para conseguir rodar o projeto primeiro faça o download da [linguagem Go](https://go.dev/dl/)
versão 1.22 ou superior na sua máquina.

Com o Go instalado na sua máquina, faça o pull do projeto

```text
git pull https://github.com/GabrielHCataldo/gopen-gateway.git
```

Depois abra o mesmo usando o próprio terminal com a IDE de sua preferência

Goland:

```text
goland gopen-gateway
```

VSCode:

```text
code gopen-gateway
```

## Air

Caso queira utilizar o hot reload para otimizar no desenvolvimento, o projeto já tem o arquivo
`.air.toml` como exemplo para você utilizar como quiser.

Veja a documentação sobre esse projeto fantástico clicando [aqui](https://github.com/cosmtrek/air).

## Gitflow

Para inicializar o desenvolvimento, você pode criar uma branch a partir da main, para um futuro
PR para a mesma.

# Agradecimentos

Esse projeto teve apoio de bibliotecas fantásticas, esse trecho dedico a cada uma listada
abaixa:

- [go-errors](https://github.com/GabrielHCataldo/go-errors)
- [go-helper](https://github.com/GabrielHCataldo/go-helper)
- [go-redis-template](https://github.com/GabrielHCataldo/go-redis-template)
- [validator](https://github.com/go-playground/validator)
- [fsnotify](https://github.com/fsnotify/fsnotify)
- [gin](https://github.com/gin-gonic/gin)
- [gjson](https://github.com/tidwall/gjson)
- [sjson](https://github.com/tidwall/sjson)
- [uuid](https://github.com/google/uuid)
- [ttlcache](https://github.com/jellydator/ttlcache)
- [godotenv](https://github.com/joho/godotenv)
- [gojsonschema](https://github.com/xeipuuv/gojsonschema)
- [air](https://github.com/cosmtrek/air)

Obrigado por contribuir para a comunidade Go e facilitar o desenvolvimento desse projeto.

# Licença Apache 2.0

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway?ref=badge_large&issueType=license)
