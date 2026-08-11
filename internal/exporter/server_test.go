package exporter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestServerServesMetricsAtRoot(t *testing.T) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewGauge(prometheus.GaugeOpts{Name: "test_metric"}))
	mux := http.NewServeMux()
	server := Server{Gatherer: registry, MetricsPath: "/"}
	if err := server.Register(mux); err != nil {
		t.Fatalf("register server: %v", err)
	}

	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), "test_metric") {
		t.Fatalf("response does not contain metrics: %s", response.Body.String())
	}
}

func TestServerRejectsMetricsPathWithoutLeadingSlash(t *testing.T) {
	server := Server{Gatherer: prometheus.NewRegistry(), MetricsPath: "metrics"}
	if err := server.Register(http.NewServeMux()); err == nil {
		t.Fatal("expected invalid metrics path to fail")
	}
}
