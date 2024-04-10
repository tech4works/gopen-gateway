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
  condições baseada em códigos de status de resposta e método http para salvar o mesmo.
- Armazenamento de cache local ou global utilizando Redis
- Configuração de limitador de tamanho global e local para cada endpoint, limitando o tamanho Header, Body e Multipart
  Memory.
- Configuração de limitador de taxa global e local para cada endpoint, limitando pelo tempo e rajada pelo IP.
- Configuração de segurança de CORS com validações de origens, método http e headers.
- Configuração global de múltiplos middlewares, para serem usados posteriormente no endpoint caso indicado.
- Filtragem personalizada de envio de headers e query para os backends do endpoint.
- Processamento de múltiplos backends, sendo eles beforewares, principais e afterwares para o endpoint.
- Configuração personalizada para abortar processo de execução dos backends pelo código de status retornado.
- Modificadores para todos os conteúdos da requisição e response (Status Code, Path, Header, Params, Query, Body)
  ao nível global (requisição) e local (atual requisição/response backend) com ações de remoção,
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
      "cache": {
        "enabled": true,
        "duration": "30s",
        "strategy-headers": [],
        "only-if-status-codes": [],
        "allow-cache-control": false
      },
      "method": "GET",
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
            "statusCode": {},
            "header": [],
            "params": [],
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

### timeout

Campo opcional, o valor padrão é `30 segundos`, esse campo é responsável pelo tempo máximo de duração do processamento
de cada requisição.

Caso a requisição ultrapasse esse tempo informado, á API Gateway irá abortar todas as transações em andamento e
retornará
o código de status `504 (Gateway Timeout)`.

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

### store

Campo opcional, valor padrão é o armazenamento local em cache, caso seja informado, o campo `redis` passa
a ser obrigatório e os outros dois campos que acompanham o mesmo `address` e `password` também.

Caso utilize o armazenamento global de cache o Redis, é indicado que os valores de endereço e senha sejam preenchidos
utilizando variável de ambiente, como no exemplo acima.

### cache

Campo opcional, se informado, o campo `duration` passa a ser obrigatório!

Caso o objeto seja informado na estrutura do endpoint, damos prioridade aos valores informados lá, caso contrário
seguiremos com os valores informados nesse campo.

O valor do cache é apenas gravado 1 vez a cada X duração informada.

Os campos `only-if-status-codes` e `only-if-methods` são utilizados para verificar se naquele endpoint habilitado
a ter cache, pode ser lido e escrito, veja mais sobre eles abaixo.

Caso a resposta não seja "fresca", ou seja, foi respondida pelo cache, o header `X-Gopen-Cache` terá o valor `true`
caso contrário o valor será `false`.

#### duration

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

#### strategy-headers

Campo opcional, a estrátegia padrão de chave de cache é pela url e método http da requisição tornando-o um cache global
por endpoint, caso informado os cabeçalhos a serem usados na estrátegia eles são agregados nos valores padrões de chave,
por exemplo, ali no exemplo foi indicado utilizar o campo `X-Forwarded-For` e o `Device` o valor final da chave
ficaria:

      GET:/users/find/479976139:177.130.228.66:95D4AF55-733D-46D7-86B9-7EF7D6634EBC

A descrição da lógica por trás dessa chave é:

      {método}:{url}:{X-Forwarded-For}:{Device}

Sem a estrátegia preenchida, a lógica padrão fica assim:

      {método}:{url}

Então o valor padrão para esse endpoint fica assim sem a estrátegia preenchida:

      GET:/users/find/479976139

Nesse exemplo tornamos o cache antes global para o endpoint em espécifico, passa a ser por cliente! Lembrando que isso
é um exemplo simples, você pode ter a estrátegia que quiser com base no header de sua aplicação.

#### only-if-methods

Campo opcional, o valor padrão é uma lista com apenas o método http `GET`, caso informada vazia, qualquer método http
será aceito.

Esse campo é responsável por decidir se irá ler e gravar o cache do endpoint (que está habilitado a ter cache) pelo
método http do mesmo.

#### only-if-status-codes

Campo opcional, o valor padrão é uma lista de códigos de status http de sucessos reconhecidos, caso informada vazia,
qualquer código de status http de resposta será aceito.

