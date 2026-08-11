# assetcache-exporter

`assetcache-exporter` is a small Prometheus exporter for Apple Content Caching. It runs on the cache Mac, reads current status from `AssetCacheManagerUtil`, reads Apple's local metrics database in SQLite read-only mode, and exposes both at `/metrics`.

The exporter does not change Content Caching settings, poll in the background, or keep its own state. A failure in either Apple data source is reported without suppressing metrics from the other source.

## Install

Download the `.pkg` from the latest [release build](https://github.com/woodleighschool/assetcache-exporter/actions/workflows/release.yaml) and open it. The package installs and starts the exporter as a LaunchDaemon.

The exporter listens on `:9200` and serves metrics at `/metrics`.

## Run from source

```sh
mise run build
./assetcache-exporter
```

It listens on `:9200` by default. Common flags are:

```text
--web.listen-address=:9200
--web.telemetry-path=/metrics
--collector.timeout=5s
--version
```

Prometheus can scrape the Mac with a normal static target or `ScrapeConfig`:

```yaml
scrape_configs:
    - job_name: assetcache
      static_configs:
          - targets: [cache.example.edu:9200]
```

## Metrics

The main metrics are:

| Metric                                       | Meaning                                                     |
| -------------------------------------------- | ----------------------------------------------------------- |
| `assetcache_active`                          | Whether Content Caching is currently active.                |
| `assetcache_cache_used_bytes`                | Logical bytes currently held in the cache.                  |
| `assetcache_cache_limit_bytes`               | Configured cache limit.                                     |
| `assetcache_cache_pressure_ratio`            | Maximum cache pressure during the last hour.                |
| `assetcache_served_bytes_total{destination}` | Cumulative bytes returned to clients, children, and peers.  |
| `assetcache_stored_bytes_total{source}`      | Cumulative bytes stored from origin, parents, and peers.    |
| `assetcache_interval_served_bytes{source}`   | Bytes served in the latest recorded Apple metrics interval. |
| `assetcache_interval_requests{source}`       | Requests received in the latest recorded interval.          |
| `assetcache_metrics_last_timestamp_seconds`  | Timestamp of that latest interval.                          |
| `assetcache_status_up`                       | Whether `AssetCacheManagerUtil` was readable.               |
| `assetcache_metrics_db_up`                   | Whether `Metrics.db` was readable and compatible.           |
| `assetcache_scrape_success`                  | Whether both sources succeeded.                             |

Interval observations are gauges, not counters. Apple omits zero rows while the cache is idle, so an old `assetcache_metrics_last_timestamp_seconds` does not by itself mean collection has failed. Use `assetcache_metrics_db_up` to distinguish a readable idle database from a source failure.

## Dashboard

A Grafana dashboard is available at [`dashboard.yaml`](dashboard.yaml).

## Data sources

The exporter reads:

```text
/usr/bin/AssetCacheManagerUtil -j status
/Library/Application Support/Apple/AssetCache/Metrics/Metrics.db
```

The SQLite connection uses `mode=ro`. Apple documents `Metrics.db` as an implementation surface that may change between macOS releases; incompatible schemas set `assetcache_metrics_db_up` to zero while the HTTP server and status collection continue running.

## Development

Use the Mise tasks as the repository interface:

```sh
mise run deps
mise run build
mise run test
mise run lint
mise run fmt-check
mise run workflow-lint
```
