package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
)


var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

var registered = false

func init() {
	if registered {
		return
	}
	registered = true

	// ลงทะเบียน metric ที่เราสร้างไว้
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	
	
	// ✅ ลงทะเบียน go runtime metric แบบปลอดภัย
	if err := prometheus.Register(prometheus.NewGoCollector()); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			_ = are.ExistingCollector // ใช้อันเดิมได้เลย
		} else {
			panic(err) // Error อื่นให้ panic ไปเลย
		}
	}
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		requestCount.WithLabelValues(c.Request.Method, c.FullPath(), fmt.Sprint(status)).Inc()
		requestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
	}
}