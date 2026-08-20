# assetcache-exporter

Prometheus exporter for Apple Content Caching. It reads current status from `AssetCacheManagerUtil` and Apple's local metrics database, then serves both at `/metrics`.

It does not change Content Caching, poll in the background, or keep its own state. The two Apple sources are read independently, so one can fail without hiding metrics from the other.

## 🚀 Usage

Download the macOS `.pkg` attached to the [latest release](https://github.com/woodleighschool/assetcache-exporter/releases/latest) and open it. The package installs a LaunchDaemon, listens on `:9200`, and runs as the local `_assetcache` account.

To run from source:

```bash
mise run build
./assetcache_exporter
```

Prometheus can scrape the Mac directly:

```yaml
scrape_configs:
  - job_name: assetcache
    static_configs:
      - targets: [cache.example.edu:9200]
```

A Grafana dashboard is available in [`dashboard.json`](dashboard.json).

## ⚙️ Configuration

| Flag                   | Default    | Purpose                            |
| ---------------------- | ---------- | ---------------------------------- |
| `--web.listen-address` | `:9200`    | HTTP listen address                |
| `--web.telemetry-path` | `/metrics` | Metrics path                       |
| `--collector.timeout`  | `5s`       | Timeout for each Apple data source |
| `--version`            |            | Print build information            |

## 📈 Metrics

| Metric                                       | Meaning                                                   |
| -------------------------------------------- | --------------------------------------------------------- |
| `assetcache_active`                          | Whether Content Caching is active                         |
| `assetcache_cache_used_bytes`                | Logical bytes held in the cache                           |
| `assetcache_cache_limit_bytes`               | Configured cache limit                                    |
| `assetcache_cache_pressure_ratio`            | Maximum cache pressure during the last hour               |
| `assetcache_served_bytes_total{destination}` | Cumulative bytes returned to clients, children, and peers |
| `assetcache_stored_bytes_total{source}`      | Cumulative bytes stored from origin, parents, and peers   |
| `assetcache_interval_served_bytes{source}`   | Bytes served in the latest Apple metrics interval         |
| `assetcache_interval_requests{source}`       | Requests received in the latest interval                  |
| `assetcache_metrics_last_timestamp_seconds`  | Timestamp of the latest interval                          |
| `assetcache_status_up`                       | Whether `AssetCacheManagerUtil` was readable              |
| `assetcache_metrics_db_up`                   | Whether `Metrics.db` was readable and compatible          |
| `assetcache_scrape_success`                  | Whether both sources succeeded                            |

Interval observations are gauges. Apple can omit zero rows while idle, so an old interval timestamp does not itself mean collection failed.

The exporter opens `/Library/Application Support/Apple/AssetCache/Metrics/Metrics.db` read-only. An incompatible schema sets `assetcache_metrics_db_up` to zero while the HTTP server and status collection continue.

## 🧑‍💻 Development

```bash
mise run deps
mise run build
mise run test
mise run lint
mise run fmt-check
mise run workflow-lint
```

## 📄 License

Licensed under the [Apache License 2.0](LICENSE).
