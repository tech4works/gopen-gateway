/*
 * Copyright 2024 Tech4Works
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package telemetry configura os exportadores OpenTelemetry (traces, métricas e logs)
// para o gopen-gateway. Segue o mesmo padrão do sdk-manas boot, com suporte a
// OTLP HTTP exporters, autenticação via headers genéricos ou Basic Auth,
// Go runtime metrics e skip de exportação em ambiente local.
package telemetry

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelglobal "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/tech4works/checker"
)

const defaultApplicationName = "gopen-gateway"

// resolvedServiceName armazena o nome do serviço resolvido para uso no Tracer.
var resolvedServiceName string

// Setup inicializa o OpenTelemetry com exportadores OTLP para traces, métricas e logs.
//
// O service.name usa APPLICATION_NAME diretamente (fallback: "gopen-gateway").
// Informações de projeto e ambiente são propagadas via entity tags (tags.environment,
// tags.project, tags.role, tags.version) para agrupamento no New Relic e plataformas similares.
//
// Variáveis de ambiente usadas:
//   - APPLICATION_NAME — nome da aplicação (fallback: "gopen-gateway")
//   - PROJECT_NAME — nome do projeto (usado em tags.project)
//
// Retorna uma função de shutdown que deve ser chamada no encerramento da aplicação.
func Setup(ctx context.Context, serviceVersion, environment string) (func(context.Context) error, error) {
	serviceName := resolveServiceName(environment)
	resolvedServiceName = serviceName
	applicationName := resolveApplicationName()

	// Configura propagadores globais W3C TraceContext + Baggage.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	attrs := []attribute.KeyValue{
		attribute.String("service.name", serviceName),
		attribute.String("service.version", serviceVersion),
		attribute.String("deployment.environment.name", environment),
		attribute.String("service.instance.id", resolveInstanceID()),
		attribute.String("application.name", applicationName),
		attribute.String("tags.environment", environment),
		attribute.String("tags.project", os.Getenv("PROJECT_NAME")),
		attribute.String("tags.role", "server"),
		attribute.String("tags.version", serviceVersion),
	}

	// Tags extras opcionais
	if projectName := os.Getenv("PROJECT_NAME"); checker.IsNotEmpty(projectName) {
		attrs = append(attrs, attribute.String("project.name", projectName))
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes("", attrs...),
	)
	if checker.NonNil(err) {
		return nil, err
	}

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")

	// Sem OTLP ou ambiente local — TracerProvider sem exporter
	if checker.IsEmpty(endpoint) || checker.Equals(environment, "local") {
		tp := sdktrace.NewTracerProvider(sdktrace.WithResource(res))
		otel.SetTracerProvider(tp)
		return tp.Shutdown, nil
	}

	return setupOTLP(ctx, res, endpoint)
}

// Tracer retorna o tracer nomeado do gateway usando o service name resolvido.
func Tracer() trace.Tracer {
	if checker.IsNotEmpty(resolvedServiceName) {
		return otel.Tracer(resolvedServiceName)
	}
	return otel.Tracer(defaultApplicationName)
}

// setupOTLP configura exportadores OTLP para traces, métricas e logs.
// Retorna shutdown function que encerra todos os providers.
func setupOTLP(ctx context.Context, res *resource.Resource, endpoint string) (func(context.Context) error, error) {
	// Monta opções base com endpoint por sinal
	traceOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpointURL(endpoint + "/v1/traces"),
	}
	logOpts := []otlploghttp.Option{
		otlploghttp.WithEndpointURL(endpoint + "/v1/logs"),
	}
	metricOpts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpointURL(endpoint + "/v1/metrics"),
	}

	// Resolve headers de autenticação
	headers := resolveOTLPHeaders()
	if len(headers) > 0 {
		traceOpts = append(traceOpts, otlptracehttp.WithHeaders(headers))
		logOpts = append(logOpts, otlploghttp.WithHeaders(headers))
		metricOpts = append(metricOpts, otlpmetrichttp.WithHeaders(headers))
	}

	// Trace exporter
	traceExporter, err := otlptracehttp.New(ctx, traceOpts...)
	if checker.NonNil(err) {
		return nil, fmt.Errorf("failed to create otlp trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
	)
	otel.SetTracerProvider(tp)

	// Metric exporter
	metricExporter, err := otlpmetrichttp.New(ctx, metricOpts...)
	if checker.NonNil(err) {
		return nil, fmt.Errorf("failed to create otlp metric exporter: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)
	otel.SetMeterProvider(mp)

	// Log exporter
	logExporter, err := otlploghttp.New(ctx, logOpts...)
	if checker.NonNil(err) {
		return nil, fmt.Errorf("failed to create otlp log exporter: %w", err)
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)
	otelglobal.SetLoggerProvider(lp)

	// Go runtime metrics (goroutines, GC, memória)
	if err := runtime.Start(); checker.NonNil(err) {
		fmt.Printf("WARNING: failed to start Go runtime metrics: %s\n", err)
	}

	fmt.Printf("OTLP exporters configured: endpoint=%s (traces, metrics, logs)\n", endpoint)

	// Shutdown composto: encerra todos os providers
	shutdown := func(ctx context.Context) error {
		var firstErr error
		if err := tp.Shutdown(ctx); checker.NonNil(err) && checker.IsNil(firstErr) {
			firstErr = err
		}
		if err := mp.Shutdown(ctx); checker.NonNil(err) && checker.IsNil(firstErr) {
			firstErr = err
		}
		if err := lp.Shutdown(ctx); checker.NonNil(err) && checker.IsNil(firstErr) {
			firstErr = err
		}
		return firstErr
	}

	return shutdown, nil
}

// resolveOTLPHeaders resolve os headers de autenticação para os exportadores OTLP.
//
// Prioridade:
//  1. OTEL_EXPORTER_OTLP_HEADERS (formato "key=value,key2=value2")
//  2. OTEL_EXPORTER_OTLP_INSTANCE_ID + OTEL_EXPORTER_OTLP_API_TOKEN (Basic Auth)
//
// Retorna nil se nenhuma autenticação estiver configurada.
func resolveOTLPHeaders() map[string]string {
	// Prioridade 1: headers genéricos (padrão OTel — funciona com New Relic, Elastic, etc.)
	rawHeaders := os.Getenv("OTEL_EXPORTER_OTLP_HEADERS")
	if checker.IsNotEmpty(rawHeaders) {
		headers := make(map[string]string)
		for _, pair := range strings.Split(rawHeaders, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
			if checker.Equals(len(parts), 2) {
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
		return headers
	}

	// Prioridade 2: Basic Auth (Grafana Cloud)
	instanceID := os.Getenv("OTEL_EXPORTER_OTLP_INSTANCE_ID")
	apiToken := os.Getenv("OTEL_EXPORTER_OTLP_API_TOKEN")
	if checker.IsNotEmpty(instanceID) && checker.IsNotEmpty(apiToken) {
		authValue := "Basic " + base64.StdEncoding.EncodeToString([]byte(instanceID+":"+apiToken))
		return map[string]string{"Authorization": authValue}
	}

	return nil
}

// resolveServiceName retorna o nome do serviço para o OTLP resource.
// Usa apenas APPLICATION_NAME (fallback: "gopen-gateway").
// Informações de projeto e ambiente são separadas via entity tags.
func resolveServiceName(environment string) string {
	return resolveApplicationName()
}

// resolveApplicationName retorna APPLICATION_NAME ou fallback "gopen-gateway".
func resolveApplicationName() string {
	if appName := os.Getenv("APPLICATION_NAME"); checker.IsNotEmpty(appName) {
		return appName
	}
	return defaultApplicationName
}

// resolveInstanceID retorna identificador único para a instância do serviço.
// Prioridade: INSTANCE_REPLICA_ID > hostname > "unknown".
func resolveInstanceID() string {
	if replicaID := os.Getenv("INSTANCE_REPLICA_ID"); checker.IsNotEmpty(replicaID) {
		return replicaID
	}
	if hostname, err := os.Hostname(); checker.IsNil(err) && checker.IsNotEmpty(hostname) {
		return hostname
	}
	return "unknown"
}
