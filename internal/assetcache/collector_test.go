package assetcache

import (
	"context"
	"errors"
	"net"
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
	if hasMetric(families, "assetcache_tcp_connect_success") || hasMetric(families, "assetcache_listen_port") {
		t.Fatal("did not expect endpoint observations after status source failure")
	}
}

func TestCollectorTCPObservation(t *testing.T) {
	for _, tt := range []struct {
		name, network, ip string
	}{
		{"IPv4", "tcp4", "127.0.0.1"},
		{"IPv6", "tcp6", "::1"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			listener, err := net.ListenTCP(tt.network, &net.TCPAddr{IP: net.ParseIP(tt.ip)})
			if err != nil {
				if tt.network == "tcp6" {
					t.Skipf("IPv6 loopback unavailable: %v", err)
				}
				t.Fatalf("listen: %v", err)
			}
			t.Cleanup(func() { _ = listener.Close() })
			port := listener.Addr().(*net.TCPAddr).Port
			status := Status{Activated: true, Active: true, Port: &port, PrivateAddresses: []string{"invalid", "192.0.2.10", tt.ip}}

			for _, connected := range []bool{true, false} {
				if !connected {
					if err := listener.Close(); err != nil {
						t.Fatalf("close listener: %v", err)
					}
				}
				collector := NewCollector(fakeStatusReader{status: status}, fakeMetricsReader{}, time.Second, nil)
				registry := prometheus.NewPedanticRegistry()
				registry.MustRegister(collector)
				families, err := registry.Gather()
				if err != nil {
					t.Fatalf("gather metrics: %v", err)
				}
				if got := metricValue(t, families, "assetcache_tcp_connect_success"); got != boolFloat(connected) {
					t.Fatalf("TCP success = %v, want connected=%v", got, connected)
				}
				if got := metricValue(t, families, "assetcache_listen_port"); got != float64(port) {
					t.Fatalf("listen port = %v, want %d", got, port)
				}
				for _, name := range []string{"assetcache_active", "assetcache_status_up", "assetcache_metrics_db_up", "assetcache_scrape_success"} {
					if got := metricValue(t, families, name); got != 1 {
						t.Fatalf("%s = %v, want 1 regardless of TCP outcome", name, got)
					}
				}
			}
		})
	}
}

func TestCollectorOmitsUnattemptedTCPObservation(t *testing.T) {
	port, zero, invalid := 54321, 0, 65536
	for _, tt := range []struct {
		name    string
		status  Status
		timeout time.Duration
	}{
		{"inactive", Status{Port: &port, PrivateAddresses: []string{"127.0.0.1"}}, time.Second},
		{"missing port", Status{Active: true, PrivateAddresses: []string{"127.0.0.1"}}, time.Second},
		{"zero port", Status{Active: true, Port: &zero, PrivateAddresses: []string{"127.0.0.1"}}, time.Second},
		{"invalid port", Status{Active: true, Port: &invalid, PrivateAddresses: []string{"127.0.0.1"}}, time.Second},
		{"missing addresses", Status{Active: true, Port: &port}, time.Second},
		{"unavailable address", Status{Active: true, Port: &port, PrivateAddresses: []string{"192.0.2.10"}}, time.Second},
		{"invalid address", Status{Active: true, Port: &port, PrivateAddresses: []string{"invalid"}}, time.Second},
		{"expired budget", Status{Active: true, Port: &port, PrivateAddresses: []string{"127.0.0.1"}}, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			registry := prometheus.NewPedanticRegistry()
			registry.MustRegister(NewCollector(fakeStatusReader{status: tt.status}, fakeMetricsReader{}, tt.timeout, nil))
			families, err := registry.Gather()
			if err != nil {
				t.Fatalf("gather metrics: %v", err)
			}
			if hasMetric(families, "assetcache_tcp_connect_success") {
				t.Fatal("unattempted TCP check must be absent")
			}
			if got := metricValue(t, families, "assetcache_active"); got != boolFloat(tt.status.Active) {
				t.Fatalf("active = %v, want %v", got, tt.status.Active)
			}
			if got := metricValue(t, families, "assetcache_scrape_success"); got != 1 {
				t.Fatalf("scrape success = %v, want 1", got)
			}
			if tt.status.Port == nil {
				if hasMetric(families, "assetcache_listen_port") {
					t.Fatal("unreported port must be absent")
				}
			} else if got := metricValue(t, families, "assetcache_listen_port"); got != float64(*tt.status.Port) {
				t.Fatalf("listen port = %v, want %d", got, *tt.status.Port)
			}
		})
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
