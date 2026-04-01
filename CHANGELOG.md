# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-04-01

### Added
- Initial release of the official DeathByCaptcha Go client.
- `HttpClient` — REST over HTTPS with `username`/`password` and `authtoken` support.
- `SocketClient` — persistent TCP socket client with auto-reconnect and thread-safety.
- `Client` interface with `GetUser`, `GetBalance`, `GetCaptcha`, `GetText`, `Report`, `Upload`, `Decode`, `Close`.
- `GetStatus` helper on `HttpClient` for service health checks.
- `TokenParams` convenience helper for token CAPTCHA parameter maps.
- Full polling schedule (`[1,1,2,3,2,2,3,2,2]` then 3 s fixed).
- Support for all 21 CAPTCHA types (types 0–25).
- Unit tests with ≥ 80 % line coverage (`go test ./deathbycaptcha/...`).
- Integration tests gated behind `DBC_USERNAME`/`DBC_PASSWORD` environment variables (`-tags=integration`).
- 22 usage examples in `examples/`.
- GitHub Actions workflows: unit tests (Go 1.25 & 1.26), coverage badge, integration, and pkg.go.dev publish.
- `LICENSE` (MIT), `RESPONSIBLE_USE.md`, `CHANGELOG.md`, and `README.md`.
