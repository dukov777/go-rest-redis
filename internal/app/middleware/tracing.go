package middleware

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// InitJaeger creates a new trace provider instance
func InitJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

// TracingMiddleware intercepts requests to add tracing information
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer, closer := InitJaeger("my-go-project")
		defer closer.Close()

		// Start a new span referring to the span context of the incoming request
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan(r.URL.Path, opentracing.ChildOf(spanCtx))
		defer span.Finish()

		// Store the span in the context of the request
		ctx := opentracing.ContextWithSpan(r.Context(), span)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		// Optionally log additional information about the request
		span.SetTag("http.method", r.Method)
		span.SetTag("http.remote_addr", r.RemoteAddr)
	})
}
