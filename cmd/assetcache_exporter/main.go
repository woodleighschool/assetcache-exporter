package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/bootstrap"

	"github.com/woodleighschool/assetcache-exporter/internal/assetcache"
)

var timeout = kingpin.Flag("collector.timeout", "Maximum time allowed for each local data source.").Default("5s").Duration()

func main() {
	runner := bootstrap.New(bootstrap.Config{
		App:                   kingpin.CommandLine,
		Name:                  "assetcache_exporter",
		Description:           "Prometheus exporter for Apple Content Caching",
		DefaultAddress:        ":9200",
		ReadHeaderTimeout:     5 * time.Second,
		MetricsHandlerFactory: newMetricsHandler,
	})
	if err := runner.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newMetricsHandler(b *bootstrap.Bootstrap) (http.Handler, error) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		versioncollector.NewCollector("assetcache_exporter"),
		assetcache.NewCollector(
			assetcache.NewStatusSource(),
			assetcache.NewMetricsSource(),
			*timeout,
			b.Logger,
		),
	)
	if !b.DisableExporterMetrics {
		registry.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
	}

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		MaxRequestsInFlight: b.MaxRequests,
	})
	if !b.DisableExporterMetrics {
		handler = promhttp.InstrumentMetricHandler(registry, handler)
	}
	return handler, nil
}
