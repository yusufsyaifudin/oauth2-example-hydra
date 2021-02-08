package tracer

import (
	"io"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

type nullCloser struct{}

func (*nullCloser) Close() error { return nil }

// New returns a new tracer
func New(enabled bool, serviceName, hostPort string, probability float64) (opentracing.Tracer, io.Closer) {
	if !enabled {
		tracer := new(opentracing.NoopTracer)
		closer := new(nullCloser)
		return tracer, closer
	}

	serviceName = strings.TrimSpace(serviceName)

	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: probability,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  hostPort, // localhost:5775
		},
	}

	// using null logger or output to std out
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.NullLogger),
	)

	// if cannot initiate using real service (because jaeger instance is down, then use noop tracer)
	if err != nil {
		tracer = new(opentracing.NoopTracer)
		closer = new(nullCloser)
		return tracer, closer
	}

	return tracer, closer
}
