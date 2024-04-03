<img src="logo-gateway.png" alt="">

[![Project status](https://img.shields.io/badge/version-v1.0.0_beta-yellow.svg)](https://github.com/GabrielHCataldo/gopen-gateway/releases/tag/v1.0.0-beta)
[![Open Source Helpers](https://www.codetriage.com/gabrielhcataldo/gopen-gateway/badges/users.svg)](https://www.codetriage.com/gabrielhcataldo/gopen-gateway)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/GabrielHCataldo/gopen-gateway)](https://www.tickgit.com/browse?repo=github.com/GabrielHCataldo/gopen-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/GabrielHCataldo/gopen-gateway)](https://goreportcard.com/report/github.com/GabrielHCataldo/gopen-gateway)
[![GoDoc](https://godoc.org/github/GabrielHCataldo/gopen-gateway?status.svg)](https://pkg.go.dev/github.com/GabrielHCataldo/gopen-gateway/helper)

[//]: # ([![build workflow]&#40;https://github.com/GabrielHCataldo/gopen-gateway/actions/workflows/go.yml/badge.svg&#41;]&#40;https://github.com/GabrielHCataldo/gopen-gateway/actions&#41;)

---

![United States](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/us.png "United States")
[Inglés](https://github.com/GabrielHCataldo/gopen-gateway/blob/main/README.md) |
![Brazil](https://raw.githubusercontent.com/stevenrskelton/flag-icon/master/png/16/country-4x3/br.png "Brazil")
[Portugués](https://github.com/GabrielHCataldo/gopen-gateway/blob/main/README.pt-br.md) |

El proyecto GOPEN fue creado con el objetivo de ayudar a los desarrolladores a tener un API Gateway robusto y fácil de usar,
con la oportunidad de trabajar en mejoras, uniendo a la comunidad y lo más importante, sin gastar nada. Él era
desarrollado, ya que muchos API Gateways gratuitos en el mercado no satisfacen muchas necesidades mínimas
a una aplicación, induciéndola a actualizarse.

Con este nuevo API Gateway no necesitarás equilibrar tus placas para ahorrar en tu infraestructura y arquitectura,
Vea a continuación todos los recursos disponibles:

- Configuración Json simplificada para múltiples entornos.
- Configuración rápida de variables de entorno para múltiples entornos.
- Versionado mediante configuración json.
- Ejecución vía docker con recarga en caliente opcional.
- Configuración de tiempo de espera global y local para cada endpoint.
- Configuración de caché global y local para cada endpoint, con personalización de la estrategia de claves de almacenamiento.
- Almacenamiento en caché local o global usando Redis
- Configuración del limitador de tamaño global y local para cada punto final, limitando el tamaño del encabezado, cuerpo y multiparte
  Memoria.
- Configuración de limitador de velocidad global y local para cada endpoint, limitando por tiempo y burst por IP.
- Configuración de seguridad CORS con validaciones de origen, método http y encabezados.
- Configuración global de múltiples middlewares, para ser utilizados posteriormente en el endpoint si así se indica.
- Filtrado personalizado para enviar encabezados y consultas a los backends de los puntos finales.
- Procesamiento de múltiples backends, incluido beforeware, main y afterware para el endpoint.
- Configuración personalizada para cancelar el proceso de ejecución del backend mediante el código de estado devuelto.
- Modificadores para todos los contenidos de solicitud y respuesta (Código de estado, Ruta, Encabezado, Parámetros, Consulta, Cuerpo)
  a nivel global (solicitud/respuesta de punto final) y nivel local (solicitud/respuesta de backend actual) con acciones de eliminación,
  adición, cambio, sustitución y cambio de nombre.
- Obtener el valor a modificar de las variables de entorno, la solicitud actual, el historial de respuestas del punto final,
  o incluso el valor pasado en la configuración.
- Ejecute los modificadores en el contexto que desee, antes o después de una solicitud de backend, usted decide.
- Realizar las modificaciones reflejadas en todas las solicitudes/respuestas posteriores, utilizando las mismas a nivel global.
- Omita la respuesta de un backend si es necesario, no se imprimirá en la respuesta del endpoint.
- Omita el cuerpo de la solicitud de su backend si es necesario.
- Agregue sus múltiples respuestas de backends si lo desea, pudiendo personalizar el nombre del campo que se asignará
  respuesta de fondo.
- Agrupe el cuerpo de su respuesta de backend en un campo de respuesta de endpoint específico.
- Personalización del tipo de respuesta del endpoint, que puede ser JSON, TEXTO y XML.
- Tener más observabilidad con el registro automático del ID de seguimiento en el encabezado de solicitudes y registros posteriores.
  estructurado.

Usabilidad y documentación
-----------
---


¿Cómo contribuir?
------------
---


Agradecimientos
------------
---

