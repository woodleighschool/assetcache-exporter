# AGENTS.md

Repository guidance for assetcache-exporter.

## Approach

- Stay within the requested scope and preserve unrelated local changes.
- This is a small local exporter, not a monitoring platform. Prefer direct code and a small dependency surface.
- Do not add remote targets, polling, persistence, dynamic metric generation, or configuration machinery without a concrete use case.
- Follow the shared Woodstar tooling baseline while keeping the relaxed Go lint profile appropriate for a small service.

## Repository Map

- Process composition: `cmd/assetcache_exporter`
- Apple status and metrics database collection: `internal/assetcache`
- HTTP surface: `internal/exporter`
- LaunchDaemon resources: `packaging`
- Grafana dashboard: `dashboard.yaml`

Keep the two Apple sources independently fallible and read `Metrics.db` strictly read-only.

## Commands

Use Mise tasks as the repository contract.

- Dependencies: `mise run deps`
- Build: `mise run build`
- Tests: `mise run test`
- Lint: `mise run lint`; fixes: `mise run lint-fix`
- Format: `mise run format`; check: `mise run fmt-check`
- Module and workflow checks: `mise run tidy-check`, `mise run workflow-lint`

## Engineering Rules

- Prefer concrete Go types, small consumer-owned interfaces, and wrapped errors.
- Keep metric names and labels stable unless the requested change deliberately changes the scrape contract.
- Treat `Metrics.db` rows as interval gauges. Do not expose them as Prometheus counters.
- Tests use synthetic JSON and temporary SQLite databases; they must not call live Macs.
- Keep local addresses, cache GUIDs, and other machine identifiers out of logs, fixtures, and version control.

## Commits

- Use focused Conventional Commits.
- Do not push, publish releases, or change Content Caching state unless explicitly requested.
- Report checks run, skipped checks, and unresolved failures.
