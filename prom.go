package cbl

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	apiCalledLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "server_api_latency",
			Help:       "server API latency",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"api_name", "method"},
	)
)

func init() {
	prometheus.MustRegister(apiCalledLatency)
}

// PromGinMiddleware prometheus for gin, register code:
// `router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})))`
func PromGinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start) // nanosecond
			apiCalledLatency.WithLabelValues(c.Request.RequestURI, c.Request.Method).Observe(float64(duration))
		}()

		c.Next()
	}
}
