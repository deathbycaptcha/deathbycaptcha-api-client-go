# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [4.7.0] - 2026-04-01

### Added

#### Core library (`deathbycaptcha/`)
- `HttpClient` — full-featured REST client over HTTPS.
  - Authentication via `username`/`password` or `authtoken`.
  - Methods: `GetUser`, `GetBalance`, `GetCaptcha`, `GetText`, `Report`, `Upload`, `Decode`, `Close`, `GetStatus`.
  - Correct `User-Agent` header (`DBC/Go v4.7.0`) on every request.
- `SocketClient` — persistent TCP socket client.
  - Auto-reconnect with mutex-guarded connection lifecycle.
  - Same `Client` interface as `HttpClient`; thread-safe for concurrent goroutines.
- `Client` interface — unified abstraction covering both transports.
- `TokenParams(captchaType int, extraFields map[string]interface{}) map[string]interface{}` — convenience builder for token-based CAPTCHA types.
- Polling schedule `[1, 1, 2, 3, 2, 2, 3, 2, 2]` seconds (then fixed 3 s) matching all other official DBC clients.
- `AccessDeniedException`, `ServiceOverloadException` — typed sentinel errors.
- Support for all 21 CAPTCHA types:
  | Type | Name |
  |------|------|
  | 0 | Image CAPTCHA |
  | 1 | reCAPTCHA v2 |
  | 3 | Fun CAPTCHA |
  | 4 | reCAPTCHA v3 |
  | 5 | hCaptcha |
  | 6 | KeyCaptcha |
  | 7 | GeeTest |
  | 8 | GeeTest v4 |
  | 9 | Capy Puzzle |
  | 10 | CloudFlare Turnstile |
  | 11 | Amazon WAF |
  | 12 | Cyber SiARA |
  | 13 | MT Captcha |
  | 14 | Friendly Captcha |
  | 15 | Cutcaptcha |
  | 16 | Tencent CAPTCHA |
  | 17 | atbCAPTCHA |
  | 18 | Lemin Cropped CAPTCHA |
  | 19 | Arkose Labs FunCaptcha |
  | 20 | Imperva / Incapsula |
  | 21 | DataDome |

#### Tests
- Unit tests with ≥ 80 % line coverage (`go test ./deathbycaptcha/...`).
- HTTP transport tests using `httptest.NewServer` (no live credentials required).
- Socket transport tests using a local `net.Listener` mock.
- Integration tests in `tests/integration/` gated behind `DBC_USERNAME` / `DBC_PASSWORD`
  environment variables; skipped automatically when credentials are absent
  (`-tags=integration`).
- JPEG fixture at `tests/fixtures/test.jpg` for image CAPTCHA integration test.

#### Examples (`examples/`)
- 22 ready-to-run `.go` programs covering every CAPTCHA type plus `get_balance`.
- Each example reads credentials from `DBC_USERNAME` / `DBC_PASSWORD` environment
  variables (or `DBC_AUTHTOKEN`) — no hardcoded secrets.

#### CI / GitHub Actions (`.github/workflows/`)
- `unit-tests-go125.yml` — matrix unit-test run on Go 1.25.
- `unit-tests-go126.yml` — matrix unit-test run on Go 1.26.
- `coverage.yml` — generates and publishes a coverage badge to GitHub Pages.
- `integration.yml` — integration tests; skipped when repository secrets are absent.
- `publish.yml` — triggers pkg.go.dev index refresh on GitHub release.

#### Documentation & project files
- `README.md` — full usage guide: quick-start, per-type code snippets, FAQ sections
  for reCAPTCHA v2/v3, Amazon WAF, and Cloudflare Turnstile; mirrors structure of
  all other official DBC client READMEs.
- `LICENSE` — MIT.
- `RESPONSIBLE_USE.md` — usage policy.
- `.gitignore` — Go-specific ignore rules (binaries, coverage artefacts, `go.work`,
  vendor, secrets).
