package monitoring

import (
	"go-fiber-gorm/pkg/logger"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var (
	registry           *prometheus.Registry
	httpRequestsTotal  *prometheus.CounterVec
	httpRequestLatency *prometheus.HistogramVec
	dbQueryTotal       *prometheus.CounterVec
	dbQueryLatency     *prometheus.HistogramVec
	initOnce           sync.Once
)

// InitMetrics initializes all Prometheus metrics
func InitMetrics() {
	initOnce.Do(func() {
		logger.Info("Initializing Prometheus metrics")

		// Create a new registry
		registry = prometheus.NewRegistry()

		// HTTP request metrics
		httpRequestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		)

		httpRequestLatency = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		)

		// Database query metrics
		dbQueryTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"operation", "table"},
		)

		dbQueryLatency = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		)

		// Register metrics
		registry.MustRegister(
			httpRequestsTotal,
			httpRequestLatency,
			dbQueryTotal,
			dbQueryLatency,
		)
	})
}

// SetupMonitoring configures monitoring endpoints
func SetupMonitoring(app *fiber.App) {
	// Initialize metrics
	InitMetrics()

	// Dashboard endpoint
	app.Get("/dashboard", monitor.New())

	// Prometheus metrics endpoint
	promHandler := fasthttpadaptor.NewFastHTTPHandler(
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	)

	app.Get("/metrics", func(c *fiber.Ctx) error {
		promHandler(c.Context())
		return nil
	})

	logger.Info("Monitoring endpoints configured: /dashboard and /metrics")
}

// IncrementHTTPRequests increments the HTTP request counter
func IncrementHTTPRequests(method, path string, status int) {
	httpRequestsTotal.WithLabelValues(method, path, string(rune(status))).Inc()
}

// ObserveHTTPRequestLatency observes HTTP request latency
func ObserveHTTPRequestLatency(method, path string, latency float64) {
	httpRequestLatency.WithLabelValues(method, path).Observe(latency)
}

// IncrementDBQueries increments the database query counter
func IncrementDBQueries(operation, table string) {
	dbQueryTotal.WithLabelValues(operation, table).Inc()
}

// ObserveDBQueryLatency observes database query latency
func ObserveDBQueryLatency(operation, table string, latency float64) {
	dbQueryLatency.WithLabelValues(operation, table).Observe(latency)
}
