package tracingutil

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	otelp "go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/trace"

	// "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/configs"
)

func InitTelemetry(c *configs.BaseConfig, zapLogger *zap.Logger) *otelp.Exporter {
	var (
		err error
		pe  *otelp.Exporter
	)

	if c.StatsEnabled {
		registry := prometheus.NewRegistry()
		registry.MustRegister(prometheus.NewGoCollector())

		pe, err = otelp.InstallNewPipeline(otelp.Config{
			Registry: registry,
		})
		if err != nil {
			zapLogger.Error("Failed to create the Prometheus stats exporter", zap.Error(err))
		}
	}

	if c.RemoteTrace.Enabled {
		// _, err := jaeger.InstallNewPipeline(
		// 	jaeger.WithCollectorEndpoint(c.RemoteTrace.TraceCollector),
		// 	jaeger.WithProcess(jaeger.Process{
		// 		ServiceName: c.Name,
		// 	}),
		// 	jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.TraceIDRatioBased(c.RemoteTrace.Ratio)}),
		// )
		// if err != nil {
		// 	zapLogger.Panic("err stdout.NewExporter", zap.Error(err))
		// }
	}

	return pe
}

const tracerKey = "otel-go-contrib-tracer"

type tracerKeyInt int

func Start(ctx context.Context, spanName string, opts ...trace.SpanOption) (context.Context, trace.Span) {
	switch c := ctx.(type) {
	case *gin.Context:
		tracer, ok := c.Get(tracerKey)
		if ok {
			c, span := tracer.(trace.Tracer).Start(c.Request.Context(), spanName, opts...)
			c = context.WithValue(c, tracerKeyInt(0), tracer)
			return c, span
		}
	default:
		tracer, ok := ctx.Value(tracerKeyInt(0)).(trace.Tracer)
		if ok {
			return tracer.Start(ctx, spanName, opts...)
		}
	}

	return otel.Tracer("unknown tracer").Start(ctx, spanName, opts...)
}
