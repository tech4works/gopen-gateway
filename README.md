<img src="assets/logo.png" alt="">

[![Project status](https://img.shields.io/badge/version-v1.0.0-orange.svg)](https://github.com/GabrielHCataldo/gopen-gateway/releases/tag/v1.0.0)
[![Playground](https://img.shields.io/badge/%F0%9F%8F%90-playground-9900cc.svg)](https://github.com/GabrielHCataldo/gopen-gateway)
[![Docker](https://badgen.net/badge/icon/docker?icon=docker&label)](https://hub.docker.com/r/gabrielhcataldo/gopen-gateway)
[![Open Source Helpers](https://www.codetriage.com/gabrielhcataldo/gopen-gateway/badges/users.svg)](https://www.codetriage.com/gabrielhcataldo/gopen-gateway)
[![TODOs](https://img.shields.io/badge/TODOs-13-red.svg?style=flat)](https://github.com/GabrielHCataldo/gopen-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/GabrielHCataldo/gopen-gateway)](https://goreportcard.com/report/github.com/GabrielHCataldo/gopen-gateway)
[![GoDoc](https://godoc.org/github/GabrielHCataldo/gopen-gateway?status.svg)](https://pkg.go.dev/github.com/GabrielHCataldo/gopen-gateway/helper)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/GabrielHCataldo/gopen-gateway)](https://github.com/GabrielHCataldo/gopen-gateway/blob/main/go.mod)
[![HitCount](https://hits.dwyl.com/GabrielHCataldo/gopen-gateway.svg)](http://hits.dwyl.com/GabrielHCataldo/gopen-gateway)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway.svg?type=small)](https://app.fossa.com/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway?ref=badge_small)

![Brazil](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/br.png)
[Portuguese](https://github.com/GabrielHCataldo/gopen-gateway/blob/main/README.pt-br.md) |
![Spain](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/es.png)
[Spanish](https://github.com/GabrielHCataldo/gopen-gateway/blob/main/README.es.md)

The GOPEN project was created with the aim of helping developers to have a robust and easy to handle API Gateway, with
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
- Run the modifiers in the context you prefer, before a backend request or after, it's your decision.
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

## Folder structure

In the project structure, you will need to have a root folder called "gopen". Inside this folder, you should have
subfolders containing the names of your environments. You can name these subfolders as you wish. Each of these
subfolders should contain at least one ".json" configuration file for the API Gateway. For example, the structure might
look like this:

### GO project

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

### Docker image project

    nome-do-seu-projeto
    | - docker-compose.yml
    | - gopen
      | - dev
      |   - .json
      |   - .env
      | - prd
      |   - .json
      |   - .env

## Configuration JSON

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

#### limiter.max-body-size

Optional field, it is of type `byteUnit`, default value is `3MB`, field is responsible for limiting the size of the
request body.

If the size of the body exceeds the informed value, the API Gateway will abort the request with status
code `413 (Request entity too large)`.

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

Optional field, it is of type `byteUnit`, default value is `5MB`, this field is responsible for limiting the size of the
multipart/form body of the request, usually used for sending files, images, etc.

If the size of the body exceeds the informed value, the API Gateway will abort the request with status
code `413 (Request entity too large)`. Learn more about this return by clicking [here](#413-request-entity-too-large).

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

Optional field, if informed, the `capacity` field becomes mandatory, this object is responsible for limiting the request
rate by IP, this limit is imposed by obtaining the maximum capacity by the `capacity` field for X duration, informed on
the `every` field.

If this capacity is exceeded, the API Gateway will abort the request for security, returning `429 (Too Many Requests)`.
Learn more about this return by clicking [here](#429-too-many-requests).

#### limiter.rate.capacity

Optional field, if the rate object is informed, it becomes mandatory, the default value is `5`, and the minimum that can
be informed is `1`, it indicates the maximum capacity of requests.

#### limiter.rate.every

Optional field, the default value is `1 second`, it indicates the value of the duration of the verification of the
maximum capacity of requests.

### security-cors

Campo opcional, usado para segurança do CORS da API Gateway, todos os campos por padrão são vazios, não restringindo
os valores de origin, methods e headers.

Caso queira restringir, e a requisição não corresponda com as configurações impostas, a API Gateway por segurança
irá abortar a requisição retornando `403 (Forbidden)`.

#### security-cors.allow-origins

Campo opcional, do tipo lista de string, os itens da lista precisam indicar quais IPs de origem a API Gateway
permite receber nas requisições.

#### security-cors.allow-methods

Optional field, of the list of string type, the items in the list need to indicate which HTTP methods the API Gateway
allows to receive in the requests.

#### security-cors.allow-headers

Optional field, of the list of string type, the items in the list need to indicate which HTTP header fields the API
Gateway allows to receive in the requests.

### middlewares

Optional field, responsible for configuring your application middlewares, it's a map with string keys mentioning the
name of your middleware, this name can be used in your [endpoints](#endpoints) as `beforewares` and `afterwares`.

The key value is a [backend](#backendname) object, however, with a note, this object will have its response, if
not [aborted](#aborted-response), automatically omitted by the endpoint, only its response header may be returned, if
any, since middleware responses are not displayed to the final HTTP client, however, its response will be stored
throughout the HTTP request made at the endpoint, being able to have its request and response values obtained and
manipulated.

For example, a `beforeware` when mentioned in the endpoint, it will be used as a pre-request middleware, that is, it
will be called before the main backends of the endpoint, so we can, for example, have a device handling middleware, as
in the configuration JSON above, where it will call this middleware backend configured in the endpoint as `beforewares`,
validating and saving the device from header information of the request, if the backend responds a failure status code,
in the example, the gateway will abort all following backends returning what the device backend responded, if it has
returned a successful status code, it will modify the header of all subsequent requests (`propagate:true`), adding
the `X-Device-Id` field, with the value of the id of the body of the backend's own response, which can be used in the
other following backends of the endpoint.

To learn more about the `modifiers` [see](#backendmodifiers).

To better understand this powerful tool in practice, see the examples of middlewares used as `beforewares`
and `afterwares` made in the [playground](https://github.com/GabrielHCataldo/gopen-gateway-playground) project.

### endpoints

Mandatory field, it's a list of objects, represents each endpoint of the API Gateway that will be registered to listen
and serve HTTP requests.

See how the basic flow of an endpoint works in the image below:

<img src="assets/basic-endpoint-flowchart-en.png" alt="basic endpoint flow">

Below we will list and explain each field of this very important object:

### endpoint.@comment

Optional field, of string type, free field for annotations related to your endpoint.

### endpoint.path

Mandatory field, of string type, responsible for the URI path of the endpoint, for example `"/users/:id"`.

If you want to have dynamic parameters in this endpoint, just use the standard `":parameter name"` for
example `"/users/:id/status/:status"`, the API Gateway will understand that we will have 2 dynamic parameters of this
endpoint, these values can be passed on to the underlying backends.

Example using the dynamic parameter for the underlying backends:

- Endpoint
    - path: `"/users/:id/status/:status"`
    - result: `"/users/1/status/removed"`
- Backend 1
    - path: `"/users/:id"`
    - result: `"/users/1"`
- Backend 2
    - path: `"/users/:id/status/:status"`
    - result: `"/users/1/status/removed"`

In the example above, we see that the parameter can be used as you wish as a path in the backend requests of the
endpoint in question.

### endpoint.method

Mandatory field, of string type, responsible for defining the HTTP method for which the endpoint will be registered.

#### endpoint.timeout

It's similar to the [timeout](#timeout) field, however, it will be applied only to the endpoint in question.

If omitted, it will inherit the value of the [timeout](#timeout) field.

### endpoint.cache

Optional field, of object type, by default it will come empty only with the `enabled` field filled with the
value `false`.

If informed, the `enabled` field becomes mandatory, the other fields, if omitted, will inherit from the [cache](#cache)
configuration at root if it exists and if filled.

If by chance, you have omitted the `duration` field both in the current configuration and in the [cache](#cache)
configuration at root, the `enabled` field is ignored, always considering it as `false` as the cache duration was not
informed in both settings.

#### endpoint.cache.enabled

Mandatory field, of boolean type, indicates whether you want to have cache in your endpoint, `true` for enabled, `false`
for not enabled.

If `true` but the `duration` field is not informed in the current configuration nor at the [root](#cache), this field
will be ignored, always considering it as `false`.

#### endpoint.cache.ignore-query

Optional field, of boolean type, if not informed the default value is `false`.

If the value is `true`, the API Gateway will ignore the URL query parameters when creating the storage key, otherwise
it will consider the URL query parameters, alphabetically ordering the keys and values.

#### endpoint.cache.duration

It's similar to the [cache.duration](#cacheduration) field, however, it will be applied only to the endpoint in
question.

If omitted, it will inherit the value of the [cache.duration](#cacheduration) field.

If it is omitted in both configurations, the `enabled` field will be ignored, always considering it as `false`.

#### endpoint.cache.strategy-headers

It's similar to the [cache.strategy-headers](#cachestrategy-headers) field, however, it will be applied only to the
endpoint in question.

If omitted, it will inherit the value of the [cache.strategy-headers](#cachestrategy-headers) field.

If informed as empty, the value will not be inherited, however, the [default](#cachestrategy-headers) value will be
applied to the endpoint in question.

#### endpoint.cache.only-if-status-codes

It's similar to the [cache.only-if-status-codes](#cacheonly-if-status-codes) field, however, it will be applied only to
the endpoint in question.

If omitted, it will inherit the value of the [cache.only-if-status-codes](#cacheonly-if-status-codes) field.

If informed as empty, the value will not be inherited, however, the [default](#cacheonly-if-status-codes) value will be
applied to the endpoint in question.

#### endpoint.cache.allow-cache-control

It's similar to the [cache.allow-cache-control](#cacheallow-cache-control) field, however, it will be applied only to
the endpoint in question.

If omitted, it will inherit the value of the [cache.allow-cache-control](#cacheallow-cache-control) field.

#### endpoint.limiter

It's similar to the [limiter](#limiter) field, however, it will be applied only to the endpoint in question.

If omitted, it will inherit the value of the [limiter](#limiter) field.

#### endpoint.response-encode

Optional field, string type, the default value is empty, indicating that the endpoint response will be encoded following
the API Gateway [response logic](#response-logic), without forcing the indicated encoding.

```text
- Accepted values:
  - JSON 
  - XML
  - TEXT
```

#### endpoint.aggregate-responses

Optional field, of boolean type, the default value is `false`, indicating that the endpoint response will not be
aggregated.

If informed with the value `true` and there is more than one response from the backends mentioned in the endpoint, it
will aggregate the backend responses. See more about API Gateway response rules by [clicking here](#response-logic).

#### endpoint.abort-if-status-codes

Optional field, of integer list type, the default value is null, indicating that any backend executed on the endpoint
that responded with HTTP status code greater than or equal to `400 (Bad request)` will be aborted.

If informed, and a backend returns the HTTP status code indicated in the configuration, the endpoint will be aborted.
This means that the other configured backends, set after it, will not be executed, and the endpoint will return its
response to the end client.

If you don't want any HTTP status code to be aborted on the endpoint, simply enter the field empty.

See how the endpoint will be responded to after a backend has been aborted by [clicking here](#aborted-response).

#### endpoint.beforewares

Optional field, of string list type, the default value is empty, indicating that the endpoint has no pre-request
middlewares.

If informed, the endpoint will run requests, position by position, starting at the beginning of the list. If the string
value of the position to be executed is correctly configured in the [middlewares](#middlewares) field, the backend
configured in it will be executed. Otherwise, it will only ignore the position and print a warning log.

#### endpoint.afterwares

Optional field, of string list type, the default value is empty, indicating that the endpoint has no post-request
middlewares.

If informed, the endpoint will run requests, position by position, starting at the beginning of the list. If the string
value of the position to be executed is correctly configured in the [middlewares](#middlewares) field, the backend
configured in it will be executed. Otherwise, it will ignore the position and just print a warning log.

#### endpoint.backends

Required field, of object list type, responsible for the main execution of the endpoint, the name itself says it all,
it's a list that indicates all the necessary services for the endpoint to return the expected response.

See how the basic flow of a backend works in the image below:

<img src="assets/basic-backend-flowchart-en.png" alt="basic backend flow chart">

Below we will list and explain each field of this important object:

### backend.@comment

Optional field, of string type, free field for notes related to your backend.

### backend.name

Optional field, of string type, is responsible for the name of your backend service, it is used to name the field of its
aggregate response, if the [backend.extra-config.group-response](#backendextra-configgroup-response) field is set
to `true`.

### backend.hosts

Required field, of string list type, is responsible for the hosts of your service that the API Gateway will call
together with the [backend.path](#backendpath) field.

In a way, we can have a "dumb" load balancer, because the backend will randomly choose which host will be called from
this list, so we can inform the same host multiple times to balance the calls, see:

````
50% each
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

Required field, of string type, the value indicates the URL of the backend service path.

We use one of the informed [backend.hosts](#backendhosts) and join it with the provided path, for example, in the hosts
field we have the value:

```json
[
  "https://instance-01",
  "https://instance-02"
]
```

And in this path field we have the value:

```text
/users/status
```

The backend will construct the following request URL:

```text
https://instance-02/users/status
```

See how the host is balanced [by clicking here](#backendhosts).

### backend.method

Required field, of string type, the value indicates which HTTP method the backend service expects.

### backend.forward-queries

Optional field, of string list type, the default value is empty, indicating that any query parameter will be passed to
the backend service.

If informed, only the indicated fields will be passed to the backend service, for example, we receive a request with the
following URL:

```text
/users?id=23&email=gabrielcataldo@gmail.com&phone=47991271234
```

In this example we have the `forward-queries` field with the following values:

```json
[
  "email",
  "phone"
]
```

The backend request URL was:

```text
/users?email=gabrielcataldo@gmail.com&phone=47991271234
```

We saw that the `id` query parameter was not passed to the backend service, as it was not mentioned in the list.

### backend.forward-headers

Optional field, of string list type, the default value is empty, indicating that any received header will be passed to
the backend service.

If informed, only the indicated fields will be passed to the backend service, for example, we receive a request with the
following header:

```json
{
  "Device": "95D4AF55-733D-46D7-86B9-7EF7D6634EBC",
  "User-Agent": "IOS",
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

In this example we have the `forward-headers` field with the following values:

```json
[
  "User-Agent",
  "Authorization"
]
```

The backend request header was:

```json
{
  "User-Agent": "IOS",
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

We saw that the `Device` field of the received header was not passed to the backend service because it was not mentioned
in the list.

### backend.extra-config

Optional field, of object type, indicates additional backend service configuration, see below about the fields and their
responsibilities:

#### backend.extra-config.omit-request-body

Optional field, of boolean type, the default value is `false`, indicating that the request body will be passed to the
backend if it has one.

If `true` is informed the request body will not be passed to the backend.

#### backend.extra-config.group-response

Optional field, of boolean type, the default value is `false`, indicating that the backend response body does not need
to be grouped into a JSON field for the final client response.

If informed with the value `true` the response body, if any, will be grouped into a JSON field of the final response,
the name of the field will be the [name](#backendname) of the backend service if filled, if not we have a naming
standard that is `backend-position in the list` which would, for example, be `backend-0`.

To understand the importance of this field, see more about the [response logic](#response-logic) of the API Gateway.

#### backend.extra-config.omit-response

Optional field, of boolean type, the default value is `false`, indicating that the response from the backend in question
will not be omitted to the final client.

If informed with the value `true`, the entire response from the backend in question will be omitted, be careful, because
if you only have this backend, and it is omitted, the API Gateway will respond by default HTTP status
code `204 (No Content)`.

To understand the importance of this field, see more about the [response logic](#response-logic) of the API Gateway.

### backend.modifiers

Optional field, of object type, the default value is empty, indicating that there is no modification process in this
backend in question.

Below we will list and explain each field of this so important object:

### modifiers.@comment

Optional field, of string type, free field for notes related to your modifiers.

### modifiers.status-code

Optional field, of integer type, the default value is `0`, indicating that there is nothing to be modified in the HTTP
response status code from the backend.

If informed, the HTTP response status code from the backend will be modified by the inserted value, this may or may not
influence the final response of the endpoint, see the [response-logic](#response-logic) of the API Gateway to learn
more.

### modifiers.header

Optional field, of list of object type, the default value is empty, responsible for the modifications of the request and
response header of the backend.

See below the fields of this object and their responsibility:

#### header.@comment

Optional field, of string type, free field for notes related to your modifier.

#### header.context

Mandatory field, of string type, is responsible for indicating which context the modification should act on.

Accepted values:

`REQUEST`: To act in the pre-request of the backend.

`RESPONSE`: To act post-request of the backend.

It's important to remember that if the value is `REQUEST` you may use in the [header.scope](#headerscope) field only the
value `REQUEST`.

#### header.scope

Optional field, of string type, the default value will be based on the informed [header.context](#headercontext) field,
the value indicates which scope we should change, whether the request scope or the response scope.

Accepted values:

`REQUEST`: To modify the request scope, this type of scope can have global action propagating this change in the
following backend requests, just use the [header.propagate](#headerpropagate) field as `true`.

`RESPONSE`: To change the backend response scope.

#### header.action

Mandatory field, of string type, responsible for the action to be taken in the header modification.

Accepted values:

`ADD`: Adds the key informed in the [header.key](#headerkey) field if it does not exist, and adds the value informed in
the [header.value](#headervalue) field.

`APD`: Adds the value informed in the [header.value](#headervalue) field if the key informed in
the [header.key](#headerkey) field exists.

`SET`: Sets the value of the key informed in the [header.key](#headerkey) field by the value passed in
the [header.value](#headervalue) field.

`RPL`: Replaces the value of the key informed in the [header.key](#headerkey) field by the value passed in
the [header.value](#headervalue) field if it exists.

`REN`: Renames the key informed in the [header.key](#headerkey) field by the value passed in
the [header.value](#headervalue) field if it exists.

`DEL`: Removes the key informed in the [header.key](#headerkey) field if it exists.

#### header.key

Mandatory field, of string type, used to indicate which header key should be modified.

#### header.value

Mandatory field, of string type, used as the value to be used to modify the key indicated in
the [header.key](#headerkey) field.

We have possibilities of using [dynamic values](#dynamic-values-for-modification), and
of [environment variables](#environment-variables) for this field.

NOTE: it only becomes optional if [query.action](#queryaction) has the `DEL` value.

#### header.propagate

Optional field, of boolean type, the default value is `false` indicating that the modifier should not propagate this
change to the following backends.

If informed as `true` this modification will be propagated to the following backends.

IMPORTANT: This field is only accepted if the [scope](#headerscope) has the `REQUEST` value.

### modifiers.param

Optional field, of object list type, default value is empty, responsible for modifying request parameters for the
backend.

See below the fields of this object and their responsibilities:

#### param.@comment

Optional field, of string type, free field for annotations related to your modifier.

#### param.context

Mandatory field, of string type, is responsible for indicating which context the modification should act.

Accepted values:

`REQUEST`: To act on the backend's pre-request.

`RESPONSE`: To act post-backend request.

#### param.action

Mandatory field, of string type, responsible for taking action in modifying the request parameters.

Accepted values:

`SET`: Sets the value of the key informed in the [param.key](#paramkey) field by the value passed in
the [param.value](#paramvalue) field.

`RPL`: Replaces the value of the key informed in the [header.key](#headerkey) field by the value passed in
the [header.value](#headervalue) field if it exists.

`REN`: Renames the key informed in the [param.key](#paramkey) field by the value passed in
the [param.value](#paramvalue) field if it exists.

`DEL`: Removes the key informed in the [param.key](#paramkey) field if it exists.

#### param.key

Mandatory field, of string type, used to indicate which parameter key should be modified.

#### param.value

Mandatory field, of string type, used as the value to modify the key indicated in the [param.key](#paramkey) field.

We have possibilities for using [dynamic values](#dynamic-values-for-modification),
and [environment variables](#environment-variables) for this field.

NOTE: it only becomes optional if [query.action](#queryaction) has the `DEL` value.

#### param.propagate

Optional field, of boolean type, the default value is `false` indicating that the modifier should not propagate this
change to the following backends.

If informed as `true` this modification will be propagated to the following backends.

### modifiers.query

Optional field, of object list type, default value is empty, responsible for modifying the query parameters of the
request for the backend.

See below the fields of this object and their responsibilities:

#### query.@comment

Optional field, of string type, free field for annotations related to your modifier.

#### query.context

Mandatory field, of string type, responsible for indicating in which context the modification should take place.

Accepted values:

`REQUEST`: To act in the backend's pre-request.

`RESPONSE`: To act post-backend request.

#### query.action

Mandatory field, of string type, responsible for the action to be taken in modifying the query parameters of the
request.

Accepted values:

`ADD`: It adds the key informed in the [query.key](#querykey) field if it does not exist, and the value reported in
the [query.value](#queryvalue) field is added.

`APD`: It adds the value informed in the [query.value](#queryvalue) field if the key informed in
the [query.key](#querykey) field exists.

`SET`: It sets the value of the key informed in the [query.key](#querykey) field by the value passed in
the [query.value](#queryvalue) field.

`RPL`: Replaces the value of the key informed in the [query.key](#querykey) field by the value passed in
the [query.value](#queryvalue) field if it exists.

`REN`: Renames the key informed in the [query.key](#querykey) field by the value passed in
the [query.value](#queryvalue) field if it exists.

`DEL`: It removes the key informed in the [query.key](#querykey) field if it exists.

#### query.key

Mandatory field, of string type, used to indicate which query parameter key should be modified.

#### query.value

Mandatory field, of string type, used as the value to be used to modify the key indicated in the [query.key](#querykey)
field.

We have possibilities of using [dynamic values](#dynamic-values-for-modification),
and [environment variables](#environment-variables) for this field.

NOTE: it becomes optional only if [query.action](#queryaction) has the value `DEL`.

#### query.propagate

Optional field, of boolean type, the default value is `false` indicating that the modifier should not propagate this
change in question to the following backends.

If informed as `true` this modification will be propagated to the following backends.

### modifiers.body

Optional field, of list of object type, default value is empty, responsible for modifying the body of request or
response from backend.

See below the fields of this object and their responsibility:

#### body.@comment

Optional field, of string type, free field for annotations related to your modifier.

#### body.context

Mandatory field, of string type, responsible for indicating in which context the modification should take place.

Accepted values:

`REQUEST`: To act on the backend's pre-request.

`RESPONSE`: To act post-backend request.

It's important to remember that if the value is `REQUEST`, you can only use the value `REQUEST` in
the [body.scope](#bodyscope) field.

#### body.scope

Optional field, of string type, the default value will be based on the [body.context](#bodycontext) field informed, the
value indicates which scope we should change, whether the request scope or response scope.

Accepted values:

`REQUEST`: To modify the request scope, this type of scope can have a global action propagating this change in the
following backend requests, just use the [body.propagate](#bodypropagate) field as `true`.

`RESPONSE`: To modify the backend response scope.

#### body.action

Mandatory field, of string type, responsible for the action to be taken in the body modification.

Accepted values if body is JSON:

`ADD`: Adds the key informed in the [body.key](#bodykey) field if it does not exist, and aggregates the value informed
in the [body.value](#bodyvalue) field.

`APD`: Adds the value informed in the [body.value](#bodyvalue) field if the key informed in the [body.key](#bodykey)
field exists.

`SET`: Sets the value of the key informed in the [body.key](#bodykey) field by the value passed in
the [body.value](#bodyvalue) field.

`RPL`: Replaces the value of the key informed in the [body.key](#bodykey) field by the value passed in
the [body.value](#bodyvalue) field if it exists.

`REN`: Renames the key informed in the [body.key](#bodykey) field to the value passed in the [body.value](#bodyvalue)
field if it exists.

`DEL`: Removes the key informed in the [body.key](#bodykey) field if it exists.

Accepted values if body is TEXT:

`ADD`: Adds the value informed in the [body.value](#bodyvalue) field to text.

`APD`: Adds the value informed in the [body.value](#bodyvalue) field if body is not empty.

`RPL`: Will replace all values similar to the key informed in the [body.key](#bodykey) field with the value provided in
the [body.value](#bodyvalue) field.

`DEL`: Removes all values similar to the key informed in the [body.key](#bodykey) field.

#### body.key

Mandatory field, of string type, used to indicate which header key should be modified.

NOTE: it becomes optional if your body is of TEXT type and [body.action](#bodyaction) has the value `ADD`.

#### body.value

Mandatory field, of string type, used as value to be used to modify the key indicated in the [body.key](#bodykey) field.

We have possibilities of using [dynamic values](#dynamic-values-for-modification),
and [environment variables](#environment-variables) for this field.

NOTE: It becomes optional only if [body.action](#bodyaction) has the value `DEL`.

#### body.propagate

Optional field, of boolean type, the default value is `false`, indicating that the modifier should not propagate this
change to subsequent backends.

If informed as `true`, this modification will be propagated to the following backends.

IMPORTANT: This field is only accepted if the [scope](#bodyscope) has the value `REQUEST`.

## Runtime JSON

When started, the GOPEN API Gateway generates a JSON file, based on the [configuration JSON](#configuration-json),
located in the `runtime` folder at the root of your workspace.

This JSON indicates the application's understanding after reading the [configuration JSON](#configuration-json), all
[#configuration variables](#environment-variables) will already have their values replaced, if any.

This json can also be read using the static route [/settings](#settings).

## Static routes

The GOPEN API Gateway has some static endpoints, that is, regardless of any configuration made, there are currently
three endpoints registered in its routes, see the details for each one and their responsibilities below:

### `/ping`

Endpoint to know if the API Gateway is alive, returns `404 (Not Found)` if it's off, and `200 (OK)` if it's running.

### `/version`

Endpoint that returns the version obtained from the [version](#version) configuration, returns `404 (Not Found)` if it
has not been informed in the [configuration json](#configuration-json), otherwise returns `200 (OK)` with the value in
the body as text.

### `/settings`

Endpoint returns some information about the project, such as version, version date, number of contributors, and a
summary of how many endpoints, middlewares, backends, and modifiers are currently configured and the configuration json
that is actively running.

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

## Environment variables

Environment variables can be easily instantiated using the .env file, in the folder indicated by the dynamic startup
environment as mentioned in the topic [FOLDER STRUCTURE](#folder-structure).

If you prefer to insert the values using docker-compose it will also work correctly, the point is that the API Gateway
will read the value recorded on the machine, regardless of how it was entered there.

The values can be used in the API Gateway JSON configuration, just use the syntax `$NAME` as a string value, see in the
example below.

A snippet of a configuration JSON, we have the following values:

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

And on our machine we have the following environment variables:

```dotenv
VERSION=1.0.0

REDIS_URL=redis-18885.c259.us-east-1-4.ec2.cloud.redislabs.com:18985
REDIS_PASSWORD=12345

TIMEOUT=5m
```

The API Gateway generates a [runtime JSON](#runtime-json) file when running the application, see the result of the same
after starting the application:

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

We see that all values with the syntax `$NAME` were replaced with their due values, if a value has been mentioned by
this syntax, but it does not exist in the environment variables, the same informed value will be kept.

### Note

In this first version, we can use this syntax only in string type fields, if you want
to [contribute](#how-to-contribute) with a solution that will accept this syntax regardless of the type, we would be
grateful!

## Dynamic values for modification

We can use request and response values from the endpoint runtime, as it was configured. These values can be obtained by
a specific syntax, we have the following possibilities to obtain these values, see:

### Request

When mentioning the syntax `#request...` you will be getting the values from the received request, if you mention the
field `#request.history...` you will be getting the values from
the [temporary history of requests](#temporary-history-of-requests).

`#request.header...`

This snippet of the syntax will get the value indicated from the request header, for example,
`#request.header.X-Forwarded-For.0` will get the first value of the field `X-Forwarded-For` from the request header if
it exists, replacing the syntax with the value, the result was `127.0.0.1`.

`#request.params...`

This snippet of the syntax will get the indicated value from the request parameters, for example,
`#request.params.id` will get the value of the `id` field from the request parameters if it exists, replacing the syntax
by the value, the result was `72761`.

`#request.query...`

This snippet of the syntax will get the indicated value from the request search parameters, for example,
`#request.query.email.0` will get the first value of the `email` field from the request search parameters if it exists,
replacing the syntax with the value, the result was `gabrielcataldo.adm@gmail.com`.

`#request.body...`

This snippet of the syntax will get the indicated value from the request body, for example,
`#request.body.deviceId` will get the value of the `deviceId` field from the request body if it exists, replacing the
syntax with the value, the result was `991238`.

#### Temporary History of Requests

If you want to obtain request values specifically from a backend, just add the `history` field in the syntax, and the
position of your backend in history, remembering that the middlewares `beforewares` and `afterwares` count as a backend
in history, another point is that you will only have the request information in this field if the backend has already
been processed, or the backend is in process. Here is what an example would look like:

In the example, I have only one backend, and it is being processed, so I can be using the syntax:

`#request.history.0.params.phone`

In this other syntax example we have three configured backends and two have already been processed, so we can use the
following syntax in the process of the third backend:

`#request.history.1.body.id`

In these examples mentioned we see that we can get the value of the request that is being made or has already been made
by a configured backend, and that they are stored in a type of temporary history.

### Response

When mentioning the syntax `#response...` you will be retrieving the values of the API Gateway response prepared based
on the [response logic](#response-logic), if you mention the field `#response.history...` you will be retrieving the
values from the [temporary history of responses](#temporary-history-of-responses).

`#response.statusCode`

This syntax will get the HTTP status code of the current endpoint response. If you want to get from the response of a
specific backend use the [temporary history of responses](#temporary-history-of-responses).

`#response.header...`

This snippet of the syntax will get the indicated value from the endpoint response header, for example,
`#response.header.X-Trace-Id` will get the list values of the field `X-Trace-Id` from the response header if it exists,
replacing the syntax with the value, the result was `["4ae6c92d16089e521626", "5ae7c93d16089e521628"]`.

`#response.body...`

This snippet of the syntax will get the indicated value from the endpoint response body, for example,
`#response.body.role` will get the value of the field `role` from the request body if it exists, replacing the syntax
with the value, the result was `ADMIN`.

#### Temporary History of Responses

If you want to retrieve the response values specifically from a backend, just add the `history` field in the syntax, and
the position of your backend in history, remembering that the middlewares `beforewares` and `afterwares` count as a
backend in history, another point is that you will only have the response information in this field if the backend has
already been processed. Here is what an example would look like:

In the example, I have only one backend, and it has been processed, so I can be using the syntax:

`#response.history.0.header.X-Value.0`

In this other syntax example, we have three configured backends and two have already been processed, so we can use the
following syntax in the process of the third backend:

`#request.history.1.body.users.0`

In these examples mentioned we see that we can get the value of a response from a backend that has already been
processed, and that they are stored in a kind of temporary history.

### Important

You can use based on these fields,
the [JSON path syntax](https://github.com/tidwall/gjson/blob/master/README.md#path-syntax) that fits your values, just
remember that header, query are maps of list of string, and params is a map of string.

### See also

- Backend modifiers configuration by clicking [here](#backendmodifiers).

- Learn in practice how to use dynamic values for modification using
  the [playground](https://github.com/GabrielHCataldo/gopen-gateway-playground) project that already comes with some
  examples of modifiers with dynamic values.

## Observability

The API Gateway, upon receiving a request, will add the client's IP to the `X-Forwarded-For` header of the underlying
requests, and also adds, if it does not exist, a unique generated value to the header called `X-Trace-Id` for better
observability of the received request in both the API Gateway logs and the underlying microservices.

## Response logic

When using an API Gateway, we often wonder, how will the response of this configured endpoint be returned to my client?

To facilitate understanding, we have created this topic to summarize the response logic of our API Gateway, so let's get
started.

### How does it work?

The API Gateway was developed with intelligence and flexibility in responding to an endpoint, it is based on two
important points, first, on the number of responses from backend services that have been processed, and second, on the
customization fields of the response configured in the [endpoint](#endpointcomment) and [backend](#backendname) objects.
Let's see some examples below for better understanding.

#### Single Backend

In this example, we will work with only one backend. See how the API Gateway will respond to this scenario.

Configuration JSON

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

While processing this endpoint, the API Gateway response was:

```text
HTTP/1.1 200 OK
```

Header ([Read more about response headers here](#response-header))

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

We saw in the example that the API Gateway served as a proxy by redirecting the request to the configured backend
service and mirroring its response body, and adding its values to the response header.

In the same example, let's force an unfortunate scenario in the backend response, see:

```text
HTTP/1.1 404 Not Found
```

Header ([Read more about response headers here](#response-header))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: false
Date: Tue, 23 Apr 2024 21:56:33 GMT
Content-Length: 235
```

Body

```json
{
  "file": "datastore/user.go",
  "line": 227,
  "endpoint": "/users/find/gabrielcataldo.adma@gmail.com",
  "message": "user not found"
}
```

In this case, the API Gateway also mirrored the response from the single backend call of the endpoint.

#### Using middlewares

In this example, we will use the [beforewares](#endpointbeforewares) and [afterwares](#endpointafterwares) middlewares,
since these backends are omitted to the final client if they are successful, we will simulate a call with the blocked
device so that the [beforeware](#endpointbeforewares) returns an error, and then an [afterware](#endpointafterwares)
that will also respond an error, because it does not exist, let's go!

Configuration JSON

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

While processing this sample endpoint simulating an error in the [beforeware](#endpointbeforewares) call, the API
Gateway response was:

```text
HTTP/1.1 403 Forbidden
```

Header ([Read more about response headers here](#response-header))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: true
Date: Tue, 23 Apr 2024 23:02:09 GMT
Content-Length: 154
```

Body

```json
{
  "file": "service/device.go",
  "line": 49,
  "endpoint": "/devices",
  "message": "unprocessed entity: device already exists and is not active"
}
```

We saw that the response was a mirror of the beforeware `save-device` return, because, as it returned
failure `403 (Forbidden)`, the endpoint aborted, not calling the following backends. Remember that you
can configure the HTTP status codes that will be aborted by your endpoint, just fill in the
field [endpoint.abort-if-status-codes](#endpointabort-if-status-codes).

In the following example we will force an error in the afterware `increment-attempts` and the API Gateway response was:

```text
HTTP/1.1 404 Not Found
```

Header ([Read more about response headers here](#response-header))

```text
Content-Type: text/plain
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: false
Date: Tue, 23 Apr 2024 23:16:57 GMT
Content-Length: 18
```

Body

```text
404 page not found
```

We noticed that the response also mirrored the return of the `increment-attempts` afterware, even though it is the last
call of a backend service from the endpoint, as it fell under the aborted response rule, so all other responses from
other backends were ignored and only the response from the aborted backend was returned.

Learn more about the [aborted response](#aborted-response).

#### Múltiplos backends

In this example, we will be working with three [backends](#endpointbackends) on the endpoint, so let's get started!

Configuration JSON:

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

In the example, we will execute all three backends successfully, and the API Gateway responded:

```text
HTTP/1.1 200 OK
```

Header ([Read more about response headers here](#response-header))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Tue, 23 Apr 2024 23:49:12 GMT
Content-Length: 755
```

Body

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

There are several points worth noting in this example. Firstly, regarding the format, the API Gateway understood that
your endpoint has multiple responses, and the [endpoint.aggregate-responses](#endpointaggregate-responses) field with
the value `true` was not used. Therefore, it lists the responses as JSON, adding the following fields:

`ok`: Indicates whether the response from the respective backend had an HTTP status code between `200` and `299`.

`code`: Filled with the HTTP status code responded by your backend.

These fields are only added if there are multiple responses, and
the [endpoint.aggregate-responses](#endpointaggregate-responses) field is not specified with the value `true`.

The second point to highlight is in the `"version": "v1.0.0"` section of the last backend. It responded with only a text
in the response body, which was `v1.0.0`. However, in this scenario, as mentioned, the API Gateway forces the parsing of
this value into JSON, adding a new field with the name obtained from the [backend.name](#backendname) configuration and
with its value.

The third point concerns the HTTP status code. It is returned based on the highest frequency. That is, if we have
three `200 OK` responses as in the example, the API Gateway will also return this code. If we have an equal frequency of
different codes, the last HTTP status code returned will be considered. Below are the possible scenarios of this logic:

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

The API Gateway will respond with `201 Created`.

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

The API Gateway will respond with `100 Continue`.

The fourth point to highlight is that since the endpoint has multiple responses, consequently, it has multiple response
headers. The API Gateway will aggregate all the fields and values for the final response header. Learn more about the
behavior of the response header by clicking [here](#response-header).

The last point to emphasize is that if one of these backend returns enters the scenario where the endpoint aborts the
response, it will not follow any guidelines shown in the respective topic but rather
the [aborted response logic](#aborted-response).

#### Multiple Aggregated Backends

In this example, we will use a configuration similar to the JSON configuration example above, but with
the [endpoint](#endpointaggregate-responses) field set to `true`.

Configuration JSON:

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

When processing the endpoint, the API Gateway responded with:

```text
HTTP/1.1 200 OK
```

Header ([Read more about response headers here](#response-header))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Wed, 24 Apr 2024 10:57:31 GMT
Content-Length: 665
```

Body

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

Here's the translation using Markdown code snippet:

We noticed the only difference in response from the [Multiple Aggregated Backends](#multiple-aggregated-backends) topic
is that it aggregated the values of all responses into a single JSON, and for repeated fields, the values were
aggregated into a list.

Other rules such as HTTP status code, forced conversion to JSON, among others, follow the same rules mentioned in
the [Multiple Aggregated Backends](#multiple-aggregated-backends) topic.

In this example, we can make the aggregated response a bit more organized. To achieve this, let's modify the snippet of
our second backend by adding the [backend.extra-config.group-response](#backendextra-configgroup-response) field with
the value `true` and giving it a name. Below is the modified snippet of the JSON configuration:

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

Upon processing the endpoint again, we obtained the following response:

```text
HTTP/1.1 200 OK
```

Header ([Read more about response headers here](#response-header))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Wed, 24 Apr 2024 11:23:07 GMT
Content-Length: 697
```

Body

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

With this configuration, we noticed that our aggregated response became more organized. It's important to understand
the [configuration JSON](#configuration-json) and its fields so that the GOPEN API Gateway can better meet your needs.

### Aborted response

For a response to be aborted by the API Gateway, one of the configured backends of the endpoint, either through
middlewares or main backends, when processed, in its response, the HTTP status code needs to match values in
the [endpoint.abort-if-status-codes](#endpointabort-if-status-codes) field of the endpoint itself.

It's important to understand these two points:

- When aborting the backend response, the API Gateway will only mirror the response of the backend itself, including the
  status code, headers, and body; thus, other processed responses will be ignored.
- Even if the backend is aborted, the [modifiers](#backendmodifiers) will still be processed.

We recommend using this configuration only when something unexpected occurs, such as an error response
like `500 (Internal Server Error)`.

### Response header

In the response, the API Gateway aggregates all header values responded by the configured backends in the endpoint,
regardless of the number of backends, which also includes the [middlewares](#middlewares), except for
the `Content-Length`, `Content-Type`, and `Date` fields.

#### Default Header Fields

Additionally, up to four fields are added to the header. See below for details:

`X-Gopen-Cache`: If the endpoint response is not "fresh", meaning the cached response was used, it returns `true`;
otherwise, it returns `false`.

`X-Gopen-Cache-Ttl`: If the endpoint response was made using caching, it returns the remaining time-to-live (TTL) of
that cache; otherwise, the field is not returned.

`X-Gopen-Complete`: If all backends have been processed by the endpoint, it returns `true`; otherwise, it
returns `false`.

`X-Gopen-Success`: If all backends have returned success, meaning the HTTP status code of the response is between `200`
and `299`, it returns `true`; otherwise, it returns `false`.

It's worth noting that if the response from a backend is [aborted](#aborted-response), only the header of that backend
is aggregated, and the rules of the above fields are considered.

Now let's see some examples of response headers:

#### Unique Header Fields

Backend 1 response header:

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521626
X-MS: api-user
Date: Wed, 24 Apr 2024 11:23:07 GMT
Content-Length: 102
```

Endpoint response header:

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

We saw that in the example, the [default fields](#default-header-fields) were added, and the values of the response
header were aggregated, which were `X-Value-Id` and `X-MS`.

#### Duplicate Header Fields

Backend 1 response header:

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521626
X-MS: api-user
Date: Wed, 24 Apr 2024 11:23:07 GMT
Content-Length: 102
```

Backend 2 response header:

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521638
X-MS: api-device
X-MS-Success: true
Date: Wed, 24 Apr 2024 11:23:08 GMT
Content-Length: 402
```

Endpoint response header:

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

We also noticed that in the example, the [default fields](#default-header-fields) were added, and the values of the
response header were aggregated, which were `X-Value-Id`, `X-MS`, and `X-MS-Success`. It's worth noting that the
duplicated fields were grouped and separated by commas.

### Default responses

Every API Gateway has its default responses for each error scenario, so we'll list each scenario below along with its
respective HTTP response:

#### 413 (Request Entity Too Large)

This scenario occurs when the request body size exceeds the limit allowed for the endpoint, using the
configuration [limiter.max-body-size](#limitermax-header-size) for normal body
and [limiter.max-multipart-memory-size](#limitermax-multipart-memory-size) for `form-data` type submission. You can
customize this configuration for a specific endpoint using the [endpoint.limiter](#endpointlimiter) field.

Header

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 11:56:06 GMT
Content-Length: 170
```

Body

```json
{
  "file": "infra/size_limiter.go",
  "line": 92,
  "endpoint": "/users",
  "message": "payload too large error: permitted limit is 1.0B",
  "timestamp": "2024-04-26T08:56:06.628636-03:00"
}
```

#### 429 (Too Many Requests)

This scenario occurs when the request limit is reached for a particular IP address. This limit is defined in
the [limiter.rate](#limiterrate) configuration. You can customize this configuration for a specific endpoint using
the [endpoint.limiter](#endpointlimiter) field.

Header

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 12:12:53 GMT
Content-Length: 177
```

Body

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

This scenario occurs when the size of the header exceeds the limit allowed for the endpoint, using the
configuration [limiter.max-header-size](#limitermax-header-size). You can customize this configuration for a specific
endpoint using the [endpoint.limiter](#endpointlimiter) field.

Header

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 11:39:53 GMT
Content-Length: 186
```

Body

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

This scenario is specific when an unexpected error occurs with the API Gateway. If this happens, report the
issue [here](https://github.com/GabrielHCataldo/gopen-gateway/issues) showing the response and the log printed on the
execution terminal.

Header

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 12:38:16 GMT
Content-Length: 183
```

Body

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

This scenario occurs when trying to communicate with the backend, and there is a communication failure with it.

Header

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Thu, 25 Apr 2024 01:07:36 GMT
Content-Length: 277
```

Body

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

This scenario occurs when the endpoint exceeds the time limit configured in the [timeout](#timeout) field. You can
customize this configuration for a specific endpoint using the [endpoint.timeout](#endpointtimeout) field.

Header

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 13:29:55 GMT
Content-Length: 150
```

Body

```json
{
  "file": "middleware/timeout.go",
  "line": 81,
  "endpoint": "/users/version",
  "message": "gateway timeout: 5m",
  "timestamp": "2024-04-26T10:29:55.908526-03:00"
}
```

# Usability

Use the [playground project](https://github.com/GabrielHCataldo/gopen-gateway-playground) to start exploring and using
the GOPEN API Gateway in practice!

# How to contribute?

We are delighted to see the community come together, and projects like this one are open to receiving your ideas. See
below how to participate.

## Download

To run the project, first download the [Go language](https://go.dev/dl/) version 1.22 or higher on your machine.

With Go installed on your machine, pull the project:

```text
git pull https://github.com/GabrielHCataldo/gopen-gateway.git
```

Then open it using the terminal and your preferred IDE.

Goland:

```text
goland gopen-gateway
```

VSCode:

```text
code gopen-gateway
```

## Air

If you want to use hot reload to optimize development, the project already has the `.air.toml` file as an example for
you to use as you wish.

See the documentation for this fantastic project by [clicking here](https://github.com/cosmtrek/air).

## Gitflow

To initialize development, you can create a branch from the main branch for a future PR to it.

# Thanks

This project has been supported by fantastic libraries, and this section is dedicated to each one listed below:

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

Thank you for contributing to the Go community and facilitating the development of this project.

# Apache 2.0 License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway?ref=badge_large&issueType=license)
