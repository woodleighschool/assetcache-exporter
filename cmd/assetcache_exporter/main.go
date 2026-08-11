package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"

	"github.com/woodleighschool/assetcache-exporter/internal/assetcache"
	"github.com/woodleighschool/assetcache-exporter/internal/exporter"
)

var (
	metricsPath  = kingpin.Flag("web.telemetry-path", "Path under which to expose exporter metrics.").Default("/metrics").String()
	timeout      = kingpin.Flag("collector.timeout", "Maximum time allowed for each local data source.").Default("5s").Duration()
	toolkitFlags = kingpinflag.AddFlags(kingpin.CommandLine, ":9200")
)

func main() {
	os.Exit(run())
}

func run() int {
	promslogConfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.Version(version.Print("assetcache_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promslog.New(promslogConfig)
	logger.Info("Starting assetcache_exporter", "version", version.Info())
	logger.Info("Build context", "build_context", version.BuildContext())

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		versioncollector.NewCollector("assetcache_exporter"),
		assetcache.NewCollector(
			assetcache.NewStatusSource(),
			assetcache.NewMetricsSource(),
			*timeout,
			logger,
		),
	)

	mux := http.NewServeMux()
	server := exporter.Server{
		Gatherer:    registry,
		MetricsPath: *metricsPath,
	}
	if err := server.Register(mux); err != nil {
		logger.Error("error creating HTTP handlers", "err", err)
		return 1
	}

	httpServer := &http.Server{
		Handler:           logRequests(mux, logger),
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := web.ListenAndServe(httpServer, toolkitFlags, logger); err != nil {
		logger.Error("error running HTTP server", "err", err)
		return 1
	}
	return 0
}

func logRequests(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.DebugContext(r.Context(), "request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
