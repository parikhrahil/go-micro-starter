package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parikhrahil/go-micro-starter/logger"
	"github.com/parikhrahil/go-micro-starter/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryOpts struct {
	Log         logger.Logger
	ServiceName string
}

func TelemetryHandler(opts *TelemetryOpts) gin.HandlerFunc {
	tracer := otel.GetTracerProvider().Tracer(opts.ServiceName)
	propogator := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		start := time.Now()

		ctx := propogator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		spanName := c.Request.Method + " " + c.Request.URL.Path
		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		var traceID, spanID string
		if span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
			spanID = span.SpanContext().SpanID().String()
			span.SetAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
			)
		}

		opts.Log.WithFields(logger.Field{
			Key:   "trace_id",
			Value: traceID,
		}, logger.Field{
			Key:   "span_id",
			Value: spanID,
		})

		c.Next()

		elapsedTime := time.Since(start)
		status := fmt.Sprintf("%d", c.Request.Response.StatusCode)
		span.SetAttributes(attribute.Int("http.status_code", c.Request.Response.StatusCode))

		if telemetry.GlobalMetrics != nil {
			telemetry.GlobalMetrics.HttpRequestsTotal.WithLabelValues(
				c.Request.Method, c.Request.URL.Path, status,
			).Inc()
			telemetry.GlobalMetrics.HttpRequestDuration.WithLabelValues(
				c.Request.Method, c.Request.URL.Path,
			).Observe(float64(elapsedTime.Milliseconds()))
		}
	}
}
