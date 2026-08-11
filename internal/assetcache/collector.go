package assetcache

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type statusReader interface {
	Read(context.Context) (Status, error)
}

type metricsReader interface {
	Read(context.Context) (*Interval, error)
}

// Collector collects Content Caching status and traffic observations.
type Collector struct {
	status  statusReader
	metrics metricsReader
	timeout time.Duration
	logger  *slog.Logger
}

// NewCollector creates a Content Caching collector.
func NewCollector(status statusReader, metrics metricsReader, timeout time.Duration, logger *slog.Logger) *Collector {
	return &Collector{
		status:  status,
		metrics: metrics,
		timeout: timeout,
		logger:  logger,
	}
}

// Describe sends the collector's metric descriptions to ch.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range allDescs {
		ch <- desc
	}
}

// Collect reads both local sources and emits all available metrics.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	statusOK := c.collectStatus(ch)
	metricsOK := c.collectMetrics(ch)

	ch <- prometheus.MustNewConstMetric(statusUpDesc, prometheus.GaugeValue, boolFloat(statusOK))
	ch <- prometheus.MustNewConstMetric(metricsDBUpDesc, prometheus.GaugeValue, boolFloat(metricsOK))
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, boolFloat(statusOK && metricsOK))
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(start).Seconds())
}

func (c *Collector) collectStatus(ch chan<- prometheus.Metric) bool {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	status, err := c.status.Read(ctx)
	if err != nil {
		if c.logger != nil {
			c.logger.Warn("status collection failed", "err", err)
		}
		return false
	}

	ch <- prometheus.MustNewConstMetric(activatedDesc, prometheus.GaugeValue, boolFloat(status.Activated))
	ch <- prometheus.MustNewConstMetric(activeDesc, prometheus.GaugeValue, boolFloat(status.Active))
	ch <- prometheus.MustNewConstMetric(
		statusInfoDesc,
		prometheus.GaugeValue,
		1,
		status.CacheStatus,
		status.StartupStatus,
		strconv.Itoa(status.RegistrationStatus),
	)
	ch <- prometheus.MustNewConstMetric(cacheUsedBytesDesc, prometheus.GaugeValue, status.CacheUsedBytes)
	ch <- prometheus.MustNewConstMetric(cacheLimitBytesDesc, prometheus.GaugeValue, status.CacheLimitBytes)
	ch <- prometheus.MustNewConstMetric(cacheFreeBytesDesc, prometheus.GaugeValue, status.CacheFreeBytes)
	ch <- prometheus.MustNewConstMetric(personalCacheUsedBytesDesc, prometheus.GaugeValue, status.PersonalCacheUsedBytes)
	ch <- prometheus.MustNewConstMetric(personalCacheLimitBytesDesc, prometheus.GaugeValue, status.PersonalCacheLimitBytes)
	ch <- prometheus.MustNewConstMetric(personalCacheFreeBytesDesc, prometheus.GaugeValue, status.PersonalCacheFreeBytes)
	ch <- prometheus.MustNewConstMetric(cachePressureDesc, prometheus.GaugeValue, status.MaxCachePressureRatio)
	ch <- prometheus.MustNewConstMetric(parentCachesDesc, prometheus.GaugeValue, float64(status.ParentCount))
	ch <- prometheus.MustNewConstMetric(peerCachesDesc, prometheus.GaugeValue, float64(status.PeerCount))
	if !status.TotalsStartTime.IsZero() {
		ch <- prometheus.MustNewConstMetric(totalsStartTimestampDesc, prometheus.GaugeValue, float64(status.TotalsStartTime.Unix()))
	}
	emitCounter(ch, servedBytesTotalDesc, status.BytesServedToChildren, "children")
	emitCounter(ch, servedBytesTotalDesc, status.BytesServedToClients, "clients")
	emitCounter(ch, servedBytesTotalDesc, status.BytesServedToPeers, "peers")
	emitCounter(ch, storedBytesTotalDesc, status.BytesStoredFromOrigin, "origin")
	emitCounter(ch, storedBytesTotalDesc, status.BytesStoredFromParents, "parents")
	emitCounter(ch, storedBytesTotalDesc, status.BytesStoredFromPeers, "peers")
	emitCounter(ch, importedBytesTotalDesc, status.BytesImportedTotal)
	emitCounter(ch, droppedBytesTotalDesc, status.BytesDroppedTotal)
	return true
}

func (c *Collector) collectMetrics(ch chan<- prometheus.Metric) bool {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	interval, err := c.metrics.Read(ctx)
	if err != nil {
		if c.logger != nil {
			c.logger.Warn("metrics database collection failed", "err", err)
		}
		return false
	}
	if interval == nil {
		return true
	}

	ch <- prometheus.MustNewConstMetric(metricsLastTimestampDesc, prometheus.GaugeValue, float64(interval.Timestamp.Unix()))
	ch <- prometheus.MustNewConstMetric(metricsIntervalDesc, prometheus.GaugeValue, interval.PeriodSeconds)
	emitObservations(ch, intervalServedBytesDesc, interval.BytesServed)
	emitObservations(ch, intervalRepliesDesc, interval.Replies)
	emitObservations(ch, intervalRequestsDesc, interval.Requests)
	ch <- prometheus.MustNewConstMetric(intervalImportedBytesDesc, prometheus.GaugeValue, interval.BytesImported)
	ch <- prometheus.MustNewConstMetric(intervalImportsDesc, prometheus.GaugeValue, interval.Imports)
	ch <- prometheus.MustNewConstMetric(intervalDroppedBytesDesc, prometheus.GaugeValue, interval.BytesDropped)
	ch <- prometheus.MustNewConstMetric(intervalPurgedBytesDesc, prometheus.GaugeValue, interval.BytesPurged)
	ch <- prometheus.MustNewConstMetric(intervalRejectedRequestsDesc, prometheus.GaugeValue, interval.RequestsRejectedNoSpace)
	ch <- prometheus.MustNewConstMetric(intervalCachePressureDesc, prometheus.GaugeValue, interval.CachePressureRatio)
	return true
}

func emitCounter(ch chan<- prometheus.Metric, desc *prometheus.Desc, value float64, labels ...string) {
	ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, value, labels...)
}

func emitObservations(ch chan<- prometheus.Metric, desc *prometheus.Desc, observations []Observation) {
	for _, observation := range observations {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, observation.Value, observation.Source)
	}
}

func boolFloat(value bool) float64 {
	if value {
		return 1
	}
	return 0
}
