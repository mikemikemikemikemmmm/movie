package otel

import (
	"backend/internal/config"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var GlobalTracer trace.Tracer

const NotEndSpan = "notEndSpan"

func InitTracer(ctx context.Context) (*tracesdk.TracerProvider, error) {
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(config.GetConfig().TempoURL), // Tempo OTLP HTTP 端點
		otlptracehttp.WithInsecure(),                            // 如果 Tempo 沒有 TLS
	)
	if err != nil {
		return nil, err
	}

	// 建立 TracerProvider
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("global-tracer"),
		)),
	)

	// 設置全局 TracerProvider
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	GlobalTracer = otel.Tracer("global-tracer")
	return tracerProvider, nil
}
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/metrics" {
			c.Next()
			return
		}
		tracerCtx, span := GlobalTracer.Start(c.Request.Context(), c.FullPath())
		defer func() {
			if c.GetHeader(NotEndSpan) != "true" {
				span.End()
			}
		}()
		span.SetAttributes(attribute.String("http.method", c.Request.Method))
		span.SetAttributes(attribute.String("http.url", c.FullPath()))
		c.Request = c.Request.WithContext(tracerCtx)
		c.Next()
	}
}

func InjectTraceToKafkaHeader(ctx context.Context, msg *kafka.Message) {
	if msg.Headers == nil {
		msg.Headers = []kafka.Header{}
	}

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	for k, v := range carrier {
		msg.Headers = append(msg.Headers, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}
}
