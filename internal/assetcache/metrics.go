package assetcache

import "github.com/prometheus/client_golang/prometheus"

var (
	scrapeSuccessDesc = prometheus.NewDesc(
		"assetcache_scrape_success",
		"Whether both local Content Caching data sources were collected successfully.",
		nil,
		nil,
	)
	scrapeDurationDesc = prometheus.NewDesc(
		"assetcache_scrape_duration_seconds",
		"Duration of the Content Caching scrape.",
		nil,
		nil,
	)
	statusUpDesc = prometheus.NewDesc(
		"assetcache_status_up",
		"Whether AssetCacheManagerUtil returned readable status.",
		nil,
		nil,
	)
	metricsDBUpDesc = prometheus.NewDesc(
		"assetcache_metrics_db_up",
		"Whether Apple's Content Caching metrics database was readable and compatible.",
		nil,
		nil,
	)
	activatedDesc = prometheus.NewDesc(
		"assetcache_activated",
		"Whether Content Caching is activated.",
		nil,
		nil,
	)
	activeDesc = prometheus.NewDesc(
		"assetcache_active",
		"Whether Content Caching is active.",
		nil,
		nil,
	)
	statusInfoDesc = prometheus.NewDesc(
		"assetcache_status_info",
		"Current textual Content Caching status.",
		[]string{"cache_status", "startup_status", "registration_status"},
		nil,
	)
	cacheUsedBytesDesc = prometheus.NewDesc(
		"assetcache_cache_used_bytes",
		"Logical size of content currently held in the cache.",
		nil,
		nil,
	)
	cacheLimitBytesDesc = prometheus.NewDesc(
		"assetcache_cache_limit_bytes",
		"Configured maximum cache size in bytes.",
		nil,
		nil,
	)
	cacheFreeBytesDesc = prometheus.NewDesc(
		"assetcache_cache_free_bytes",
		"Bytes currently available to Content Caching.",
		nil,
		nil,
	)
	personalCacheUsedBytesDesc = prometheus.NewDesc(
		"assetcache_personal_cache_used_bytes",
		"Logical size of personal iCloud content currently held in the cache.",
		nil,
		nil,
	)
	personalCacheLimitBytesDesc = prometheus.NewDesc(
		"assetcache_personal_cache_limit_bytes",
		"Maximum cache size available for personal iCloud content.",
		nil,
		nil,
	)
	personalCacheFreeBytesDesc = prometheus.NewDesc(
		"assetcache_personal_cache_free_bytes",
		"Bytes currently available for personal iCloud content.",
		nil,
		nil,
	)
	cachePressureDesc = prometheus.NewDesc(
		"assetcache_cache_pressure_ratio",
		"Maximum cache pressure reported during the last hour, from 0 to 1.",
		nil,
		nil,
	)
	parentCachesDesc = prometheus.NewDesc(
		"assetcache_parent_caches",
		"Number of parent content caches currently reported.",
		nil,
		nil,
	)
	peerCachesDesc = prometheus.NewDesc(
		"assetcache_peer_caches",
		"Number of peer content caches currently reported.",
		nil,
		nil,
	)
	totalsStartTimestampDesc = prometheus.NewDesc(
		"assetcache_totals_start_timestamp_seconds",
		"Unix timestamp from which AssetCacheManagerUtil cumulative totals apply.",
		nil,
		nil,
	)
	servedBytesTotalDesc = prometheus.NewDesc(
		"assetcache_served_bytes_total",
		"Cumulative bytes returned by Content Caching since the reported totals start time.",
		[]string{"destination"},
		nil,
	)
	storedBytesTotalDesc = prometheus.NewDesc(
		"assetcache_stored_bytes_total",
		"Cumulative bytes stored by Content Caching since the reported totals start time.",
		[]string{"source"},
		nil,
	)
	importedBytesTotalDesc = prometheus.NewDesc(
		"assetcache_imported_bytes_total",
		"Cumulative bytes imported since the reported totals start time.",
		nil,
		nil,
	)
	droppedBytesTotalDesc = prometheus.NewDesc(
		"assetcache_dropped_bytes_total",
		"Cumulative bytes downloaded but not cached since the reported totals start time.",
		nil,
		nil,
	)
	metricsLastTimestampDesc = prometheus.NewDesc(
		"assetcache_metrics_last_timestamp_seconds",
		"Unix timestamp of the latest recorded Content Caching metrics interval.",
		nil,
		nil,
	)
	metricsIntervalDesc = prometheus.NewDesc(
		"assetcache_metrics_interval_seconds",
		"Duration of the latest recorded Content Caching metrics interval.",
		nil,
		nil,
	)
	intervalServedBytesDesc = prometheus.NewDesc(
		"assetcache_interval_served_bytes",
		"Bytes served during the latest recorded metrics interval.",
		[]string{"source"},
		nil,
	)
	intervalRepliesDesc = prometheus.NewDesc(
		"assetcache_interval_replies",
		"Replies served during the latest recorded metrics interval.",
		[]string{"source"},
		nil,
	)
	intervalRequestsDesc = prometheus.NewDesc(
		"assetcache_interval_requests",
		"Requests received during the latest recorded metrics interval.",
		[]string{"source"},
		nil,
	)
	intervalImportedBytesDesc = prometheus.NewDesc(
		"assetcache_interval_imported_bytes",
		"Bytes imported during the latest recorded metrics interval.",
		nil,
		nil,
	)
	intervalImportsDesc = prometheus.NewDesc(
		"assetcache_interval_imports",
		"Import requests received during the latest recorded metrics interval.",
		nil,
		nil,
	)
	intervalDroppedBytesDesc = prometheus.NewDesc(
		"assetcache_interval_dropped_bytes",
		"Bytes downloaded but not cached during the latest recorded metrics interval.",
		nil,
		nil,
	)
	intervalPurgedBytesDesc = prometheus.NewDesc(
		"assetcache_interval_purged_bytes",
		"Bytes purged during the latest recorded metrics interval.",
		nil,
		nil,
	)
	intervalRejectedRequestsDesc = prometheus.NewDesc(
		"assetcache_interval_requests_rejected_no_space",
		"Requests rejected for lack of cache space during the latest recorded metrics interval.",
		nil,
		nil,
	)
	intervalCachePressureDesc = prometheus.NewDesc(
		"assetcache_interval_cache_pressure_ratio",
		"Cache pressure during the latest recorded metrics interval, from 0 to 1.",
		nil,
		nil,
	)
)

var allDescs = []*prometheus.Desc{
	scrapeSuccessDesc,
	scrapeDurationDesc,
	statusUpDesc,
	metricsDBUpDesc,
	activatedDesc,
	activeDesc,
	statusInfoDesc,
	cacheUsedBytesDesc,
	cacheLimitBytesDesc,
	cacheFreeBytesDesc,
	personalCacheUsedBytesDesc,
	personalCacheLimitBytesDesc,
	personalCacheFreeBytesDesc,
	cachePressureDesc,
	parentCachesDesc,
	peerCachesDesc,
	totalsStartTimestampDesc,
	servedBytesTotalDesc,
	storedBytesTotalDesc,
	importedBytesTotalDesc,
	droppedBytesTotalDesc,
	metricsLastTimestampDesc,
	metricsIntervalDesc,
	intervalServedBytesDesc,
	intervalRepliesDesc,
	intervalRequestsDesc,
	intervalImportedBytesDesc,
	intervalImportsDesc,
	intervalDroppedBytesDesc,
	intervalPurgedBytesDesc,
	intervalRejectedRequestsDesc,
	intervalCachePressureDesc,
}
