<img src="assets/logo.png" alt="">

[![Project status](https://img.shields.io/badge/version-v1.0.0-orange.svg)](https://github.com/tech4works/gopen-gateway/releases/tag/v1.0.0)
[![Playground](https://img.shields.io/badge/%F0%9F%8F%90-playground-9900cc.svg)](https://github.com/tech4works/gopen-gateway)
[![Docker](https://badgen.net/badge/icon/docker?icon=docker&label)](https://hub.docker.com/r/gabrielhcataldo/gopen-gateway)
[![Open Source Helpers](https://www.codetriage.com/gabrielhcataldo/gopen-gateway/badges/users.svg)](https://www.codetriage.com/gabrielhcataldo/gopen-gateway)
[![TODOs](https://img.shields.io/badge/TODOs-13-red.svg?style=flat)](https://github.com/tech4works/gopen-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/tech4works/gopen-gateway)](https://goreportcard.com/report/github.com/tech4works/gopen-gateway)
[![GoDoc](https://godoc.org/github/GabrielHCataldo/gopen-gateway?status.svg)](https://pkg.go.dev/github.com/tech4works/gopen-gateway/helper)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/GabrielHCataldo/gopen-gateway)](https://github.com/tech4works/gopen-gateway/blob/main/go.mod)
[![HitCount](https://hits.dwyl.com/GabrielHCataldo/gopen-gateway.svg)](http://hits.dwyl.com/GabrielHCataldo/gopen-gateway)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway.svg?type=small)](https://app.fossa.com/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway?ref=badge_small)

![United States](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/us.png)
[Inglés](https://github.com/tech4works/gopen-gateway/blob/main/README.md) |
![Brazil](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/br.png)
[Portugués](https://github.com/tech4works/gopen-gateway/blob/main/README.pt-br.md)

El proyecto GOPEN fue creado con el objetivo de ayudar a los desarrolladores a tener una API Gateway robusta y fácil de
manejar, con la oportunidad de trabajar en mejoras agregando a la comunidad, y lo más importante, sin coste alguno. Fue
desarrollada porque muchas APIs Gateway gratuitas en el mercado no satisfacen muchas necesidades mínimas
para una aplicación, lo que induce a hacer la actualización.

Con esta nueva API Gateway no tendrás que equilibrar platos para ahorrar en tu infraestructura y arquitectura,
y aún optimizará tu desarrollo, mira abajo todos los recursos disponibles:

- JSON de configuración simplificada para múltiples entornos.
- Configuración rápida de variables de entorno para múltiples entornos.
- Versionado a través de JSON de configuración.
- Ejecución vía docker con recarga en caliente opcional.
- Configuración de timeout global y local para cada endpoint.
- Configuración de caché global y local para cada endpoint, con personalización de la estrategia de almacenamiento de
  llave y condiciones basadas en los códigos de estado de respuesta y el método HTTP para leer y guardar el mismo.
- Almacenamiento de caché local o global utilizando Redis.
- Configuración de limitador de tamaño, global y local para cada endpoint, limitando el tamaño del Header, Body y
  Multipart Memory.
- Configuración de limitador de tasa, global y local para cada endpoint, limitando por tiempo y ráfaga por IP.
- Configuración de seguridad de CORS con validaciones de orígenes, método HTTP y headers.
- Configuración global de múltiples middlewares, para ser utilizados posteriormente en el endpoint si se indica.
- Filtrado personalizado de envío de headers y query a los backends del endpoint.
- Procesamiento de múltiples backends, siendo ellos antes de la ejecución, principales y después de la ejecución para el
  endpoint.
- Configuración personalizada para abortar el proceso de ejecución de los backends por el código de estado devuelto.
- Modificadores para todos los contenidos de la petición y respuesta (Código de Estado, Ruta, Header, Parámetros,
  Consulta, Cuerpo) a nivel global (petición/respuesta) y local (petición backend respuesta backend) con acciones de
  eliminación, añadido, cambio, sustitución y renombrado.
- Obtenga el valor a modificar de las variables de entorno, de la petición actual, del historial de respuestas del
  endpoint, o incluso del valor pasado en la configuración.
- Ejecute los modificadores en el contexto que desee, antes de una petición backend o después, usted decide.
- Haga que las modificaciones se reflejen en todas las peticiones/respuestas siguientes, usando la misma a nivel global.
- Omita la respuesta de un backend si es necesario, la misma no será utilizada en la respuesta del endpoint.
- Omita el cuerpo de la petición de su backend si es necesario.
- Agregue sus múltiples respuestas de los backends si lo desea, pudiendo personalizar el nombre del campo a ser asignado
  a la respuesta del backend.
- Agrupe el cuerpo de la respuesta de su backend en un campo específico de la respuesta del endpoint.
- Personalización del tipo de respuesta del endpoint pudiendo ser JSON, TEXT y XML.
- Obtenga más observabilidad con el registro automático del trace id en el encabezado de las peticiones siguientes y
  registros bien estructurados.

# Documentación

Para entender cómo funciona, necesitamos explicar primero la estructura de los entornos dinámicos que GOPEN acepta para
su configuración en JSON y archivo de variables de entorno, ¡así que vamos allá!

## Estructura de carpetas

En la estructura del proyecto, en su raíz necesitará tener una carpeta llamada "gopen" y dentro de ella necesita tener
las carpetas conteniendo los nombres de sus entornos, puede dar el nombre que quiera, esta carpeta necesitará tener al
menos el archivo ".json" de configuración de la API Gateway, quedará más o menos así, por ejemplo:

### Proyecto GO

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

### Proyecto usando imagen docker

    nombre-de-su-proyecto
    | - docker-compose.yml
    | - gopen
      | - dev
      |   - .json
      |   - .env
      | - prd
      |   - .json
      |   - .env

## JSON de configuración

Con base en este archivo JSON de configuración obtenido por el entorno deseado la aplicación tendrá sus endpoints y sus
reglas definidas, vea abajo un ejemplo simple con todos los campos posibles y sus conceptos y reglas:

````json
{
  "$schema": "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json",
  "@comment": "Mi primera API Gateway",
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
      "@comment": "Servicio de middleware para validar y guardar el dispositivo basándose en las informaciones del encabezado.",
      "name": "device-manager",
      "hosts": [
        "http://192.168.1.2:8051"
      ],
      "path": "/devices",
      "method": "PUT",
      "forward-headers": [],
      "modifiers": {
        "@comment": "Primer modificador",
        "header": [
          {
            "@comment": "Añade el campo X-Device-Id con el valor del id del cuerpo de la respuesta del middleware",
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
      "@comment": "Endpoint responsable de devolver el usuario encontrado por la clave, que puede ser un correo electrónico o un número de teléfono.",
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
          "@comment": "Servicio de backend encargado de obtener el usuario a través de la clave.",
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
            "@comment": "Servicio de backend encargado de obtener el usuario a través de la clave.",
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

Campo obligatorio, para ayuda en la escritura y reglas del propio JSON de configuración, el único valor aceptado es

```text
https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json`
```

### @comment

Campo opcional, de tipo string, libre para hacer anotaciones relacionadas con su proyecto.

### version

Campo opcional, utilizado para control de versión y también usado en la respuesta del endpoint estático
[/version](#version-1).

### port

Campo obligatorio, utilizado para indicar el puerto que debe escuchar la API Gateway, valor mínimo `1` y valor
máximo `65535`.

### hot-reload

Campo opcional, el valor predeterminado es `false`, si es `true` se utiliza para la carga automática cuando
ocurre algún cambio en el archivo .json y .env en el directorio del entorno seleccionado.

### store

Campo opcional, el valor por defecto es el almacenamiento local en caché, si se informa, el campo `redis` pasa
a ser obligatorio, y también el campo `address` .

Si usas Redis para el almacenamiento global de caché, se recomienda que los valores de dirección y contraseña estén
completados usando la variable de ambiente, como en el ejemplo anterior.

### timeout

Campo opcional, el valor por defecto es `30 segundos`, este campo es responsable por el tiempo máximo de procesamiento
de cada petición.

Si la solicitud supera este tiempo establecido, la API Gateway abortará todas las transacciones en curso y
devolverá el código de estado `504 (Gateway Timeout)`. Consulte más sobre esta respuesta
haciendo [clic aquí](#504-gateway-timeout).

IMPORTANTE: Si se proporciona información en el objeto de endpoint, le daremos prioridad al valor informado del
endpoint, de lo contrario seguiremos con el valor informado o predeterminado de este campo, en la raíz del JSON de
configuración.

```text
- Valores aceptados:
    - s para segundos
    - m para minutos
    - h para horas
    - ms para milisegundos
    - us (o µs) para microsegundos
    - ns para nanosegundos

- Ejemplos:
    - 10s
    - 5ms
    - 1h30m
    - 1.5h
```

### cache

Campo opcional, si se proporciona, el campo `duration` se vuelve obligatorio!

Si el objeto se proporciona en la estructura de [endpoint.cache](#cache), damos prioridad a los valores proporcionados
allí, de lo contrario, continuaremos con los valores proporcionados en este campo.

El valor de la caché solo se graba una vez cada X duración proporcionada en el campo `every`.

Los campos `only-if-status-codes` y `only-if-methods` se utilizan para verificar si en ese endpoint habilitado
para tener caché, se puede leer y escribir la caché con base en método HTTP y código de estado de respuesta, vea más
sobre ellos a continuación.

Si la respuesta no es "fresca", es decir, fue respondida por la caché, el encabezado `X-Gopen-Cache` tendrá el
valor `true` de lo contrario el valor será `false`.

#### cache.duration

Indica cuánto tiempo durará la caché, es de tipo `time.Duration`.

```text
- Valores aceptados:
    - s para segundos
    - m para minutos
    - h para horas
    - ms para milisegundos
    - us (o µs) para microsegundos
    - ns para nanosegundos

- Ejemplos:
    - 10s
    - 5ms
    - 1h30m
    - 1.5h
```

#### cache.strategy-headers

Campo opcional, la estrategia de clave de caché por defecto es por la url y método HTTP de la solicitud convirtiéndolo
en un caché global por endpoint, si se proporcionan los encabezados a ser utilizados en la estrategia se agregan a los
valores por defecto de la clave, por ejemplo, en el ejemplo se indicó utilizar el campo `X-Forwarded-For` y `Device`, el
valor final de la clave sería:

      GET:/users/find/479976139:177.130.228.66:95D4AF55-733D-46D7-86B9-7EF7D6634EBC

La descripción de la lógica detrás de esta clave es:

      método:url:X-Forwarded-For:Device

Sin la estrategia completada, la lógica predeterminada es así:

      método:url

Entonces el valor predeterminado para este endpoint sin la estrategia completada es:

      GET:/users/find/479976139

¡En este ejemplo, convertimos el caché, que antes era global para el endpoint en específico, para que ahora sea por
cliente! Recuerda que esto es solo un ejemplo simple, puedes tener la estrategia que quieras basada en el encabezado de
tu aplicación.

#### cache.only-if-methods

Campo opcional, el valor por defecto es una lista con solo el método HTTP `GET`, si está vacía, se aceptará cualquier
método HTTP.

Este campo es responsable de decidir si se leerá y se escribirá la caché del endpoint (que está habilitado para tener
caché) por el método HTTP del mismo.

#### cache.only-if-status-codes

Campo opcional, el valor por defecto es una lista de códigos de estado HTTP de éxito reconocidos, si está vacía,
cualquier código de estado HTTP de respuesta será aceptado.

Este campo es responsable de decidir si vamos a grabar o no la caché del endpoint (que está habilitado para tener caché)
por el código de estado HTTP de la respuesta del mismo.

#### cache.allow-cache-control

Campo opcional, el valor por defecto es `false`, si se establece como `true`, la API Gateway considerará el encabezado
`Cache-Control` siguiendo las siguientes reglas a partir del valor informado en la solicitud o en la respuesta de los
backends:

`no-cache`: este valor solo se considera en el header de la solicitud, si se informa, desconsideramos la lectura de la
caché y continuamos con el proceso normal para obtener la respuesta "fresca".

`no-store`: este valor es considerado solo en la respuesta escrita por sus backends, si se informa, no guardamos la
caché.

### limiter

Campo opcional, objeto responsable de las reglas de limitación de la API Gateway, ya sea de tamaño o tasa, los valores
predeterminados varían de campo a campo, vea:

#### limiter.max-header-size

Campo opcional, es del tipo `byteUnit`, el valor por defecto es `1MB`, es responsable de limitar el tamaño del
encabezado de la solicitud.

Si el tamaño del encabezado supera el valor informado, la API Gateway abortará la solicitud con el código de estado
`431 (Request header fields too large)`. Obtenga más información sobre esta devolución haciendo
[clic aquí](#431-request-header-fields-too-large).

```text
- Valores aceptados:
    - B para Byte
    - KB para KiloByte
    - MB para MegaByte
    - GB para GigaByte
    - TB para TeraByte
    - PB para PetaByte
    - EB para ExaByte
    - ZB para ZettaByte
    - YB para YottaByte

- Ejemplos:
    - 1B
    - 50KB
    - 5MB
    - 1.5GB
```

#### limiter.max-body-size

Campo opcional, es del tipo `byteUnit`, el valor por defecto es `3MB`, el campo es responsable de limitar el tamaño del
cuerpo de la solicitud.

Si el tamaño del cuerpo supera el valor informado, la API Gateway abortará la solicitud con el código de estado
`413 (Request entity too large)`.

```text
- Valores aceptados:
    - B para Byte
    - KB para KiloByte
    - MB para MegaByte
    - GB para GigaByte
    - TB para TeraByte
    - PB para PetaByte
    - EB para ExaByte
    - ZB para ZettaByte
    - YB para YottaByte

- Ejemplos:
    - 1B
    - 50KB
    - 5MB
    - 1.5GB
```

#### limiter.max-multipart-memory-size

Campo opcional, es del tipo `byteUnit`, el valor por defecto es `5MB`, este campo es responsable de limitar el tamaño
del cuerpo multipart/form de la solicitud, generalmente utilizado para el envío de archivos, imágenes, etc.

Si el tamaño del cuerpo supera el valor informado, la API Gateway abortará la solicitud con el código de
estado `413 (Request entity too large)`. Obtenga más información sobre esta devolución haciendo
[clic aquí](#413-request-entity-too-large).

```text
- Valores aceptados:
    - B para Byte
    - KB para KiloByte
    - MB para MegaByte
    - GB para GigaByte
    - TB para TeraByte
    - PB para PetaByte
    - EB para ExaByte
    - ZB para ZettaByte
    - YB para YottaByte

- Ejemplos:
    - 1B
    - 50KB
    - 5MB
    - 1.5GB
```

#### limiter.rate

Campo opcional, si se proporciona, el campo `capacity` se vuelve obligatorio. Este objeto es responsable de limitar la
tasa de solicitud por IP, este límite se impone obteniendo la capacidad máxima del campo `capacity` por X duración,
proporcionada en el campo `every`.

Si se supera esta capacidad, la API Gateway por seguridad abortará la solicitud, devolviendo `429 (Too many requests)`.
Obtenga más información sobre esta devolución haciendo [clic aquí](#429-too-many-requests).

#### limiter.rate.capacity

Campo opcional, si se proporciona el objeto `rate`, este se vuelve obligatorio. El valor por defecto es `5`, y el mínimo
que se puede proporcionar es `1`. Indica la capacidad máxima de solicitudes.

#### limiter.rate.every

Campo opcional, el valor por defecto es `1 segundo`, indica el valor de la duración de la verificación de la capacidad
máxima de solicitudes.

### security-cors

Campo opcional, utilizado para la seguridad CORS de la API Gateway, todos los campos son vacíos por defecto, no
restringiendo los valores de `origin`, `methods` y `headers`.

Si desea restringir y la solicitud no coincide con las configuraciones establecidas, la API Gateway por seguridad
abortará la solicitud, devolviendo `403 (Prohibido)`.

#### security-cors.allow-origins

Campo opcional, de tipo lista de strings, los elementos de la lista deben indicar los IPs de origen que la API Gateway
permite recibir en las solicitudes.

#### security-cors.allow-methods

Campo opcional, de tipo lista de strings, los elementos de la lista deben indicar los métodos HTTP que la API Gateway
permite recibir en las solicitudes.

#### security-cors.allow-headers

Campo opcional, de tipo lista de strings, los elementos de la lista deben indicar qué campos del header HTTP la API
Gateway permite recibir en las solicitudes.

### middlewares

Campo opcional, es responsable de la configuración de sus middlewares de aplicación, es un mapa con claves en cadena de
texto mencionando el nombre de su middleware, ese nombre podrá ser utilizado en sus [endpoints](#endpoints)
como `beforewares` y `afterwares`.

El valor de la clave es un objeto de [backend](#backendname), sin embargo, con una observación, este objeto tendrá su
respuesta, en caso de no ser [abortada](#respuesta-abortada), omitida automáticamente por el endpoint, solo su header de
respuesta podrá ser devuelto, en caso de tener uno, ya que las respuestas de los middleware se deben quedar ocultas al
cliente final HTTP, sin embargo, su respuesta será almacenada a lo largo de la solicitud HTTP hecha en el endpoint,
pudiendo obtener y manipular sus valores de solicitud y respuesta.

Por ejemplo, un `beforeware` cuando se menciona en el endpoint, se utilizará como middleware de pre-solicitudes, es
decir, se llamará antes de los backends principales del endpoint, entonces podemos, por ejemplo, tener un middleware de
manipulación de dispositivo, como en el JSON de configuración de arriba, donde llamará a este backend de middleware
configurado en el endpoint como `beforewares`, validando y guardando el dispositivo a partir de informaciones del header
de la solicitud, en caso de que el backend responda un código de estado de falla, en el ejemplo, el gateway abortará
todos los backends siguientes retornando lo que el backend de dispositivo respondió, en caso de que haya retornado un
código de estado de éxito, modificará el header de todas las solicitudes siguientes (`propagate:true`), añadiendo el
campo `X-Device-Id`, con el valor del ID del cuerpo de respuesta del propio backend, pudiendo ser utilizado en los otros
backends siguientes del endpoint.

Para saber más sobre los `modifiers` [veja](#backendmodifiers).

Para entender mejor esta herramienta poderosa en la práctica, vea los ejemplos de middlewares usados como `beforewares`
y `afterwares` hechos en el proyecto de [playground](https://github.com/tech4works/gopen-gateway-playground).

### endpoints

Campo obligatorio, es una lista de objetos, representa cada endpoint de la API Gateway que se registrará para escuchar y
servir las solicitudes HTTP.

Vea cómo funciona el flujo básico de un endpoint en la imagen a continuación:

<img src="assets/basic-endpoint-flowchart-es.png" alt="flujo básico de un endpoint">

A continuación vamos a listar y explicar cada campo de este objeto tan importante:

### endpoint.@comment`

Campo opcional, de tipo string, campo libre para anotaciones relacionadas con su endpoint.

### endpoint.path

Campo obligatorio, de tipo string, responsable del camino URI del endpoint, ejemplo `"/users/:id"`.

Si desea tener parámetros dinámicos en este endpoint, simplemente use el estándar `":nombre del parámetro"` por
ejemplo `"/users/:id/status/:status"`, la API Gateway entenderá que tendremos 2 parámetros dinámicos de este endpoint,
estos valores pueden ser pasados a los backends subyacentes.

Ejemplo usando el parámetro dinámico para los backends subyacentes:

- Endpoint
    - path: `"/users/:id/status/:status"`
    - resultado: `"/users/1/status/removed"`
- Backend 1
    - path: `"/users/:id"`
    - resultado: `"/users/1"`
- Backend 2
    - path: `"/users/:id/status/:status"`
    - resultado: `"/users/1/status/removed"`

En el ejemplo anterior vemos que el parámetro puede ser utilizado como quiera como path en las solicitudes de backend
del endpoint en cuestión.

### endpoint.method

Campo obligatorio, de tipo string, responsable de definir qué método HTTP será registrado el endpoint.

### endpoint.timeout

Es similar al campo [timeout](#timeout), sin embargo, será aplicado solo para el endpoint en cuestión.

Si se omite, se heredará el valor del campo [timeout](#timeout).

### endpoint.cache

Campo opcional, de tipo objeto, por defecto vendrá vacío solo con el campo `enabled` lleno con el valor `false`.

Si se informa, el campo `enabled` se vuelve obligatorio, los otros campos, si se omiten, heredará de la
configuración [cache](#cache) en la raíz si existe y si está llenada.

Si por casualidad, ha omitido el campo `duration` tanto en la configuración actual como en la
configuración [cache](#cache) en la raíz, el campo `enabled` es ignorado considerándose siempre como `false` ya que no
se informó la duración del caché en ambas configuraciones.

### endpoint.cache.enabled

Campo obligatorio, de tipo booleano, indica si desea que tenga caché en su endpoint, `true` para habilitado, `false`
para no habilitado.

Si está en `true` pero no se informa el campo `duration` en la configuración actual ni en [la raíz](#cache), este campo
será ignorado considerándose siempre como `false`.

### endpoint.cache.ignore-query

Campo opcional, de tipo booleano, si no se informa el valor por defecto es `false`.

Si el valor es `true`, la API Gateway ignorará los parámetros de búsqueda de la URL al crear la clave de almacenamiento,
de lo contrario, considerará los parámetros de búsqueda de la URL ordenando alfabéticamente las claves y valores.

### endpoint.cache.duration

Es similar al campo [cache.duration](#cacheduration), sin embargo, se aplicará solo para el endpoint en cuestión.

Si se omite, heredará el valor del campo [cache.duration](#cacheduration).

Si se omite en las dos configuraciones, el campo `enabled` será ignorado considerándose siempre como `false`.

### endpoint.cache.strategy-headers

Es similar al campo [cache.strategy-headers](#cachestrategy-headers), sin embargo, se aplicará solo para el endpoint en
cuestión.

Si se omite, heredará el valor del campo [cache.strategy-headers](#cachestrategy-headers).

Si se informa vacío, el valor no será heredado, sin embargo, se aplicará el valor [default](#cachestrategy-headers) para
el endpoint en cuestión.

### endpoint.cache.only-if-status-codes

Es similar al campo [cache.only-if-status-codes](#cacheonly-if-status-codes), sin embargo, se aplicará solo para el
endpoint en cuestión.

Si se omite, heredará el valor del campo [cache.only-if-status-codes](#cacheonly-if-status-codes).

Si se informa vacío, el valor no será heredado, sin embargo, se aplicará el valor [default](#cacheonly-if-status-codes)
para el endpoint en cuestión.

### endpoint.cache.allow-cache-control

Es similar al campo [cache.allow-cache-control](#cacheallow-cache-control), sin embargo, se aplicará solo para el
endpoint en cuestión.

Si se omite, heredará el valor del campo [cache.allow-cache-control](#cacheallow-cache-control).

### endpoint.limiter

Es similar al campo [limiter](#limiter), sin embargo, se aplicará solo para el endpoint en cuestión.

Si se omite, heredará el valor del campo [limiter](#limiter).

### endpoint.response-encode

Campo opcional, de tipo string, el valor por defecto es vacío, indicando que la respuesta del endpoint será codificada
siguiendo la [lógica de respuesta](#lógica-de-respuesta) de la API Gateway, sin forzar la codificación indicada.

```text
- Valores aceptados:
  - JSON 
  - XML
  - TEXT
```

### endpoint.aggregate-responses

Campo opcional, de tipo booleano, el valor por defecto es `false`, indicando que la respuesta del endpoint no será
agregada.

Si se informa con el valor `true` y tiene más de una respuesta de los backends informados en el endpoint él agregará las
respuestas de los backends, vea más sobre las reglas de respuesta de la API Gateway haciendo
[clic aquí](#lógica-de-respuesta).

### endpoint.abort-if-status-codes

Campo opcional, de tipo lista de enteros, el valor por defecto es nulo, indicando que cualquier backend ejecutado en el
endpoint que haya respondido el código de estado HTTP mayor o igual a `400 (Bad request)` será abortado.

Si se informa, y un backend retorna el código de estado HTTP indicado en la configuración, el endpoint será abortado,
esto significa que los otros backends configurados después del mismo, no serán ejecutados, y el endpoint retornará la
respuesta del mismo al cliente final.

Si desea que ningún código de estado HTTP sea abortado en el endpoint, simplemente informe el campo vacío.

Vea cómo el endpoint será respondido después de que un backend sea abortado haciendo [clic aquí](#lógica-de-respuesta).

### endpoint.beforewares

Campo opcional, de tipo lista de string, el valor por defecto es vacío, indicando que el endpoint no tiene ningún
middleware de pre-solicitudes.

Si se informa, el endpoint ejecutará las solicitudes, posición por posición, comenzando en el inicio de la lista. Si el
valor en string de la posición a ser ejecutada está configurado en el campo [middlewares](#middlewares) correctamente,
se ejecutará el backend configurado en el mismo. De lo contrario ignorará la posición solo imprimiendo un log de
atención.

### endpoint.afterwares

Campo opcional, de tipo lista de string, el valor por defecto es vacío, indicando que el endpoint no tiene ningún
middleware de post-solicitudes.

Si se informa, el endpoint ejecutará las solicitudes, posición por posición, comenzando en el inicio de la lista. Si el
valor en string de la posición a ser ejecutada está configurado en el campo [middlewares](#middlewares) correctamente,
se ejecutará el backend configurado en el mismo. De lo contrario ignorará la posición solo imprimiendo un log de
atención.

### endpoint.backends

Campo obligatorio, de tipo lista de objetos, responsable de la ejecución principal del endpoint, el propio nombre ya
dice todo, es una lista que indica todos los servicios necesarios para que el endpoint retorne la respuesta esperada.

Vea cómo funciona el flujo básico de un backend en la imagen a continuación:

<img src="assets/basic-backend-flowchart-es.png" alt="flujo básico de un backend">

A continuación, vamos a listar y explicar cada campo de este objeto tan importante:

### backend.@comment

Campo opcional, de tipo string, campo libre para anotaciones relacionados con su backend.

### backend.name

Campo opcional, de tipo string, es responsable del nombre de su servicio backend, se utiliza para dar nombre al campo de
respuesta agregada del mismo, en caso de que el
campo [backend.extra-config.group-response](#backendextra-configgroup-response) esté como `true`.

### backend.hosts

Campo obligatorio, de tipo lista de strings, es responsable de los hosts de su servicio que la API Gateway llamará en
conjunto con el campo [backend.path](#backendpath).

De alguna manera podemos tener un balanceador de carga "tonto", ya que el backend seleccionará al azar en esta lista qué
host será llamado, con esto podemos informar múltiples veces el mismo host para equilibrar las llamadas, vea:

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

Campo obligatorio, de tipo string, el valor indica la URL del camino del servicio backend.

Usamos uno de los [backend.hosts](#backendhosts) informados y los juntamos con el path proporcionado, por ejemplo, en el
campo hosts tenemos el valor

```text
[
  "https://instance-01", 
  "https://instance-02"
]
```

Y en este campo de path, tenemos el valor

```text
/users/status
```

El backend construirá la siguiente URL de solicitud

```text
https://instance-02/users/status
```

Vea cómo se equilibra el host haciendo [clic aquí](#backendhosts).

### backend.method

Campo obligatorio, de tipo string, el valor indica qué método HTTP espera el servicio backend.

### backend.forward-queries

Campo opcional, tipo lista de string, el valor predeterminado está vacío, lo que indica que cualquier parámetro de
búsqueda será reenviado al servicio backend.

Si se informa, solo los campos indicados serán reenviados al servicio backend, por ejemplo, recibimos una
solicitud con la siguiente URL

````
/users?id=23&email=gabrielcataldo@gmail.com&phone=47991271234
````

En este ejemplo, tenemos el campo `forward-queries` con los siguientes valores

````
[
  "email",
  "phone"
]
````

La URL de solicitud al backend fue

````
/users?email=gabrielcataldo@gmail.com&phone=47991271234
````

Vemos que el parámetro de búsqueda `id` no fue reenviado al servicio backend, ya que no fue mencionado en la lista.

### backend.forward-headers

Campo opcional, tipo lista de cadenas, el valor predeterminado está vacío, lo que indica que cualquier encabezado
recibido será reenviado al servicio backend.

Si se informa, solo los campos indicados serán reenviados al servicio backend, por ejemplo, recibimos una solicitud con
el siguiente encabezado

````
{
  "Device": "95D4AF55-733D-46D7-86B9-7EF7D6634EBC",
  "User-Agent": "IOS",
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
````

En este ejemplo, tenemos el campo `forward-headers` con los siguientes valores:

````
[
  "User-Agent",
  "Authorization"
]
````

El encabezado de la solicitud al backend fue:

```json
{
  "User-Agent": "IOS",
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

Vemos que el campo `Device` del encabezado recibido no se reenvió al servicio backend, ya que no se mencionó en la
lista.

### backend.extra-config

Campo opcional, de tipo objeto, indica configuraciones adicionales del servicio backend. A continuación, vea los campos
y sus responsabilidades:

#### backend.extra-config.omit-request-body

Campo opcional, de tipo booleano, el valor predeterminado es `false`, lo que indica que el cuerpo de la solicitud se
pasará al backend si lo hay.

Si se informa `true`, el cuerpo de la solicitud no se pasará al backend.

#### backend.extra-config.group-response

Campo opcional, de tipo booleano, el valor predeterminado es `false`, lo que indica que no será necesario agrupar el
cuerpo de la respuesta del backend en un campo JSON para la respuesta al cliente final.

Si se informa `true`, el cuerpo de la respuesta, en caso de que exista, se agrupará en un campo JSON de la respuesta
final, el nombre del campo será el [nombre](#backendname) del servicio backend si está lleno, si no tenemos un estándar
de nomenclatura que es `backend-posición en la lista` que sería, por ejemplo, `backend-0`.

Para entender la importancia de este campo, vea más sobre la [lógica de respuesta](#lógica-de-respuesta) de la API
Gateway.

#### backend.extra-config.omit-response

Campo opcional, de tipo booleano, el valor predeterminado es `false`, lo que indica que la respuesta del backend de la
pregunta no será omitida para el cliente final.

Si se informa `true`, toda la respuesta del backend en cuestión será omitida, tenga cuidado, porque si solo tiene este
backend, y se omite, la API Gateway responderá por defecto con el código de estado HTTP `204 (No Content)`.

Para entender la importancia de este campo, vea más sobre la [lógica de respuesta](#lógica-de-respuesta) de la API
Gateway.

### backend.modifiers

Campo opcional, de tipo objeto, el valor predeterminado es vacío, indicando que no hay ningún proceso de modificación en
este backend en cuestión.

A continuación, enumeraremos y explicaremos cada campo de este importante objeto:

### modifiers.@comment

Campo opcional, de tipo string, espacio libre para anotaciones relacionadas con tus modificadores.

### modifiers.status-code

Campo opcional, de tipo entero, el valor predeterminado es `0`, indicando que no hay nada que modificar en el código de
estado HTTP de la respuesta del backend.

Si se informa, el código de estado HTTP de la respuesta del backend se modificará por el valor insertado, esto puede
influir o no en la respuesta final del endpoint, vea la [lógica-de-respuesta](#lógica-de-respuesta) de la API Gateway
para saber más.

### modifiers.header

Campo opcional, de tipo lista de objetos, valor predeterminado es vacío, responsable por las modificaciones de
encabezado de la solicitud y respuesta del backend.

Veamos a continuación los campos de este objeto y sus responsabilidades:

#### header.@comment

Campo opcional, de tipo string, espacio libre para anotaciones relacionadas con tu modificador.

#### header.context

Campo obligatorio, de tipo string, es responsable de indicar en qué contexto debe actuar la modificación.

Valores aceptados:

`REQUEST`: Para actuar en la pre-solicitud al backend.

`RESPONSE`: Para actuar en la post-solicitud al backend.

Es importante recordar que si el valor es `REQUEST`, solo se puede usar el valor `REQUEST` en el
campo [header.scope](#headerscope).

#### header.scope

Campo opcional, de tipo string, el valor predeterminado se basará en el campo [header.context](#headercontext)
informado. El valor indica qué ámbito debemos cambiar, ya sea el ámbito de la solicitud o de la respuesta.

Valores aceptados:

`REQUEST`: Para modificar el ámbito de la solicitud, este tipo de ámbito puede tener una actuación global propagando
este cambio en las subsiguientes solicitudes al backend, solo es necesario utilizar el
campo [header.propagate](#headerpropagate) como `true`.

`RESPONSE`: Para modificar el ámbito de la respuesta del backend.

#### header.action

Campo obligatorio, de tipo string, responsable de la acción a tomar en la modificación del encabezado.

Valores aceptados:

`ADD`: Agrega la clave informada en el campo [header.key](#headerkey) si no existe, y agrega el valor informado en el
campo [header.value](#headervalue).

`APD`: Añade el valor informado en el campo [header.value](#headervalue) si la clave informada en el
campo [header.key](#headerkey) existe.

`SET`: Define el valor de la clave informada en el campo [header.key](#headerkey) con el valor pasado en el
campo [header.value](#headervalue).

`RPL`: Reemplaza el valor de la clave informada en el campo [header.key](#headerkey) con el valor pasado en el
campo [header.value](#headervalue) si existe.

`REN`: Renombra la clave informada en el campo [header.key](#headerkey) con el valor pasado en el
campo [header.value](#headervalue) si existe.

`DEL`: Elimina la clave informada en el campo [header.key](#headerkey) si existe.

#### header.key

Campo obligatorio, de tipo string, utilizado para indicar cuál clave del encabezado debe ser modificada.

#### header.value

Campo obligatorio, de tipo cadena (string), utilizado como valor para modificar la clave indicada en el
campo [header.key](#headerkey).

Existen posibilidades de uso de [valores dinámicos](#valores-dinámicos-para-modificación), y
de [variables de entorno](#variables-de-entorno) para este campo.

NOTA: se vuelve opcional solo si [query.action](#queryaction) tiene el valor `DEL`.

#### header.propagate

Campo opcional, de tipo booleano, el valor predeterminado es `false`, indicando que el modificador no debe propagar este
cambio a los siguientes backends.

Si se informa como `true`, este cambio se propagará a los siguientes backends.

IMPORTANTE: Este campo solo es aceptado si el [ámbito](#headerscope) tiene el valor `REQUEST`.

### modifiers.param

Campo opcional, de tipo lista de objetos, el valor predeterminado es vacío, responsable de las modificaciones de
parámetros de la solicitud para el servidor backend.

A continuación, puede encontrar los campos de este objeto y sus responsabilidades:

#### param.@comment

Campo opcional, de tipo string, campo libre para anotaciones relacionadas con su modificador.

#### param.context

Campo obligatorio, de tipo string, es responsable de indicar en qué contexto debe actuar la modificación.

Valores aceptados:

`REQUEST`: Para actuar en la pre-solicitud al backend.

`RESPONSE`: Para actuar después de la solicitud al backend.

#### param.action

Campo obligatorio, de tipo string, responsable de la acción a ser tomada en la modificación de los parámetros
de la solicitud.

Valores aceptados:

`SET`: Establece el valor de la clave informada en el campo [param.key](#paramkey) por el valor pasado en el
campo [param.value](#paramvalue).

`RPL`: Reemplaza el valor de la clave informada en el campo [header.key](#headerkey) por el valor pasado en el
campo [header.value](#headervalue) si existe.

`REN`: Renombra la clave informada en el campo [param.key](#paramkey) por el valor pasado en el
campo [param.value](#paramvalue) si existe.

`DEL`: Elimina la clave informada en el campo [param.key](#paramkey) si existe.

#### param.key

Campo obligatorio, de tipo string, utilizado para indicar qué clave de parámetro debe ser modificada.

#### param.value

Campo obligatorio, de tipo string, utilizado como valor para ser usado para modificar la clave indicada en el
campo [param.key](#paramkey).

Tenemos posibilidades de uso de [valores dinámicos](#valores-dinámicos-para-modificación), y
de [variables de entorno](#variables-de-entorno) para este campo.

NOTA: se vuelve opcional solo si [query.action](#queryaction) tiene el valor `DEL`.

#### param.propagate

Campo opcional, de tipo booleano, el valor predeterminado es `false`, lo que indica que el modificador no debe propagar
este cambio a los backends subsiguientes.

Si se informa como `true`, esta modificación se propagará a los siguientes backends.

### modifiers.query

Campo opcional, de tipo lista de objetos, el valor predeterminado es vacío, responsable de las modificaciones de
parámetros de búsqueda de la solicitud para el backend.

Vea a continuación los campos de este objeto y sus responsabilidades:

#### query.@comment

Campo opcional, de tipo string, un campo libre para anotaciones relacionadas con su modificador.

#### query.context

Campo obligatorio, de tipo string, es responsable de indicar en qué contexto debe actuar la modificación.

Valores aceptados:

`REQUEST`: Para actuar en la pre-solicitud del backend.

`RESPONSE`: Para actuar después de la solicitud al backend.

#### query.action

Campo obligatorio, de tipo string, responsable de la acción que se debe tomar en la modificación de los
parámetros de búsqueda de la solicitud.

Valores aceptados:

`ADD`: Añade la clave informada en el campo [query.key](#querykey) en caso de que no exista, y agrega el valor informado
en el campo [query.value](#queryvalue).

`APD`: Aumenta el valor informado en el campo [query.value](#queryvalue) en caso de que la clave informada en el
campo [query.key](#querykey) exista.

`SET`: Establece el valor de la clave informada en el campo [query.key](#querykey) por el valor dado en el
campo [query.value](#queryvalue).

`RPL`: Reemplaza el valor de la clave informada en el campo [query.key](#querykey) por el valor dado en el
campo [query.value](#queryvalue) en caso de existir.

`REN`: Renombra la clave informada en el campo [query.key](#querykey) por el valor dado en el
campo [query.value](#queryvalue) en caso de existir.

`DEL`: Elimina la clave informada en el campo [query.key](#querykey) en caso de existir.

#### query.key

Campo obligatorio, de tipo string, utilizado para indicar cuál clave de parámetro de búsqueda debe ser modificada.

#### query.value

Campo obligatorio, de tipo string, utilizado como valor para modificar la clave indicada en el
campo [query.key](#querykey).

Existen posibilidades de uso de [valores dinámicos](#valores-dinámicos-para-modificación) y
de [variables de entorno](#variables-de-entorno) para este campo.

NOTA: Se vuelve opcional solo si [query.action](#queryaction) tiene el valor `DEL`.

#### query.propagate

Campo opcional, de tipo booleano, el valor predeterminado es `false`, lo que indica que el modificador no debería
propagar este cambio a los siguientes backends.

Si se informa como `true`, esta modificación será propagada a los backends siguientes.

### modifiers.body

Campo opcional, de tipo lista de objetos, el valor predeterminado es vacío, responsable de las modificaciones del cuerpo
de la solicitud o respuesta del backend.

Vea a continuación los campos de este objeto y sus responsabilidades:

#### body.@comment

Campo opcional, de tipo string, espacio libre para anotaciones relacionadas con su modificador.

#### body.context

Campo obligatorio, de tipo string, es responsable de indicar en qué contexto debe actuar la modificación.

Valores aceptados:

`REQUEST`: Para actuar antes de la solicitud al backend.

`RESPONSE`: Para actuar después de la solicitud al backend.

Es importante recordar que si el valor es `REQUEST`, solo puede usar el valor `REQUEST` en el
campo [body.scope](#bodyscope).

#### body.scope

Campo opcional, de tipo string, el valor predeterminado se basará en el campo [body.context](#bodycontext)
proporcionado, indicando qué ámbito debemos cambiar, ya sea el ámbito de la solicitud o de la respuesta.

Valores aceptados:

`REQUEST`: Para modificar el ámbito de la solicitud, este tipo de ámbito puede tener un efecto global propagando este
cambio a las siguientes solicitudes al backend, simplemente use el campo [body.propagate](#bodypropagate) como `true`.

`RESPONSE`: Para modificar el ámbito de la respuesta del backend.

#### body.action

Campo obligatorio, de tipo string, responsable de la acción a tomar en la modificación del cuerpo (body).

Valores aceptados si el cuerpo es JSON:

`ADD`: Agrega la clave proporcionada en el campo [body.key](#bodykey) si no existe, y agrega el valor proporcionado en
el campo [body.value](#bodyvalue).

`APD`: Añade el valor proporcionado en el campo [body.value](#bodyvalue) si la clave proporcionada en el
campo [body.key](#bodykey) existe.

`SET`: Define el valor de la clave proporcionada en el campo [body.key](#bodykey) con el valor pasado en el
campo [body.value](#bodyvalue).

`RPL`: Reemplaza el valor de la clave proporcionada en el campo [body.key](#bodykey) con el valor pasado en el
campo [body.value](#bodyvalue) si existe.

`REN`: Renombra la clave proporcionada en el campo [body.key](#bodykey) con el valor pasado en el
campo [body.value](#bodyvalue) si existe.

`DEL`: Elimina la clave proporcionada en el campo [body.key](#bodykey) si existe.

Valores aceptados si el cuerpo es TEXTO:

`ADD`: Agrega el valor proporcionado en el campo [body.value](#bodyvalue) al texto.

`APD`: Añade el valor proporcionado en el campo [body.value](#bodyvalue) si el cuerpo no está vacío.

`RPL`: Reemplazará todos los valores que se parezcan a la clave proporcionada en el campo [body.key](#bodykey) con el
valor pasado en el campo [body.value](#bodyvalue).

`DEL`: Elimina todos los valores que se parezcan a la clave proporcionada en el campo [body.key](#bodykey).

#### body.key

Campo obligatorio, de tipo string, utilizado para indicar qué clave del encabezado debe ser modificada.

NOTA: se vuelve opcional si el cuerpo (body) es del tipo TEXTO y [body.action](#bodyaction) tiene el valor `ADD`.

#### body.value

Campo obligatorio, de tipo string, se utiliza como valor que se usará para modificar la clave indicada en el
campo [body.key](#bodykey).

Tenemos posibilidades de usar [valores dinámicos](#valores-dinámicos-para-modificación)
y [variables de entorno](#variables-de-entorno) para este campo.

NOTA: Se vuelve opcional solo si [body.action](#bodyaction) tiene el valor `DEL`.

#### body.propagate

Campo opcional, de tipo booleano, el valor predeterminado es `false`, indicando que el modificador no debe propagar este
cambio a los backends siguientes.

Si se informa como `true`, esta modificación se propagará a los siguientes backends.

IMPORTANTE: Este campo solo se acepta si [body.scope](#bodyscope) tiene el valor `REQUEST`.

## JSON en tiempo de ejecución

El GOPEN API Gateway, cuando se inicia, genera un archivo JSON, basado en
el [JSON de configuración](#json-de-configuración), ubicado en la carpeta `runtime` en la raíz de su área de trabajo.

Este JSON indica cuál fue la comprensión de la aplicación al leer el [JSON de configuración](#json-de-configuración),
todas las [#variables de entorno](#variables-de-entorno) ya tendrán sus valores reemplazados, si existen.

Este json también se puede leer utilizando la ruta estática [/settings](#settings).

## Rutas estáticas

El GOPEN API Gateway tiene algunos endpoints estáticos, es decir, independientemente de cualquier configuración
realizada, tendremos actualmente hay tres puntos finales registrados en sus rutas, consulte cada uno y sus
responsabilidades a continuación:

### `/ping`

Endpoint para saber si API Gateway está activo, la ruta devuelve `404 (No encontrado)` si está desactivado, y
`200 (OK)` si está al aire.

### `/version`

Endpoint que devuelve la versión obtenida en la configuración [versión](#version), devuelve `404 (Not Found)` si no se
ha encontrado informado en el [JSON de configuración](#json-de-configuración), de lo contrario devuelve `200 (OK)` con
el valor en el cuerpo como texto.

### `/settings`

Endpoint devuelve cierta información sobre el proyecto, como la versión, la fecha de la versión, el número de
contribuyentes y un resumen de cuántos puntos finales, middleware, backends y modificadores están configurados
actualmente y el json de configuración que está funcionando activamente.

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

## Variables de entorno

Se pueden crear instancias de variables de entorno fácilmente utilizando el archivo .env, en la carpeta indicada por el
entorno. Inicialización dinámica como se menciona en el tema [ESTRUCTURA DE CARPETAS](#estructura-de-carpetas).

Si prefieres insertar los valores usando docker-compose también funcionará correctamente, la cuestión es que la API
Gateway leerá el valor escrito en la máquina, independientemente de cómo se haya ingresado.

Los valores se pueden usar en la configuración JSON de API Gateway, solo use la sintaxis `$NAME` como
un valor de cadena, consulte el ejemplo siguiente.

Un extracto de una configuración JSON, tengo los siguientes valores:

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

Y en nuestra máquina tenemos las siguientes variables de entorno:

```dotenv
VERSION=1.0.0

REDIS_URL=redis-18885.c259.us-east-1-4.ec2.cloud.redislabs.com:18985
REDIS_PASSWORD=12345

TIMEOUT=5m
```

La API Gateway genera un archivo de [JSON en tiempo de ejecución](#json-en-tiempo-de-ejecución) al ejecutar la
aplicación, vea el resultado del mismo después de iniciar la aplicación:

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

Observamos que todos los valores con la sintaxis `$NOME` fueron reemplazados por sus respectivos valores. Si un valor se
menciona con esta sintaxis, pero no existe en las variables de entorno, se mantendrá el mismo valor informado.

### Observación

En esta primera versión, podemos utilizar esta sintaxis solo en campos de tipo string. Si
desea [contribuir](#cómo-contribuir) con una solución que aceptará esta sintaxis independientemente del tipo, ¡estaremos
agradecidos!

## Valores dinámicos para modificación

Podemos utilizar los valores de solicitud y respuesta del tiempo de ejecución del punto final, tal como se configuró.
Estos valores pueden ser obtenidos a través de una sintaxis específica, tenemos las siguientes posibilidades para
obtener estos valores, mira:

### Solicitud

Cuando mencionas la sintaxis `#request...` estarás obteniendo los valores de la solicitud recibida, si mencionas el
campo `#request.history...` estarás obteniendo los valores
del [historial temporal de solicitudes](#historial-temporal-de-solicitudes).

`#request.header...`

Este fragmento de sintaxis obtendrá del encabezado de la solicitud el valor indicado, por ejemplo,
`#request.header.X-Forwarded-For.0` obtendrá el primer valor del campo `X-Forwarded-For` del encabezado de la solicitud
si existe, sustituyendo la sintaxis por el valor, el resultado fue `127.0.0.1`.

`#request.params...`

Este fragmento de sintaxis obtendrá de los parámetros de la solicitud el valor indicado, por ejemplo,
`#request.params.id` obtendrá el valor del campo `id` de los parámetros de la solicitud si existe,
sustituyendo la sintaxis por el valor, el resultado fue `72761`.

`#request.query...`

Este fragmento de sintaxis obtendrá de los parámetros de búsqueda de la solicitud el valor indicado, por ejemplo,
`#request.query.email.0` obtendrá el primer valor del campo `email` de los parámetros de búsqueda de la solicitud si
existe, sustituyendo la sintaxis por el valor, el resultado fue `gabrielcataldo.adm@gmail.com`.

`#request.body...`

Este fragmento de sintaxis obtendrá del cuerpo de la solicitud el valor indicado, por ejemplo,
`#request.body.deviceId` obtendrá el valor del campo `deviceId` del cuerpo de la solicitud si existe,
sustituyendo la sintaxis por el valor, el resultado fue `991238`.

#### Historial temporal de solicitudes

Si quieres obtener los valores de la solicitud específicamente de un backend, simplemente añade el campo `history` en la
sintaxis, y la posición de tu backend en el historial, recordando que los middlewares `beforewares` y `afterwares`
cuentan como un backend en el historial, otro punto es que solo tendrás la información de la solicitud en este campo si
el backend ya ha sido procesado, o el backend está en proceso. Mira cómo quedaría un ejemplo a continuación:

En el ejemplo, tengo solo un backend y está siendo procesado, entonces puedo estar utilizando la sintaxis:

`#request.history.0.params.phone`

En este otro ejemplo de sintaxis tenemos tres backends configurados y dos ya han sido procesados, entonces podemos
utilizar la siguiente sintaxis en el proceso del tercer backend:

`#request.history.1.body.id`

En estos ejemplos citados vemos que podemos obtener el valor de la solicitud que está siendo realizada o ya fue
realizada por un backend configurado, y que están almacenados en un tipo de historial temporal.

### Respuesta

Cuando mencionas la sintaxis `#response...` estarás obteniendo los valores de la respuesta del API Gateway preparada con
base en la [lógica de respuesta](#lógica-de-respuesta), si mencionas el campo `#response.history...` estarás obteniendo
los valores del [historial temporal de respuestas](#historial-temporal-de-respuestas).

`#response.statusCode`

Esta sintaxis obtendrá el código de estado HTTP de la respuesta actual del endpoint. Si quieres obtener la respuesta de
un backend específico, usa el [historial temporal de respuestas](#historial-temporal-de-respuestas).

`#response.header...`

Este fragmento de la sintaxis obtendrá del encabezado de respuesta del endpoint el valor indicado, por ejemplo,
`#response.header.X-Trace-Id` obtendrá los valores en lista del campo `X-Trace-Id` del encabezado de respuesta si
existe, sustituyendo la sintaxis por el valor, el resultado fue `["4ae6c92d16089e521626", "5ae7c93d16089e521628"]`.

`#response.body...`

Este fragmento de la sintaxis obtendrá del cuerpo de respuesta del endpoint el valor indicado, por ejemplo,
`#response.body.role` obtendrá el valor del campo `role` del cuerpo de la respuesta si existe,
sustituyendo la sintaxis por el valor, el resultado fue `ADMIN`.

#### Historial temporal de respuestas

Si deseas obtener los valores de la respuesta específicamente de un backend, simplemente añade el campo `history` en la
sintaxis, y la posición de tu backend en el historial, recuerda que los middlewares `beforewares` y `afterwares` cuentan
como un backend en el historial, otro punto es que solo tendrás la información de la respuesta en este campo si el
backend ya ha sido procesado. Mira cómo se vería un ejemplo a continuación:

En el ejemplo, tengo solo un backend y ya ha sido procesado, entonces puedo estar utilizando la sintaxis:

`#response.history.0.header.X-Value.0`

En este otro ejemplo de sintaxis tenemos tres backends configurados y dos ya han sido procesados, entonces podemos
utilizar la siguiente sintaxis en el proceso del tercer backend:

`#request.history.1.body.users.0`

En estos ejemplos citados vemos que podemos obtener el valor de la respuesta de un backend que ya fue procesado,
y que están almacenados en un tipo de historial temporal.

### Importante

Puedes utilizar con base en estos campos,
la [sintaxis de JSON path](https://github.com/tidwall/gjson/blob/master/README.md#path-syntax) que se ajusta a tus
valores, solo recuerda que, los objetos header, query son mapas de lista de string, y params es un mapa de string.

### Véase también

- Configuración de modificadores de backend [clicando aquí](#backendmodifiers).

- Aprende en la práctica cómo utilizar los valores dinámicos para modificaciones usando el
  proyecto [playground](https://github.com/tech4works/gopen-gateway-playground) que ya viene con algunos ejemplos
  de modificadores con valores dinámicos.

## Observabilidad

Cuando la API Gateway recibe una solicitud, añadirá la IP del cliente en el encabezado `X-Forwarded-For` de las
solicitudes subyacentes, y también añade, si no existe, un valor único generado al encabezado llamado `X-Trace-Id` para
una mejor observabilidad de la solicitud recibida tanto en los logs de la API Gateway como en los microservices
subyacentes.

## Lógica de respuesta

Cuando utilizamos una API Gateway nos preguntamos, ¿cómo se devolverá a mi cliente
la respuesta de ese endpoint configurado?

Para facilitar la comprensión, creamos este tema para resumir la lógica de respuesta de nuestra API Gateway,
así que empecemos.

### ¿Cómo funciona?

La API Gateway fue desarrollada con una inteligencia y flexibilidad al responder un endpoint, se basa en dos
puntos importantes, primero, en la cantidad de respuestas de servicios backend que han sido procesados, y segundo, en
los campos de personalización de la respuesta configurados en los objetos [endpoint](#endpointcomment)
y [backend](#backendname). Vamos a ver algunos ejemplos a continuación para una mejor comprensión.

#### Único backend

En este ejemplo trabajaremos solo con un único backend, mira cómo se comportará la API Gateway al responder
a este escenario:

Json de configuración

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

Al procesar este endpoint la respuesta de la API Gateway fue:

```text
HTTP/1.1 200 OK
```

Encabezado ([Ver sobre los encabezados de respuesta aquí](#encabezado-de-respuesta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Tue, 23 Apr 2024 11:37:26 GMT
Content-Length: 620
```

Cuerpo

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

Vimos que en el ejemplo, la API Gateway actuó como un proxy redirigiendo la solicitud al servicio de backend configurado
y reflejando su cuerpo de respuesta, agregando además sus valores al encabezado de la respuesta.

En este mismo ejemplo, vamos a forzar un escenario de infelicidad en la respuesta del backend, vea:

```text
HTTP/1.1 404 Not Found
```

Encabezado ([Ver sobre los encabezados de respuesta aquí](#encabezado-de-respuesta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: false
Date: Tue, 23 Apr 2024 21:56:33 GMT
Content-Length: 235
```

Cuerpo

```json
{
  "file": "datastore/user.go",
  "line": 227,
  "endpoint": "/users/find/gabrielcataldo.adma@gmail.com",
  "message": "user not found"
}
```

En este caso, la API Gateway también reflejó la respuesta de la única llamada al backend del endpoint.

#### Utilizando middlewares

En este ejemplo, vamos a utilizar los middlewares de [beforewares](#endpointbeforewares)
y [afterwares](#endpointafterwares). Como estos backends son omitidos al cliente final si tienen éxito, vamos a simular
una llamada con el dispositivo bloqueado para que el [beforeware](#endpointbeforewares) devuelva un error, y luego
un [afterware](#endpointafterwares) que también responderá con un error, ya que no existe. ¡Vamos allá!

Json de configuración

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

Al procesar este endpoint de ejemplo simulando el error en la llamada a [beforeware](#endpointbeforewares), la respuesta
del API Gateway fue

```text
HTTP/1.1 403 Forbidden
```

Encabezado ([Ver sobre los encabezados de respuesta aquí](#encabezado-de-respuesta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: true
Date: Tue, 23 Apr 2024 23:02:09 GMT
Content-Length: 154
```

Cuerpo

```json
{
  "file": "service/device.go",
  "line": 49,
  "endpoint": "/devices",
  "message": "unprocessed entity: device already exists and is not active"
}
```

Vimos que la respuesta fue el espejo del retorno del beforeware `save-device`, ya que como el mismo retornó
un error `403 (Forbidden)`, el endpoint se abortó, no llamando a los backends siguientes. Recuerda que puedes
configurar los códigos de estado HTTP que se abortarán en tu endpoint, solo tienes que llenar el
campo [endpoint.abort-if-status-codes](#endpointabort-if-status-codes).

En el siguiente ejemplo, vamos a forzar un error en el afterware `increment-attempts` y la respuesta de la API Gateway
fue

```text
HTTP/1.1 404 Not Found
```

Encabezado ([Ver sobre los encabezados de respuesta aquí](#encabezado-de-respuesta))

```text
Content-Type: text/plain
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: false
Date: Tue, 23 Apr 2024 23:16:57 GMT
Content-Length: 18
```

Cuerpo

```text
404 page not found
```

Vimos que la respuesta también fue el espejo del retorno del afterware `increment-attempts`, aunque sea la última
llamada de un servicio backend del endpoint, porque cayó en la regla de respuesta abortada. Por lo tanto, todas las
otras respuestas de los otros backends fueron ignoradas y solo se devolvió la respuesta del backend abortado.

Vea más sobre la [respuesta abortada](#respuesta-abortada).

#### Múltiples backends

En este ejemplo trabajaremos con tres [backends](#endpointbackends) principales en el endpoint, ¡vamos allá!

Json de configuración

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

En el ejemplo ejecutaremos los tres backend con éxito, la API Gateway respondió

```text
HTTP/1.1 200 OK
```

Encabezado ([Ver sobre los encabezados de respuesta aquí](#encabezado-de-respuesta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Tue, 23 Apr 2024 23:49:12 GMT
Content-Length: 755
```

Cuerpo

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

Tenemos algunos puntos en este ejemplo que vale la pena resaltar. En primer lugar, con formato, la API Gateway entendió
que su endpoint tiene múltiples respuestas y no se utilizó el
campo [endpoint.aggregate-responses](#endpointaggregate-responses) con el valor `true`, por lo que lista las respuestas
como JSON añadiendo los siguientes campos:

`ok`: Indica si la respuesta del backend en cuestión tuvo el código de estado HTTP entre `200` y `299`.

`code`: Se llena con el código de estado HTTP respondido por su backend.

Estos campos solo se agregan si hay múltiples respuestas y el
campo [endpoint.aggregate-responses](#endpointaggregate-responses) no se proporciona con el valor `true`.

El segundo punto a destacar es en la parte `"version": "v1.0.0"` del último backend, el mismo respondió solo un texto en
el cuerpo de respuesta que fue `v1.0.0`, sin embargo, para este escenario como se mencionó, la API Gateway fuerza el
análisis de este valor para un JSON, añadiendo un nuevo campo con el nombre obtenido en la
configuración [backend.name](#backendname) y con el valor del mismo.

El tercer punto es sobre el código de estado HTTP, el mismo se devuelve por la mayor frecuencia, es decir, si tenemos
tres devoluciones `200 OK` como en el ejemplo, la API Gateway también devolverá ese código. Si tenemos una devolución
igualitaria, el último código de estado HTTP devuelto será considerado, vea los escenarios posibles de esta lógica:

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

la API Gateway responderá `201 Created`.

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

la API Gateway responderá `100 Continue`.

El cuarto punto a destacar es que, como el endpoint tiene múltiples respuestas, consecuentemente tenemos múltiples
encabezados de respuesta, la API Gateway agregará todos los campos y valores para el encabezado de la respuesta final,
vea más sobre el comportamiento del encabezado de respuesta haciendo [clic aquí](#encabezado-de-respuesta).

El último punto a destacar es que si una de estas devoluciones de backend entra en el escenario en que el endpoint
aborta la respuesta, no seguirá ninguna directriz mostrada en el tema en cuestión, sino
[lógica de respuesta abortada](#respuesta-abortada).

#### Múltiples backends agregados

En este ejemplo utilizaremos una configuración parecida con el JSON de configuración del ejemplo anterior, pero con
campo [endpoint](#endpointaggregate-responses) con el valor `true`.

JSON de la configuración

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

Cuando procesamos el endpoint, la respuesta de la API Gateway fue:

```text
HTTP/1.1 200 OK
```

Encabezado ([Ver sobre los encabezados de respuesta aquí](#encabezado-de-respuesta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Wed, 24 Apr 2024 10:57:31 GMT
Content-Length: 665
```

Cuerpo

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

Vimos que la única diferencia en la respuesta del tópico [Múltiples backends](#múltiples-backends) es que agregó los
valores de todas las respuestas en un solo JSON, y los campos que se repitieron han agregado los valores en una lista.

Las demás reglas como el código de estado HTTP, la conversión forzada a JSON, entre otras, siguen las mismas reglas
mencionadas en el tópico [Múltiples backends](#múltiples-backends).

En el ejemplo, podemos hacer que la respuesta agregada sea un poco más organizada, para ello vamos a modificar la
sección de nuestro segundo backend agregando el
campo [backend.extra-config.group-response](#backendextra-configgroup-response) con el valor `true` y darle un nombre,
vea la sección del JSON de configuración modificada:

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

Cuando procesamos nuevamente el endpoint, obtuvimos la siguiente respuesta:

```text
HTTP/1.1 200 OK
```

Encabezado ([Ver sobre los encabezados de respuesta aquí](#encabezado-de-respuesta))

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: true
X-Gopen-Success: true
Date: Wed, 24 Apr 2024 11:23:07 GMT
Content-Length: 697
```

Cuerpo

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

Con esta configuración, vimos que nuestra respuesta agregada se organizó mejor, y es importante entender sobre
el [json de configuración](#json-de-configuración) y sus campos, para que el GOPEN API Gateway pueda satisfacer mejor
sus necesidades.

### Respuesta abortada

Para que una respuesta sea abortada por la API Gateway, uno de los backends configurados en el endpoint, tanto
middlewares como los principales, al ser procesados, en su respuesta, el código de estado HTTP necesita seguir los
valores en el campo [endpoint.abort-if-status-codes](#endpointabort-if-status-codes) del propio endpoint.

Es importante que entiendas estos dos puntos:

- Al abortar la respuesta del backend, la API Gateway solo reflejará la respuesta del mismo, código de estado,
  encabezado y cuerpo. Por lo tanto, se ignorarán las otras respuestas ya procesadas.
- Incluso si el backend es abortado, los [modificadores](#backendmodifiers) serán procesados.

Recomendamos utilizar esta configuración solo cuando algo no salió como se esperaba, como, por ejemplo, una respuesta de
error `500 (Internal server error)`.

### Encabezado de respuesta

En la respuesta, la API Gateway, con la excepción de los campos `Content-Length`, `Content-Type` y `Date`, agrega todos
los valores del encabezado respondidos por los backends configurados en el endpoint, independientemente de la cantidad
de backends, esto incluye también a los [middlewares](#middlewares).

#### Campos de encabezado estándar

También se agregan hasta cuatro campos en el encabezado. A continuación, te explicamos cada uno de ellos:

`X-Gopen-Cache`: Si la respuesta del endpoint no es "fresca", es decir, se utilizó la respuesta almacenada en caché,
se devuelve el valor `true`, de lo contrario se devuelve el valor `false`.

`X-Gopen-Cache-Ttl`: Si la respuesta del endpoint ha sido realizada utilizando el almacenamiento en caché, devuelve la
duración del tiempo de vida restante de esa caché, de lo contrario, el campo no se devuelve.

`X-Gopen-Complete`: Si todos los backends han sido procesados por el endpoint, se devuelve el valor `true`, en caso
contrario se devuelve el valor `false`.

`X-Gopen-Success`: Si todos los backends han devuelto éxitos, es decir, el código de estado HTTP de respuesta entre
`200` y `299`, devuelve el valor `true`, de lo contrario el valor `false`.

Recuerda que si la respuesta de un backend es [abortada](#respuesta-abortada), solo se agrega el header del mismo y
se consideran las reglas de los campos anteriores.

Ahora veamos algunos ejemplos de encabezados de retorno:

#### Campos únicos de encabezado

Encabezado de respuesta del backend 1:

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521626
X-MS: api-user
Date: Wed, 24 Apr 2024 11:23:07 GMT
Content-Length: 102
```

Encabezado de respuesta del endpoint:

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

Vemos en el ejemplo que se han agregado los [campos estándar](#campos-de-encabezado-estándar) y se agregaron los valores
del encabezado de respuesta, que fueron `X-Value-Id` y `X-MS`.

#### Campos de encabezado duplicados

Encabezado de respuesta del backend 1:

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521626
X-MS: api-user
Date: Wed, 24 Apr 2024 11:23:07 GMT
Content-Length: 102
```

Encabezado de respuesta del backend 2:

```text
Content-Type: application/json
X-Value-Id: 4ae6c92d16089e521638
X-MS: api-device
X-MS-Success: true
Date: Wed, 24 Apr 2024 11:23:08 GMT
Content-Length: 402
```

Encabezado de respuesta del endpoint:

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

Vemos en el ejemplo que también se han agregado los [campos estándar](#campos-de-encabezado-estándar) y se agregaron los
valores del encabezado de respuesta, que fueron `X-Value-Id`, `X-MS` y `X-MS-Success`. Es importante destacar que los
campos duplicados se agruparon y se separaron por una coma.

### Respuestas estándar

Cada API Gateway tiene sus respuestas estándar para cada escenario de error, entonces enumeraremos
cada escenario y su respectiva respuesta HTTP a continuación:

#### 413 (Request Entity Too Large)

Este escenario ocurre cuando el tamaño del cuerpo de la solicitud es mayor que el permitido para el punto final,
utilizando la configuración [limiter.max-body-size](#limitermax-header-size) para el cuerpo normal
y [limiter.max-multipart-memory-size](#limitermax-multipart-memory-size) para el envío del tipo `form-data`. Puedes
personalizar esta configuración para un punto final específico utilizando el campo [endpoint.limiter](#endpointlimiter).

Encabezado

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 11:56:06 GMT
Content-Length: 170
```

Cuerpo

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

Este escenario sucede cuando se alcanza el límite de solicitudes desde una IP específica, este límite se define en la
configuración [limiter.rate](#limiterrate). Puedes personalizar esta configuración para un punto final específico
utilizando el campo [endpoint.limiter](#endpointlimiter).

Encabezado

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 12:12:53 GMT
Content-Length: 177
```

Cuerpo

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

Este escenario sucede cuando el tamaño del encabezado es mayor que el permitido para el endpoint, utilizando la
configuración [limiter.max-header-size](#limitermax-header-size). Puedes personalizar esta configuración para un punto
final específico utilizando el campo [endpoint.limiter](#endpointlimiter).

Encabezado

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 11:39:53 GMT
Content-Length: 186
```

Cuerpo

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

Este escenario es específico cuando ocurra algún error inesperado con la API Gateway. En caso de que suceda esto, por
favor informe el problema [aquí](https://github.com/tech4works/gopen-gateway/issues), mostrando la respuesta y el
registro de impresiones en el terminal de ejecución.

Encabezado

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 12:38:16 GMT
Content-Length: 183
```

Cuerpo

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

Este escenario ocurre cuando se intenta comunicar con el backend y existe algún fallo en la comunicación con este.

Encabezado

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Thu, 25 Apr 2024 01:07:36 GMT
Content-Length: 277
```

Cuerpo

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

Este escenario sucede cuando el punto final excede el límite de tiempo configurado en el campo [timeout](#timeout).
Puedes personalizar esta configuración para un punto final específico utilizando el
campo [endpoint.timeout](#endpointtimeout).

Encabezado

```text
Content-Type: application/json
X-Gopen-Cache: false
X-Gopen-Complete: false
X-Gopen-Success: false
Date: Fri, 26 Apr 2024 13:29:55 GMT
Content-Length: 150
```

Cuerpo

```json
{
  "file": "middleware/timeout.go",
  "line": 81,
  "endpoint": "/users/version",
  "message": "gateway timeout: 5m",
  "timestamp": "2024-04-26T10:29:55.908526-03:00"
}
```

# Usabilidad

¡Utiliza el proyecto [playground](https://github.com/tech4works/gopen-gateway-playground) para comenzar a explorar
y utilizar en la práctica la API Gateway de GOPEN!

# ¿Cómo contribuir?

Nos alegra cuando vemos que la comunidad se apoya mutuamente, y proyectos como este, están abiertos para recibir
tus ideas. A continuación, te presentamos cómo puedes participar.

## Download

Para poder ejecutar el proyecto, primero debes descargar el [lenguaje Go](https://go.dev/dl/) versión 1.22 o superior en
tu máquina.

Con Go instalado en tu máquina, realiza el pull del proyecto:

```text
git pull https://github.com/tech4works/gopen-gateway.git
```

Luego, abre el proyecto utilizando la terminal con el IDE de tu preferencia:

Goland:

```text
goland gopen-gateway
```

VSCode:

```text
code gopen-gateway
```

## Air

Si quieres utilizar la recarga en caliente para optimizar durante el desarrollo, el proyecto ya tiene el archivo
`.air.toml` como ejemplo para que lo uses como quieras.

Mira la documentación sobre este fantástico proyecto haciendo [clic aquí](https://github.com/cosmtrek/air).

## Gitflow

Para iniciar el desarrollo, puedes crear una rama a partir de `main`, para un futuro
Pull Request hacia la misma.

# Agradecimientos

Este proyecto tuvo el apoyo de bibliotecas fantásticas, este segmento está dedicado a cada una de ellas listadas
a continuación:

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

Gracias por contribuir a la comunidad Go y facilitar el desarrollo de este proyecto.

# License Apache 2.0

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2FGabrielHCataldo%2Fgopen-gateway?ref=badge_large&issueType=license)

