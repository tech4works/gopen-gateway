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

package telemetry

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/tech4works/checker"
)

const tracerName = "gopen-gateway"

// Setup initializes the global OTel TracerProvider with an OTLP HTTP exporter.
// The exporter endpoint is read from OTEL_EXPORTER_OTLP_ENDPOINT (default: http://localhost:4318).
// Returns a shutdown function that must be called on application exit.
func Setup(ctx context.Context, serviceName, serviceVersion, environment string) (func(context.Context) error, error) {
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL(otlpEndpoint()),
	)
	if checker.NonNil(err) {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			attribute.String("deployment.environment", environment),
		),
	)
	if checker.NonNil(err) {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)

	// Configura o propagador global W3C TraceContext + Baggage.
	// Sem isso, otelhttp e otelgin não injetam nem extraem o traceparent.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

// Tracer returns the named tracer for this gateway.
func Tracer() trace.Tracer {
	return otel.Tracer(tracerName)
}

func otlpEndpoint() string {
	if ep := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); checker.IsNotEmpty(ep) {
		return ep
	}
	return "http://localhost:4318"
}
