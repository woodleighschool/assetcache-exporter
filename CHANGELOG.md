# Changelog

## [1.2.0](https://github.com/woodleighschool/assetcache-exporter/compare/1.1.0...1.2.0) (2026-08-27)


### Features

* **go:** update prometheus group ([#15](https://github.com/woodleighschool/assetcache-exporter/issues/15)) ([c523af7](https://github.com/woodleighschool/assetcache-exporter/commit/c523af747375903d473ed8f71eb052affb298fc5))


### Code Refactoring

* adopt exporter toolkit bootstrap ([18a9a64](https://github.com/woodleighschool/assetcache-exporter/commit/18a9a64872a1b49780a56711dad9f10d0b299dc8))


### Documentation

* align repository agent guidance ([68a2050](https://github.com/woodleighschool/assetcache-exporter/commit/68a2050006a6614758014ba5ffbae60904fc8d8a))
* clarify repository guidance ([5f91826](https://github.com/woodleighschool/assetcache-exporter/commit/5f91826fb639dd4fc0757ccdf9fa8edd52380565))
* clarify usage and releases ([3a594f5](https://github.com/woodleighschool/assetcache-exporter/commit/3a594f5770637af5f7a999d10269630dd88165a8))
* tighten exporter description ([3128460](https://github.com/woodleighschool/assetcache-exporter/commit/31284600b3abee4a54e521697665e8cdf1dae575))


### Build System

* target macOS arm64 ([c4b1d57](https://github.com/woodleighschool/assetcache-exporter/commit/c4b1d57a76b675a0af5156f836bb874b42c2224e))


### Continuous Integration

* **github-action:** Update action home-operations/.github/actions/workflow-lint (v1.0.2 → v1.0.3) ([#8](https://github.com/woodleighschool/assetcache-exporter/issues/8)) ([1fe2fd8](https://github.com/woodleighschool/assetcache-exporter/commit/1fe2fd8ebb7da675ceb173b070961378a50fe5e7))
* **release:** restore GoReleaser ([c951d93](https://github.com/woodleighschool/assetcache-exporter/commit/c951d93f7b874e6e6f9b747ae8285fd6ac021931))
* **release:** use shared notarization action ([5515788](https://github.com/woodleighschool/assetcache-exporter/commit/55157880786577d8489111edd7e4e966a48e5ce7))
* **release:** use shared package action ([f00b3f3](https://github.com/woodleighschool/assetcache-exporter/commit/f00b3f313ac193aa880fc0d41c7d361359d6b0cf))
* sync shared repository tooling ([8df0f63](https://github.com/woodleighschool/assetcache-exporter/commit/8df0f63a4da9170fe92a94a1bf682fd6556366cf))
* use shared App Store Connect key ([fb8b150](https://github.com/woodleighschool/assetcache-exporter/commit/fb8b15062a00c21e4a2a85628a74b236d30d7597))
* verify release package ([5c6d55c](https://github.com/woodleighschool/assetcache-exporter/commit/5c6d55cbc92cc32cd4f81b4fdbf5bbccb05e7250))


### Miscellaneous Chores

* align ignore rules ([8f856fd](https://github.com/woodleighschool/assetcache-exporter/commit/8f856fd779fb3ae166aa4d15e17870a672befb2e))
* cleanup grafana dashboard ([c099a62](https://github.com/woodleighschool/assetcache-exporter/commit/c099a622ac435ad8f9c2decd1402052e42cd5f5d))
* **go:** update toolchain to 1.27 ([5c6cfad](https://github.com/woodleighschool/assetcache-exporter/commit/5c6cfadd2c7b058dcb8ad40d2a134a5a3fee8397))
* **mise:** update tool go (1.26.6 → 1.27.0) ([#9](https://github.com/woodleighschool/assetcache-exporter/issues/9)) ([76e360b](https://github.com/woodleighschool/assetcache-exporter/commit/76e360b18187c61daf5c6f14046759c9f155fe20))
* **mise:** update tool golangci-lint (2.12.2 → 2.13.0) ([#10](https://github.com/woodleighschool/assetcache-exporter/issues/10)) ([cce94bc](https://github.com/woodleighschool/assetcache-exporter/commit/cce94bc07a0f7c47d518272a166ca74113eb6677))
* **mise:** update tool golangci-lint (2.13.0 → 2.13.1) ([#11](https://github.com/woodleighschool/assetcache-exporter/issues/11)) ([4da6552](https://github.com/woodleighschool/assetcache-exporter/commit/4da6552d279a79e35c87aac867bab4f4faa1594a))
* **mise:** update tool lefthook (2.1.10 → 2.1.11) ([#12](https://github.com/woodleighschool/assetcache-exporter/issues/12)) ([80b774f](https://github.com/woodleighschool/assetcache-exporter/commit/80b774ff788742e2e92a0890b5043657c15bdc10))
* **mise:** update tool oxfmt (0.64.0 → 0.65.0) ([#14](https://github.com/woodleighschool/assetcache-exporter/issues/14)) ([89ef3b6](https://github.com/woodleighschool/assetcache-exporter/commit/89ef3b63db1a0343a99abfcadf52efe7db8dfe57))
* **release-please:** sync configuration ([9727afc](https://github.com/woodleighschool/assetcache-exporter/commit/9727afc1c93c0faad4e3222685af5e0bf39c3a53))

## [1.1.0](https://github.com/woodleighschool/assetcache-exporter/compare/1.0.1...1.1.0) (2026-08-20)


### Features

* **go:** update module modernc.org/sqlite (v1.56.0 → v1.57.0) ([#6](https://github.com/woodleighschool/assetcache-exporter/issues/6)) ([6ccdb2b](https://github.com/woodleighschool/assetcache-exporter/commit/6ccdb2b84108809a67f3699ab1dc20722029ea26))


### Bug Fixes

* **deps:** update indirect dependencies ([a70f346](https://github.com/woodleighschool/assetcache-exporter/commit/a70f346874406dcdc217ba64d5be85345d6bf59c))
* preserve interval source series ([1b1d5f9](https://github.com/woodleighschool/assetcache-exporter/commit/1b1d5f9f8f14a5e366c60b3f21cf545ffc61877d))

## [1.0.1](https://github.com/woodleighschool/assetcache-exporter/compare/1.0.0...1.0.1) (2026-08-11)


### Bug Fixes

* publish package with releases ([3d2199c](https://github.com/woodleighschool/assetcache-exporter/commit/3d2199c074a173e00b0a2d41722684d33b834e21))
* upload package with gh ([a1aea72](https://github.com/woodleighschool/assetcache-exporter/commit/a1aea7296f196e425d935a417abe966400630490))
* use JSON for Grafana dashboard ([c4c90fe](https://github.com/woodleighschool/assetcache-exporter/commit/c4c90fe46ff4530fca05f7abb58d5412ddf0099f))

## 1.0.0 (2026-08-11)


### Features

* add Apple content cache exporter ([8017e97](https://github.com/woodleighschool/assetcache-exporter/commit/8017e97a61f3fe2752f7ee5d49bede77ee78e6db))
* add macOS installer package ([14247ce](https://github.com/woodleighschool/assetcache-exporter/commit/14247ce19f02219b3cedbef55841fc691f53c56b))

## Changelog
