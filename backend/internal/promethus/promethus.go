package promethus

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var Registry = prometheus.NewRegistry()
var httpRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of HTTP requests",
	},
	[]string{"path", "method", "status"},
)
var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.DefBuckets, // 0.005s ~ 10s
	},
	[]string{"path"},
)

func InitPromethus() {
	Registry.MustRegister(httpDuration)
}
func TimerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/metrics" {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues(c.FullPath()).Observe(duration)
	}
}
