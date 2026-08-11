package assetcache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestCollectorKeepsMetricsWhenStatusFails(t *testing.T) {
	collector := NewCollector(
		fakeStatusReader{err: errors.New("status unavailable")},
		fakeMetricsReader{interval: &Interval{
			Timestamp:     time.Unix(1_700_000_000, 0),
			PeriodSeconds: 60,
			BytesServed:   []Observation{{Source: "cache", Value: 42}},
		}},
		time.Second,
		nil,
	)
	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(collector)

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
	if got := metricValue(t, families, "assetcache_status_up"); got != 0 {
		t.Fatalf("status up = %v, want 0", got)
	}
	if got := metricValue(t, families, "assetcache_metrics_db_up"); got != 1 {
		t.Fatalf("metrics db up = %v, want 1", got)
	}
	if got := metricValue(t, families, "assetcache_scrape_success"); got != 0 {
		t.Fatalf("scrape success = %v, want 0", got)
	}
	if got := metricValue(t, families, "assetcache_interval_served_bytes"); got != 42 {
		t.Fatalf("interval served bytes = %v, want 42", got)
	}
	if got := metricType(t, families, "assetcache_interval_served_bytes"); got != dto.MetricType_GAUGE {
		t.Fatalf("interval served bytes type = %v, want gauge", got)
	}
	if hasMetric(families, "assetcache_active") {
		t.Fatal("did not expect status metrics after status source failure")
	}
}

func TestCollectorKeepsStatusWhenMetricsDatabaseFails(t *testing.T) {
	collector := NewCollector(
		fakeStatusReader{status: Status{Active: true, Activated: true}},
		fakeMetricsReader{err: errors.New("database unavailable")},
		time.Second,
		nil,
	)
	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(collector)

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
	if got := metricValue(t, families, "assetcache_status_up"); got != 1 {
		t.Fatalf("status up = %v, want 1", got)
	}
	if got := metricValue(t, families, "assetcache_metrics_db_up"); got != 0 {
		t.Fatalf("metrics db up = %v, want 0", got)
	}
	if got := metricValue(t, families, "assetcache_scrape_success"); got != 0 {
		t.Fatalf("scrape success = %v, want 0", got)
	}
	if got := metricValue(t, families, "assetcache_active"); got != 1 {
		t.Fatalf("active = %v, want 1", got)
	}
	if hasMetric(families, "assetcache_metrics_last_timestamp_seconds") {
		t.Fatal("did not expect interval metrics after database source failure")
	}
}

func TestCollectorTreatsReadableEmptyMetricsDatabaseAsUp(t *testing.T) {
	collector := NewCollector(
		fakeStatusReader{status: Status{}},
		fakeMetricsReader{},
		time.Second,
		nil,
	)
	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(collector)

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
	if got := metricValue(t, families, "assetcache_metrics_db_up"); got != 1 {
		t.Fatalf("metrics db up = %v, want 1", got)
	}
	if got := metricValue(t, families, "assetcache_scrape_success"); got != 1 {
		t.Fatalf("scrape success = %v, want 1", got)
	}
	if hasMetric(families, "assetcache_metrics_last_timestamp_seconds") {
		t.Fatal("did not expect interval timestamp for an empty database")
	}
}

type fakeStatusReader struct {
	status Status
	err    error
}

func (f fakeStatusReader) Read(context.Context) (Status, error) {
	return f.status, f.err
}

type fakeMetricsReader struct {
	interval *Interval
	err      error
}

func (f fakeMetricsReader) Read(context.Context) (*Interval, error) {
	return f.interval, f.err
}

func metricValue(t *testing.T, families []*dto.MetricFamily, name string) float64 {
	t.Helper()
	for _, family := range families {
		if family.GetName() != name || len(family.Metric) == 0 {
			continue
		}
		metric := family.Metric[0]
		if metric.Gauge != nil {
			return metric.Gauge.GetValue()
		}
		if metric.Counter != nil {
			return metric.Counter.GetValue()
		}
	}
	t.Fatalf("metric %q not found", name)
	return 0
}

func metricType(t *testing.T, families []*dto.MetricFamily, name string) dto.MetricType {
	t.Helper()
	for _, family := range families {
		if family.GetName() == name {
			return family.GetType()
		}
	}
	t.Fatalf("metric family %q not found", name)
	return dto.MetricType_UNTYPED
}

func hasMetric(families []*dto.MetricFamily, name string) bool {
	for _, family := range families {
		if family.GetName() == name {
			return true
		}
	}
	return false
}
