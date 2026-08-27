# assetcache-exporter

[![Release](https://img.shields.io/github/v/release/woodleighschool/assetcache-exporter?display_name=tag&sort=semver)](https://github.com/woodleighschool/assetcache-exporter/releases/latest)
[![CI](https://github.com/woodleighschool/assetcache-exporter/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/woodleighschool/assetcache-exporter/actions/workflows/ci.yaml)
[![Go](https://img.shields.io/github/go-mod/go-version/woodleighschool/assetcache-exporter?logo=go)](https://github.com/woodleighschool/assetcache-exporter/blob/main/go.mod)
[![License](https://img.shields.io/github/license/woodleighschool/assetcache-exporter)](https://github.com/woodleighschool/assetcache-exporter/blob/main/LICENSE)

Prometheus exporter for Apple Content Caching. It reads current status from `AssetCacheManagerUtil` and Apple's local metrics database, then serves both at `/metrics`.

Each scrape reads the two Apple sources independently, so one can fail without hiding metrics from the other.

## 🚀 Usage

Download the macOS `.pkg` from the [latest release](https://github.com/woodleighschool/assetcache-exporter/releases/latest). It installs `assetcache_exporter` in `/usr/local/bin` and a system LaunchDaemon in `/Library/LaunchDaemons`. The service runs as the macOS `_assetcache` account on `:9200`; installation loads it immediately, and launchd keeps it running.

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

`mise run build` writes `./assetcache_exporter`. Running that binary directly does not install or load the LaunchDaemon.

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