Esse campo é responsável por decidir se irá gravar o cache do endpoint (que está habilitado a ter cache) pelo
código de status http de resposta do mesmo.

#### allow-cache-control

Campo opcional, o valor padrão é `false`, caso seja informado como `true` a API Gateway irá considerar o header
`Cache-Control` seguindo as regras a seguir a partir do valor informado na requisição ou na resposta dos backends:

`no-cache`: esse valor é apenas considerado no header da requisição, caso informado desconsideramos a leitura do cache
e seguimos com o processo normal para obter a resposta "fresca".

`no-store`: esse valor é considerado apenas na resposta escrita por seus backends, caso informado não gravamos o
cache.

### limiter

Campo opcional, objeto responsável pelas regras de limitação da API Gateway, seja de tamanho ou taxa, os valores padrões
variam de campo a campo, veja:

#### max-header-size

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

#### max-body-size

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

#### max-multipart-memory-size

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

#### rate

Campo opcional, caso seja informado, o campo `capacity` torna-se obrigatório, esse objeto é responsável por limitar
a taxa de requisição pelo IP, esse limite é imposto obtendo a capacidade máxima pelo campo `capacity` por X duração,
informado no campo `every`.

Caso essa capacidade seja ultrapassada, a API Gateway por segurança abortará a requisição, retornando
`429 (Too many requests)`.

#### rate.capacity

Campo opcional, caso o objeto rate seja informado, ele passa a ser obrigatório, o valor padrão é `5`, e o mínimo
que poderá ser informado é `1`, indica a capacidade máxima de requisições.

#### rate.every

Campo opcional, o valor padrão é `1 segundo`, indica o valor da duração da verificação da capacidade máxima de
requisições.

### security-cors

Campo opcional, usado para segurança do CORS da API Gateway, todos os campos por padrão são vazios, não restringindo
os valores de origin, methods e headers.

Caso queira restringir, e a requisição não esteja de acordo com as configurações impostas, a API Gateway por segurança
irá abortar a requisição retornando `403 (Forbidden)`.

#### allow-origins

Campo opcional, do tipo lista de string, os itens da lista precisam indicar quais IPs de origem a API Gateway
permite receber nas requisições.

#### allow-methods

Campo opcional, do tipo lista de string, os itens da lista precisam indicar quais métodos http a API Gateway
permite receber nas requisições.

#### allow-headers

Campo opcional, do tipo lista de string, os itens da lista precisam indicar quais campos de cabeçalho http a API Gateway
permite receber nas requisições.

### middlewares

Campo opcional, é responsável pela configuração de seus middlewares de aplicação, é um mapa com chaves
em string mencionando o nome do seu middleware, esse nome poderá ser utilizado em seu [endpoint](#endpoint)
como `beforeware` e `afterware`.

O valor da chave é um objeto de [backend](#backend), porém, com uma observação, esse objeto terá
sua resposta de sucesso omitida automáticamente pelo endpoint, já que respostas de sucesso de middlewares não são
exibidas para o cliente final http, porém sua resposta será armazenada ao longo da requisição http feita no endpoint,
podendo ser manipulada.

Por exemplo, um `beforeware` quando mencionado no endpoint, ele será utilizado como middleware de pré-requisições, isto
é, ele será chamado antes dos backends principais do endpoint, então podemos, por exemplo, ter um middleware
de manipulação de device, como no json de configuração acima, aonde ele irá chamar esse backend de middleware
configurado no endpoint como `beforeware`, validando e salvando o dispositivo a partir de informações do header da
requisição, caso o backend responda um código de status de falha, no exemplo, o gateway abortará todos os backends
seguintes retornando o que o backend de device respondeu, caso tenha retornado um código de status de sucesso, ele irá
modificar o header de todas as requisições seguintes (`propagate:true`), adicionando o campo `X-Device-Id`, com o valor
do id do body de resposta do próprio backend, podendo ser utilizado nos outros backends seguintes do endpoint.

Para saber mais sobre os `modifiers` [veja](#modifiers).

Para entender melhor essa ferramenta poderosissíma, na prática, veja os exemplos de middlewares usados como
`beforeware` e `afterware` feitos no projeto
de [playground](https://github.com/GabrielHCataldo/gopen-gateway-playground).

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

