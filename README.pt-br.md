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
- Configuração de cache global e local para cada endpoint, com customização da estratégia da chave de armazenamento, e
  condições baseada em códigos de status de resposta e método HTTP para ler e salvar o mesmo.
- Armazenamento de cache local ou global utilizando Redis.
- Configuração de limitador de tamanho, global e local para cada endpoint, limitando o tamanho do Header, Body e
  Multipart
  Memory.
- Configuração de limitador de taxa, global e local para cada endpoint, limitando pelo tempo e rajada pelo IP.
- Configuração de segurança de CORS com validações de origens, método HTTP e headers.
- Configuração global de múltiplos middlewares, para serem usados posteriormente no endpoint caso indicado.
- Filtragem personalizada de envio de headers e query para os backends do endpoint.
- Processamento de múltiplos backends, sendo eles beforewares, principais e afterwares para o endpoint.
- Configuração personalizada para abortar processo de execução dos backends pelo código de status retornado.
- Modificadores para todos os conteúdos de requisição e response (Status Code, Path, Header, Params, Query, Body)
  ao nível global (requisição/response) e local (requisição backend/response backend) com ações de remoção,
  adição, alteração, substituição e renomeio.
- Obtenha o valor a ser modificado de variáveis de ambiente, da requisição atual, do histórico de respostas do endpoint,
  ou até mesmo do valor passado na configuração.
- Executa os modificadores no contexto que desejar, antes de uma requisição backend ou depois, você decide.
- Faça com que as modificações reflitam em todas as requisições/respostas seguintes, usando a mesma ao nível global.
- Omita a resposta de um backend caso necessite, a mesma não será utilizada na resposta do endpoint.
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

### ESTRUTURA DE PASTAS

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

--- 

### JSON DE CONFIGURAÇÃO

Com base nesse arquivo json de configuração obtido pela env desejada a aplicação terá seus endpoints e suas regras
definidas, veja abaixo um exemplo simples com todos os campos possíveis e seus conceitos e regras:

