// Package exporter provides the assetcache-exporter HTTP surface.
package exporter

import (
	"errors"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
)

// Server registers the exporter HTTP handlers.
type Server struct {
	Gatherer    prometheus.Gatherer
	MetricsPath string
}

// Register adds the metrics endpoint and landing page to mux.
func (s Server) Register(mux *http.ServeMux) error {
	metricsPath := s.MetricsPath
	if metricsPath == "" {
		metricsPath = "/metrics"
	}
	if !strings.HasPrefix(metricsPath, "/") {
		return errors.New("metrics path must start with /")
	}

	mux.Handle(metricsPath, promhttp.HandlerFor(s.Gatherer, promhttp.HandlerOpts{}))
	if metricsPath == "/" {
		return nil
	}
	landingPage, err := web.NewLandingPage(web.LandingConfig{
		Name:        "assetcache_exporter",
		Description: "Prometheus exporter for Apple Content Caching",
		Version:     version.Info(),
		Links: []web.LandingLinks{
			{Address: metricsPath, Text: "Metrics"},
		},
	})
	if err != nil {
		return err
	}
	mux.Handle("/", landingPage)
	return nil
}
