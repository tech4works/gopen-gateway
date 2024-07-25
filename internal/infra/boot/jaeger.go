package boot

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
)

type noopJaegerLogger struct{}

func (l *noopJaegerLogger) Error(_ string) {}

func (l *noopJaegerLogger) Infof(_ string, _ ...interface{}) {}

func InitJaeger(host string) (opentracing.Tracer, io.Closer, error) {
	jaegerConfig := &jaegercfg.Configuration{
		ServiceName: "gopen-gateway",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: host,
		},
	}
	return jaegerConfig.NewTracer(jaegercfg.Logger(&noopJaegerLogger{}))
}
