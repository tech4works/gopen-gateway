<img src="logo-gateway.png" alt="">

[![Project status](https://img.shields.io/badge/version-v1.0.0_beta-yellow.svg)](https://github.com/GabrielHCataldo/gopen-gateway/releases/tag/v1.0.0-beta)
[![Open Source Helpers](https://www.codetriage.com/gabrielhcataldo/gopen-gateway/badges/users.svg)](https://www.codetriage.com/gabrielhcataldo/gopen-gateway)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/GabrielHCataldo/gopen-gateway)](https://www.tickgit.com/browse?repo=github.com/GabrielHCataldo/gopen-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/GabrielHCataldo/gopen-gateway)](https://goreportcard.com/report/github.com/GabrielHCataldo/gopen-gateway)
[![GoDoc](https://godoc.org/github/GabrielHCataldo/gopen-gateway?status.svg)](https://pkg.go.dev/github.com/GabrielHCataldo/gopen-gateway/helper)

[//]: # ([![build workflow]&#40;https://github.com/GabrielHCataldo/gopen-gateway/actions/workflows/go.yml/badge.svg&#41;]&#40;https://github.com/GabrielHCataldo/gopen-gateway/actions&#41;)

---

![United States](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/us.png "United States")
[Inglês](https://github/GabrielHCataldo/gopen-gateway/blob/main/README.md) |
![Spain](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/es.png "Spain")
[Espanhol](https://github/GabrielHCataldo/gopen-gateway/blob/main/README.es.md)

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

Usabilidade e documentação
-----------
---


Como contríbuir?
------------
---


Agradecimentos
------------
---

