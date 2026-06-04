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

const tracerName = "gopen-gateway"

// Setup inicializa o OpenTelemetry com exportadores OTLP para traces, métricas e logs.
//
// Se OTEL_EXPORTER_OTLP_ENDPOINT estiver definido e o environment não for "local",
// configura exportadores completos (traces + metrics + logs) com autenticação.
// Caso contrário, configura apenas um TracerProvider local sem exportação.
//
// Autenticação (verificada na seguinte ordem de prioridade):
//  1. OTEL_EXPORTER_OTLP_HEADERS — formato padrão OTel "key=value,key2=value2"
//  2. OTEL_EXPORTER_OTLP_INSTANCE_ID + OTEL_EXPORTER_OTLP_API_TOKEN — Basic Auth
//
// Retorna uma função de shutdown que deve ser chamada no encerramento da aplicação.
func Setup(ctx context.Context, serviceName, serviceVersion, environment string) (func(context.Context) error, error) {
	// Configura propagadores globais W3C TraceContext + Baggage.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			attribute.String("service.name", serviceName),
			attribute.String("service.version", serviceVersion),
			attribute.String("deployment.environment.name", environment),
			attribute.String("service.instance.id", resolveInstanceID()),
		),
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

// Tracer retorna o tracer nomeado do gateway.
func Tracer() trace.Tracer {
	return otel.Tracer(tracerName)
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

// resolveInstanceID retorna identificador único para a instância do serviço.
// Prioridade: RAILWAY_REPLICA_ID > hostname > "unknown".
func resolveInstanceID() string {
	if replicaID := os.Getenv("RAILWAY_REPLICA_ID"); checker.IsNotEmpty(replicaID) {
		return replicaID
	}
	if hostname, err := os.Hostname(); checker.IsNil(err) && checker.IsNotEmpty(hostname) {
		return hostname
	}
	return "unknown"
}