````json
{
  "$schema": "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json",
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
      "name": "Save device",
      "hosts": [
        "http://192.168.1.2:8051"
      ],
      "path": "/devices",
      "method": "PUT",
      "forward-headers": [],
      "modifiers": {
        "header": [
          {
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
      "@comment": "Feature: Find user by key",
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
      "beforeware": [
        "save-device"
      ],
      "backends": [
        {
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

Campo obrigatório, para o auxílio na escrita e regras do próprio json de configuração.

### version

Campo opcional, usado para retorno do endpoint estático `/version`.

### port

Campo obrigatório, utilizado para indicar a porta a ser ouvida pela API Gateway, valor mínimo `1` e valor
máximo `65535`.

### hot-reload

Campo opcional, o valor padrão é `false`, caso seja `true` é utilizado para o carregamento automático quando
houver alguma alteração no arquivo .json e .env na pasta do ambiente selecionado.

### store

Campo opcional, valor padrão é o armazenamento local em cache, caso seja informado, o campo `redis` passa
a ser obrigatório e o campo `address` também.

Caso utilize o armazenamento global de cache o Redis, é indicado que os valores de endereço e senha sejam preenchidos
utilizando variável de ambiente, como no exemplo acima.

### timeout

Campo opcional, o valor padrão é `30 segundos`, esse campo é responsável pelo tempo máximo de duração do processamento
de cada requisição.

Caso a requisição ultrapasse esse tempo informado, á API Gateway irá abortar todas as transações em andamento e
retornará o código de status `504 (Gateway Timeout)`.

IMPORTANTE: Caso seja informado no objeto de endpoint, damos prioridade ao valor informado do endpoint, caso contrário
seguiremos com o valor informado ou padrão desse campo, na raiz do json de configuração.

```
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
```

### cache

Campo opcional, se informado, o campo `duration` passa a ser obrigatório!

Caso o objeto seja informado na estrutura do [endpoint.cache](#cache), damos prioridade aos valores informados lá,
caso contrário seguiremos com os valores informados nesse campo.

O valor do cache é apenas gravado 1 vez a cada X duração informada no campo `every`.

Os campos `only-if-status-codes` e `only-if-methods` são utilizados para verificar se naquele endpoint habilitado
a ter cache, pode ser lido e escrito o cache com base no método HTTP e código de status de resposta, veja mais sobre
eles abaixo.

Caso a resposta não seja "fresca", ou seja, foi respondida pelo cache, o header `X-Gopen-Cache` terá o valor `true`
caso contrário o valor será `false`.

#### cache.duration

Indica o tempo que o cache irá durar, ele é do tipo `time.Duration`.

```
- Valores aceitos:
    - s para segundos
    - m para minutos
    - h para horas
    - ms para milissegundos
    - us (ou µs) para microssegundos
    - ns para nanossegundos

- Exemplos:
    - 1h
    - 15.5ms
    - 1h30m
    - 1.5m
```

#### cache.strategy-headers

Campo opcional, a estrátegia padrão de chave de cache é pela url e método HTTP da requisição tornando-o um cache global
por endpoint, caso informado os cabeçalhos a serem usados na estrátegia eles são agregados nos valores padrões de chave,
por exemplo, ali no exemplo foi indicado utilizar o campo `X-Forwarded-For` e o `Device` o valor final da chave
ficaria:

      GET:/users/find/479976139:177.130.228.66:95D4AF55-733D-46D7-86B9-7EF7D6634EBC

A descrição da lógica por trás dessa chave é:

      método:url:X-Forwarded-For:Device

Sem a estrátegia preenchida, a lógica padrão fica assim:

      método:url

Então o valor padrão para esse endpoint fica assim sem a estrátegia preenchida:

      GET:/users/find/479976139

Nesse exemplo tornamos o cache antes global para o endpoint em espécifico, passa a ser por cliente! Lembrando que isso
é um exemplo simples, você pode ter a estrátegia que quiser com base no header de sua aplicação.

#### cache.only-if-methods

Campo opcional, o valor padrão é uma lista com apenas o método HTTP `GET`, caso informada vazia, qualquer método HTTP
será aceito.

Esse campo é responsável por decidir se irá ler e gravar o cache do endpoint (que está habilitado a ter cache) pelo
método HTTP do mesmo.

#### cache.only-if-status-codes

Campo opcional, o valor padrão é uma lista de códigos de status HTTP de sucessos reconhecidos, caso informada vazia,
qualquer código de status HTTP de resposta será aceito.

Esse campo é responsável por decidir se irá gravar o cache do endpoint (que está habilitado a ter cache) pelo
código de status HTTP de resposta do mesmo.

#### cache.allow-cache-control

Campo opcional, o valor padrão é `false`, caso seja informado como `true` a API Gateway irá considerar o header
`Cache-Control` seguindo as regras a seguir a partir do valor informado na requisição ou na resposta dos backends:

`no-cache`: esse valor é apenas considerado no header da requisição, caso informado desconsideramos a leitura do cache
e seguimos com o processo normal para obter a resposta "fresca".

`no-store`: esse valor é considerado apenas na resposta escrita por seus backends, caso informado não gravamos o
cache.

### limiter

Campo opcional, objeto responsável pelas regras de limitação da API Gateway, seja de tamanho ou taxa, os valores padrões
variam de campo a campo, veja:

#### limiter.max-header-size

Campo opcional, ele é do tipo `byteUnit`, valor padrão é `1MB`, é responsável por limitar o tamanho do cabeçalho de
requisição.

Caso o tamanho do cabeçalho ultrapasse o valor informado, a API Gateway irá abortar a requisição com o código de status
`431 (Request header fields too large)`.

```
- Valores aceitos:
    - B para Byte
    - KB para KiloByte
    - MB para Megabyte
    - GB para Gigabyte
    - TB para Terabyte
    - PB para Petabyte
    - EB para Exabyte
    - ZB para Zettabyte
    - YB para Yottabyte

- Exemplos:
    - 1B
    - 50KB
    - 5MB
    - 1.5GB
```

#### limiter.max-body-size

Campo opcional, ele é do tipo `byteUnit`, valor padrão é `3MB`, campo é responsável por limitar o tamanho do corpo
da requisição.

Caso o tamanho do corpo ultrapasse o valor informado, a API Gateway irá abortar a requisição com o código de status
`413 (Request entity too large)`.

```
- Valores aceitos:
    - B para Byte
    - KB para KiloByte
    - MB para Megabyte
    - GB para Gigabyte
    - TB para Terabyte
    - PB para Petabyte
    - EB para Exabyte
    - ZB para Zettabyte
    - YB para Yottabyte

- Exemplos:
    - 1B
    - 50KB
    - 5MB
    - 1.5GB
```

#### limiter.max-multipart-memory-size

Campo opcional, ele é do tipo `byteUnit`, valor padrão é `5MB`, esse campo é responsável por limitar o tamanho do
corpo multipart/form da requisição, geralmente utilizado para envio de arquivos, imagens, etc.

Caso o tamanho do corpo ultrapasse o valor informado, a API Gateway irá abortar a requisição com o código de status
`413 (Request entity too large)`.

```
- Valores aceitos:
  - B para Byte
  - KB para KiloByte
  - MB para Megabyte
  - GB para Gigabyte
  - TB para Terabyte
  - PB para Petabyte
  - EB para Exabyte
  - ZB para Zettabyte
  - YB para Yottabyte

- Exemplos:
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
`429 (Too many requests)`.

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
como `beforeware` e `afterware`.

O valor da chave é um objeto de [backend](#backendname), porém, com uma observação, esse objeto terá
sua resposta de sucesso omitida automáticamente pelo endpoint, já que respostas de sucesso de middlewares não são
exibidas para o cliente final HTTP, porém sua resposta será armazenada ao longo da requisição HTTP feita no endpoint,
podendo ser manipulada.

Por exemplo, um `beforeware` quando mencionado no endpoint, ele será utilizado como middleware de pré-requisições, isto
é, ele será chamado antes dos backends principais do endpoint, então podemos, por exemplo, ter um middleware
de manipulação de device, como no json de configuração acima, aonde ele irá chamar esse backend de middleware
configurado no endpoint como `beforeware`, validando e salvando o dispositivo a partir de informações do header da
requisição, caso o backend responda um código de status de falha, no exemplo, o gateway abortará todos os backends
seguintes retornando o que o backend de device respondeu, caso tenha retornado um código de status de sucesso, ele irá
modificar o header de todas as requisições seguintes (`propagate:true`), adicionando o campo `X-Device-Id`, com o valor
do id do body de resposta do próprio backend, podendo ser utilizado nos outros backends seguintes do endpoint.

Para saber mais sobre os `modifiers` [veja](#backendmodifiers).

Para entender melhor essa ferramenta poderosa, na prática, veja os exemplos de middlewares usados como
`beforeware` e `afterware` feitos no projeto
de [playground](https://github.com/GabrielHCataldo/gopen-gateway-playground).

### endpoints

Campo obrigatório, é uma lista de objeto, representa cada endpoint da API Gateway que será registrado para ouvir e
servir as requisições HTTP.

Veja abaixo como funciona o fluxo básico de um endpoint na imagem abaixo:

#### TODO: colocar imagem

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

```
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

Campo opcional, do tipo lista de inteiros, o valor padrão é vazio, indicando que qualquer backend executado no endpoint
que tenha respondido o status code maior ou igual a `400 (Bad request)` será abortado.

Caso informado, e um backend retorna o status code indicado na configuração, o endpoint será abortado, isso significa
que os outros backends configurados após o mesmo, não serão executados, e o endpoint irá retornar a resposta do mesmo
ao cliente final.

Veja como o endpoint será respondido após um backend ser abortado clicando [aqui](#lógica-de-resposta).

#### endpoint.beforeware

Campo opcional, do tipo lista de string, o valor padrão é vazio, indicando que o endpoint não tem nenhum middleware
de pré-requisições.

Caso informado, o endpoint irá executar as requisições, posição por posição, começando no início da lista. Caso o valor
em string da posição a ser executada estiver configurada no campo [middlewares](#middlewares) corretamente, será
executado
o backend configurado no mesmo. Caso contrário irá ignorar a posição apenas imprimindo um log de atenção.

#### endpoint.afterware

Campo opcional, do tipo lista de string, o valor padrão é vazio, indicando que o endpoint não tem nenhum middleware
de pós-requisições.

Caso informado, o endpoint irá executar as requisições, posição por posição, começando no início da lista. Caso o valor
em string da posição a ser executada estiver configurada no campo [middlewares](#middlewares) corretamente, será
executado o backend configurado no mesmo. Caso contrário irá ignorar a posição apenas imprimindo um log de atenção.

#### endpoint.backends

Campo obrigatório, do tipo lista de objeto, responsável pela execução principal do endpoint, o próprio nome já diz tudo,
é uma lista que indica todos os serviços necessários para que o endpoint retorne a resposta esperada.

Veja abaixo como funciona o fluxo básico de um backend na imagem abaixo:

#### TODO: colocar imagem

Abaixo iremos listar e explicar cada campo desse objeto tão importante:

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

Utilizamos um dos [backend.hosts](#backendhosts) informados e juntamos com o path fornecido, por exemplo, no campo hosts
temos o valor

```
[
  "https://instance-01", 
  "https://instance-02"
]
```

E nesse campo path temos o valor

```
/users/status
```

O backend irá construir a seguinte URL de requisição

```
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
requisição com o seguinte cabeçalho

````
{
  "Device": "95D4AF55-733D-46D7-86B9-7EF7D6634EBC",
  "User-Agent": "IOS",
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
````

Nesse exemplo temos o campo `forward-headers` com os seguintes valores

````
[
  "User-Agent",
  "Authorization"
]
````

O cabeçalho de requisição ao backend foi

```
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
ser agrupada em um campo json para a resposta ao cliente final.

Caso informado com o valor `true` o body de resposta caso tenha, será agrupado em um campo json da resposta final,
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

Veja abaixo como funciona o fluxo básico de um modificador na imagem abaixo:

#### TODO: colocar imagem

Abaixo iremos listar e explicar cada campo desse objeto tão importante:

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

#### header.context

Campo obrigatório, do tipo string, é responsável por indicar qual contexto a modificação deve atuar.

Valores aceitos:

`REQUEST` para atuar na pré-requisição do backend.

`RESPONSE` para atuar pós-requisição do backend.

Importante lembrar que caso o valor for `REQUEST` poderá utilizar no campo [header.scope](#headerscope) apenas o valor
`REQUEST`.

#### header.scope

Campo opcional, do tipo string, o valor padrão será baseado no campo [header.context](#headercontext) informado, o valor
indica qual escopo devemos alterar, se o escopo de requisição ou de resposta.

Valores aceitos:

`REQUEST` para modificar o escopo de requisição, esse tipo de escopo pode ter uma atuação global propagando essa mudança
nas requisições backends seguintes, basta utilizar o campo [header.propagate](#headerpropagate) como `true`.

`RESPONSE` para modificar o escopo de resposta do backend.

#### header.action

Campo obrigatório, do tipo string, responsável pela ação a ser tomada na modificação do cabeçalho.

Valores aceitos:

`ADD` adiciona a chave informada no campo [header.key](#headerkey) caso não exista, e agrega o valor informado no
campo [header.value](#headervalue).

`SET` modifica o valor da chave informada no campo [header.key](#headerkey) pelo valor passado no
campo [header.value](#headervalue).

`DEL` remove a chave informada no campo [header.key](#headerkey).

`REN` renomeia a chave informada no campo [header.key](#headerkey) pelo valor passado no
campo [header.value](#headervalue).

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

#### param.context

Campo obrigatório, do tipo string, é responsável por indicar qual contexto a modificação deve atuar.

Valores aceitos:

`REQUEST` para atuar na pré-requisição do backend.

`RESPONSE` para atuar pós-requisição do backend.

#### param.action

Campo obrigatório, do tipo string, responsável pela ação a ser tomada na modificação dos parâmetros da requisição.

Valores aceitos:

`SET` modifica o valor da chave informada no campo [param.key](#paramkey) pelo valor passado no
campo [param.value](#paramvalue).

`DEL` remove a chave informada no campo [param.key](#paramkey).

`REN` renomeia a chave informada no campo [param.key](#paramkey) pelo valor passado no
campo [param.value](#paramvalue).

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

#### query.context

Campo obrigatório, do tipo string, é responsável por indicar qual contexto a modificação deve atuar.

Valores aceitos:

`REQUEST` para atuar na pré-requisição do backend.

`RESPONSE` para atuar pós-requisição do backend.

#### query.action

Campo obrigatório, do tipo string, responsável pela ação a ser tomada na modificação dos parâmetros de busca da
requisição.

Valores aceitos:

`ADD` adiciona a chave informada no campo [query.key](#querykey) caso não exista, e agrega o valor informado no
campo [query.value](#queryvalue).

`SET` modifica o valor da chave informada no campo [query.key](#querykey) pelo valor passado no
campo [query.value](#queryvalue).

`DEL` remove a chave informada no campo [query.key](#querykey).

`REN` renomeia a chave informada no campo [query.key](#querykey) pelo valor passado no
campo [query.value](#queryvalue).

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

#### body.context

Campo obrigatório, do tipo string, é responsável por indicar qual contexto a modificação deve atuar.

Valores aceitos:

`REQUEST` para atuar na pré-requisição do backend.

`RESPONSE` para atuar pós-requisição do backend.

Importante lembrar que caso o valor for `REQUEST` poderá utilizar no campo [body.scope](#bodyscope) apenas o valor
`REQUEST`.

#### body.scope

Campo opcional, do tipo string, o valor padrão será baseado no campo [body.context](#bodycontext) informado, o valor
indica qual escopo devemos alterar, se o escopo de requisição ou de resposta.

Valores aceitos:

`REQUEST` para modificar o escopo de requisição, esse tipo de escopo pode ter uma atuação global propagando essa mudança
nas requisições de backend seguintes, basta utilizar o campo [body.propagate](#bodypropagate) como `true`.

`RESPONSE` para modificar o escopo de resposta do backend.

#### body.action

Campo obrigatório, do tipo string, responsável pela ação a ser tomada na modificação do body.

Valores aceitos se o body for JSON:

`SET` modifica o valor da chave informada no campo [body.key](#bodykey) pelo valor passado no
campo [body.value](#bodyvalue).

`REN` renomeia a chave informada no campo [body.key](#bodykey) pelo valor passado no
campo [body.value](#bodyvalue).

`DEL` remove a chave informada no campo [body.key](#bodykey).

Valores aceitos se o body for TEXTO:

`ADD` agrega o valor informado no campo [body.value](#bodyvalue) ao texto.

`SET` irá substituir todos os valores semelhantes à chave informada no campo [body.key](#bodykey) pelo valor passado no
campo [body.value](#bodyvalue).

`DEL` remove todos os valores semelhantes à chave informada no campo [body.key](#bodykey).

#### body.key

Campo obrigatório, do tipo string, utilizado para indicar qual chave do cabeçalho deve ser modificada.

OBS: se torna opcional se seu body for do tipo TEXTO e [body.action](#bodyaction) tiver o valor `ADD`.

#### body.value

Campo obrigatório, do tipo string, utilizado como valor a ser usado para modificar a chave indicada no
campo [header.key](#headerkey).

Temos possibilidades de utilização de [valores dinâmicos](#valores-dinâmicos-para-modificação),
e de [variáveis de ambiente](#variáveis-de-ambiente) para esse campo.

OBS: se torna opcional apenas se [header.action](#headeraction) tiver o valor `DEL`.

#### body.propagate

Campo opcional, do tipo booleano, o valor padrão é `false` indicando que o modificador não deve propagar essa mudança
em questão para os backends seguintes.

Caso informado como `true` essa modificação será propagada para os seguintes backends.

IMPORTANTE: Esse campo só é aceito se o [escopo](#bodyscope) tiver o valor `REQUEST`.

---

### VARIÁVEIS DE AMBIENTE

---

### VALORES DINÂMICOS PARA MODIFICAÇÃO

---

### LÓGICA DE RESPOSTA

--- 

Usabilidade
-----------
---
Use o projeto [playground](https://github.com/GabrielHCataldo/gopen-gateway-playground) para começar a explorar e
utilizar na prática o Gopen API Gateway!


Como contríbuir?
------------
---


Agradecimentos
------------
---

