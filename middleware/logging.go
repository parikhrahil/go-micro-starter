package middleware

import (
	"fmt"
	"time"

	"github.com/parikhrahil/go-micro-starter/logger"

	"github.com/gin-gonic/gin"
)

func LogHandler(log logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		ip := ctx.ClientIP()
		start := time.Now()

		ctx.Next()

		elapsedTimeMS := fmt.Sprintf("%.2fms", time.Since(start).Seconds()*1000)

		if raw != "" {
			path = path + "?" + raw
		}

		fields := []logger.Field{
			{
				Key:   "method",
				Value: ctx.Request.Method,
			},
			{
				Key:   "path",
				Value: path,
			},
			{
				Key:   "status",
				Value: ctx.Writer.Status(),
			},
			{
				Key:   "latency",
				Value: elapsedTimeMS,
			},
			{
				Key:   "ip",
				Value: ip,
			},
		}
		log.Info("HTTP Request", fields...)
	}
}
