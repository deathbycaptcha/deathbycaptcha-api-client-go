# [DeathByCaptcha](https://deathbycaptcha.com/)


<p align="center">
  <a href="https://github.com/deathbycaptcha/deathbycaptcha-api-client-python"><img alt="Python" src="https://img.shields.io/badge/Python-3776AB?style=for-the-badge&logo=python&logoColor=white"></a>
  <a href="https://github.com/deathbycaptcha/deathbycaptcha-api-client-nodejs"><img alt="Node.js" src="https://img.shields.io/badge/Node.js-339933?style=for-the-badge&logo=nodedotjs&logoColor=white"></a>
  <a href="https://github.com/deathbycaptcha/deathbycaptcha-api-client-dotnet"><img alt=".NET" src="https://img.shields.io/badge/.NET-512BD4?style=for-the-badge&logo=dotnet&logoColor=white"></a>
  <a href="https://github.com/deathbycaptcha/deathbycaptcha-api-client-java"><img alt="Java" src="https://img.shields.io/badge/Java-ED8B00?style=for-the-badge&logo=openjdk&logoColor=white"></a>
  <a href="https://github.com/deathbycaptcha/deathbycaptcha-api-client-php"><img alt="PHP" src="https://img.shields.io/badge/PHP-777BB4?style=for-the-badge&logo=php&logoColor=white"></a>
  <a href="https://github.com/deathbycaptcha/deathbycaptcha-api-client-perl"><img alt="Perl" src="https://img.shields.io/badge/Perl-39457E?style=for-the-badge&logo=perl&logoColor=white"></a>
  <a href="https://github.com/deathbycaptcha/deathbycaptcha-api-client-c11"><img alt="C" src="https://img.shields.io/badge/C-A8B9CC?style=for-the-badge&logo=c&logoColor=black"></a>
  <a href="https://github.com/deathbycaptcha/deathbycaptcha-api-client-cpp"><img alt="C++" src="https://img.shields.io/badge/C%2B%2B-00599C?style=for-the-badge&logo=cplusplus&logoColor=white"></a>
  <a href="https://github.com/deathbycaptcha/deathbycaptcha-api-client-go"><img alt="Go" src="https://img.shields.io/badge/%E2%80%BA-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white&labelColor=555555"></a>
</p>


## 📖 Introduction

The [DeathByCaptcha](https://deathbycaptcha.com) Go client is the official library for the DeathByCaptcha **captcha solving service**. It provides a simple, well-documented interface for integrating CAPTCHA solving into automation workflows — a particularly common need when you use it as a **captcha solver for web scraping**, where CAPTCHAs block access to the pages you need to extract data from. It supports both the HTTPS API (encrypted transport — recommended when security is a priority) and the socket-based API (faster and lower latency, recommended for high-throughput production workloads).

Compatible with **Go 1.21+** (tested on Go 1.25 and Go 1.26).

Key features:

- 🧩 Send image, audio and modern token-based CAPTCHA types (reCAPTCHA v2/v3, Turnstile, GeeTest, etc.).
- 🔄 Unified `Client` interface across HTTP and socket transports — switching implementations is a one-line change.
- 🔐 Built-in support for proxies, timeouts and advanced token parameters for modern CAPTCHA flows.
- ⚡ Concurrent-safe: both clients use `sync.Mutex` internally.

Quick start example (HTTP):

```go
import dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/v4/deathbycaptcha"

client := dbc.NewHttpClient("your_username", "your_password")
defer client.Close()

captcha, err := client.Decode("path/to/captcha.jpg", dbc.DefaultTimeout, map[string]string{})
if err != nil {
    log.Fatal(err)
}
if captcha != nil {
    fmt.Println("Solved:", *captcha.Text)
}
```

> **🚌 Transport options:** Use `HttpClient` for encrypted HTTPS communication — credentials and data travel over TLS. Use `SocketClient` for lower latency and higher throughput — it is faster but communicates over a plain TCP connection to `api.dbcapi.me` on ports `8123–8130`.

---

### Tests Status

[![Unit Tests Go 1.25](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/actions/workflows/unit-tests-go125.yml/badge.svg)](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/actions/workflows/unit-tests-go125.yml)
[![Unit Tests Go 1.26](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/actions/workflows/unit-tests-go126.yml/badge.svg)](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/actions/workflows/unit-tests-go126.yml)
[![API Integration](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/actions/workflows/integration.yml/badge.svg)](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/actions/workflows/integration.yml)
[![Coverage](https://img.shields.io/endpoint?url=https://deathbycaptcha.github.io/deathbycaptcha-api-client-go/coverage-badge.json)](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/actions/workflows/coverage.yml)

---

## 🗂️ Index

- [Installation](#installation)
    - [From pkg.go.dev (Recommended)](#from-pkggodev-recommended)
    - [From GitHub Repository](#from-github-repository)
- [How to Use DBC API Clients](#how-to-use-dbc-api-clients)
    - [Common Clients' Interface](#common-clients-interface)
    - [Available Methods](#captcha-methods)
- [Credentials & Configuration](#credentials--configuration)
    - [Quick Setup](#quick-setup)
- [CAPTCHA Types Quick Reference & Examples](#captcha-types-reference)
    - [Quick Start](#quick-start)
    - [Type Reference](#sample-index-by-captcha-type)
    - [Per-Type Code Snippets](#quick-type-snippets)
- [CAPTCHA Types Extended Reference](#captcha-types-extended-reference)
    - [reCAPTCHA Image-Based API — Deprecated (Types 2 & 3)](#recaptcha-image-based-api)
    - [reCAPTCHA Token API (v2 & v3)](#recaptcha-token-api)
    - [reCAPTCHA v2 API FAQ](#recaptcha-v2-api-faq)
    - [What is reCAPTCHA v3?](#what-is-recaptcha-v3)
    - [reCAPTCHA v3 API FAQ](#recaptcha-v3-api-faq)
    - [Amazon WAF API (Type 16)](#amazon-waf-api-faq)
    - [Cloudflare Turnstile API (Type 12)](#cloudflare-turnstile-api-faq)


<a id="installation"></a>
## 🛠️ Installation

<a id="from-pkggodev-recommended"></a>
### 📦 From pkg.go.dev (Recommended)

```bash
go get github.com/deathbycaptcha/deathbycaptcha-api-client-go/v4/deathbycaptcha
```

<a id="from-github-repository"></a>
### 🐙 From GitHub Repository

For the latest development version or to contribute:

```bash
git clone https://github.com/deathbycaptcha/deathbycaptcha-api-client-go.git
cd deathbycaptcha-api-client-go
go mod download
```

<a id="how-to-use-dbc-api-clients"></a>
## 🚀 How to Use DBC API Clients

<a id="common-clients-interface"></a>
### 🔌 Common Clients' Interface

All clients must be instantiated with your DeathByCaptcha credentials — either *username* and *password*, or an *authtoken* (available in the DBC user panel). Replace `NewHttpClient` with `NewSocketClient` to use the socket transport instead.

```go
import dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/v4/deathbycaptcha"

// Username + password (HTTPS transport — encrypted, recommended when security matters)
client := dbc.NewHttpClient(username, password)

// Username + password (socket transport — faster, lower latency, recommended for high throughput)
// client := dbc.NewSocketClient(username, password)

// Authtoken (HTTPS transport)
client := dbc.NewHttpClientWithToken(authtoken)

// Authtoken (socket transport)
// client := dbc.NewSocketClientWithToken(authtoken)
```

| Transport | Constructor | Best for |
|---|---|---|
| HTTPS | `dbc.NewHttpClient` | Encrypted TLS transport — safer for credential handling and network-sensitive environments |
| Socket | `dbc.NewSocketClient` | Plain TCP — faster and lower latency, recommended for high-throughput production workloads |

All clients share the same `Client` interface. Below is a summary of every available method and its signature.

<a id="captcha-methods"></a>

| Method | Signature | Returns | Description |
|---|---|---|---|
| `Upload` | `Upload(captcha interface{}, params map[string]string)` | `(*Captcha, error)` | Upload a CAPTCHA for solving without waiting. `captcha` is a file path or `nil`; `params` carry type-specific fields (`type`, `token_params`, `waf_params`, etc.). |
| `Decode` | `Decode(captcha interface{}, timeout int, params map[string]string)` | `(*Captcha, error)` | Upload and poll until solved or timed out. Preferred method for most integrations. |
| `GetCaptcha` | `GetCaptcha(id int)` | `(*Captcha, error)` | Fetch status and result of a previously uploaded CAPTCHA by its numeric ID. |
| `GetText` | `GetText(id int)` | `(string, error)` | Convenience wrapper — return only the solved text from `GetCaptcha`. |
| `Report` | `Report(id int)` | `(bool, error)` | Report a CAPTCHA as incorrectly solved to request a refund. Only report genuine errors. |
| `GetBalance` | `GetBalance()` | `(float64, error)` | Return the current account balance in US cents. |
| `GetUser` | `GetUser()` | `(*User, error)` | Return user info: ID, rate, balance, banned status. |
| `Close` | `Close()` | — | Release resources (socket connection). Safe to call on `HttpClient` (no-op). |

### 📬 CAPTCHA Result Object

All methods that return a solved CAPTCHA return a `*Captcha` struct with the following fields:

| Field | Type | Description |
|---|---|---|
| `CaptchaID` | `int` | Numeric CAPTCHA ID assigned by DBC |
| `Text` | `*string` | Solved text or token (the value you inject into the page); `nil` if not yet solved |
| `IsCorrect` | `bool` | Whether DBC considers the solution correct |

```go
// Example result
captcha, err := client.Decode("captcha.jpg", dbc.DefaultTimeout, map[string]string{})
if captcha != nil {
    fmt.Printf("ID: %d  Text: %s  Correct: %v\n",
        captcha.CaptchaID, *captcha.Text, captcha.IsCorrect)
}
```

### 💡 Full Usage Example

```go
package main

import (
    "fmt"
    "log"
    "os"

    dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/v4/deathbycaptcha"
)

func main() {
    client := dbc.NewHttpClient(os.Getenv("DBC_USERNAME"), os.Getenv("DBC_PASSWORD"))
    defer client.Close()

    bal, err := client.GetBalance()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Balance: %.4f US cents\n", bal)

    captcha, err := client.Decode("path/to/captcha.jpg", dbc.DefaultTimeout, map[string]string{})
    if err != nil {
        if _, ok := err.(*dbc.AccessDeniedException); ok {
            log.Fatal("Access denied — check your credentials and/or balance")
        }
        log.Fatal(err)
    }
    if captcha != nil {
        fmt.Printf("Solved CAPTCHA %d: %s\n", captcha.CaptchaID, *captcha.Text)

        // Report only if you are certain the solution is wrong:
        // client.Report(captcha.CaptchaID)
    }
}
```

<a id="credentials--configuration"></a>
## 🔑 Credentials & Configuration

<a id="quick-setup"></a>
### ⚡ Quick Setup

```go
// Username + password
client := dbc.NewHttpClient("your_username", "your_password")

// Auth token (from DBC user panel)
client := dbc.NewHttpClientWithToken("your_authtoken")

// Socket variants
client := dbc.NewSocketClient("your_username", "your_password")
client := dbc.NewSocketClientWithToken("your_authtoken")
```

Never commit real credentials to source control. Use environment variables:

```bash
# ① Export credentials in your shell
export DBC_USERNAME=your_username
export DBC_PASSWORD=your_password

# ② Run tests locally (fastest)
go test ./deathbycaptcha/... -v

# ③ Run integration tests with real credentials
DBC_USERNAME=your_user DBC_PASSWORD=your_pass \
  go test -tags=integration ./tests/integration/... -v -timeout 300s

# ④ Push to repo for GitHub Actions
git push
```

<a id="captcha-types-reference"></a>
## 🧩 CAPTCHA Types Quick Reference & Examples

This section covers every supported CAPTCHA type, how to run the corresponding example scripts, and ready-to-copy code snippets. Start with the Quick Start below, then use the Type Reference to find the type you need.

<a id="quick-start"></a>
### 🏁 Quick Start

1. **📦 Install the library** (see [Installation](#installation))
2. **📂 Navigate to the `examples/` directory** and run the script for the type you need:

```bash
cd examples
go run example.Normal_Captcha.go

# Balance check:
go run get_balance.go
```

> ⚠️ Always run examples from the `examples/` directory so the `images/` folder is accessible to scripts that need it.

Before running any script, add your DBC credentials inside it:

```go
username := "your_username"
password := "your_password"
```

<a id="sample-index-by-captcha-type"></a>
### 📋 Type Reference

The table below maps every supported type to its use case, a code snippet, and the corresponding example file in `examples/`.

| Type ID | CAPTCHA Type | Use Case | Quick Use | Go Sample |
| --- | --- | --- | --- | --- |
| 0 | Standard Image | Basic image CAPTCHA | [snippet](#sample-type-0-standard-image) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Normal_Captcha.go) |
| 2 | ~~reCAPTCHA Coordinates~~ | Deprecated — do not use for new integrations | — | — |
| 3 | ~~reCAPTCHA Image Group~~ | Deprecated — do not use for new integrations | — | — |
| 4 | reCAPTCHA v2 Token | reCAPTCHA v2 token solving | [snippet](#sample-type-4-recaptcha-v2-token) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.reCAPTCHA_v2.go) |
| 5 | reCAPTCHA v3 Token | reCAPTCHA v3 with risk scoring | [snippet](#sample-type-5-recaptcha-v3-token) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.reCAPTCHA_v3.go) |
| 25 | reCAPTCHA v2 Enterprise | reCAPTCHA v2 Enterprise tokens | [snippet](#sample-type-25-recaptcha-v2-enterprise) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.reCAPTCHA_v2_Enterprise.go) |
| 8 | GeeTest v3 | GeeTest v3 verification | [snippet](#sample-type-8-geetest-v3) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Geetest_v3.go) |
| 9 | GeeTest v4 | GeeTest v4 verification | [snippet](#sample-type-9-geetest-v4) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Geetest_v4.go) |
| 11 | Text CAPTCHA | Text-based question solving | [snippet](#sample-type-11-text-captcha) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Textcaptcha.go) |
| 12 | Cloudflare Turnstile | Cloudflare Turnstile token | [snippet](#sample-type-12-cloudflare-turnstile) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Turnstile.go) |
| 13 | Audio CAPTCHA | Audio CAPTCHA solving | [snippet](#sample-type-13-audio-captcha) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Audio.go) |
| 14 | Lemin | Lemin CAPTCHA | [snippet](#sample-type-14-lemin) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Lemin.go) |
| 15 | Capy | Capy CAPTCHA | [snippet](#sample-type-15-capy) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Capy.go) |
| 16 | Amazon WAF | Amazon WAF verification | [snippet](#sample-type-16-amazon-waf) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Amazon_Waf.go) |
| 17 | Siara | Siara CAPTCHA | [snippet](#sample-type-17-siara) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Siara.go) |
| 18 | MTCaptcha | MTCaptcha CAPTCHA | [snippet](#sample-type-18-mtcaptcha) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Mtcaptcha.go) |
| 19 | Cutcaptcha | Cutcaptcha CAPTCHA | [snippet](#sample-type-19-cutcaptcha) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Cutcaptcha.go) |
| 20 | Friendly Captcha | Friendly Captcha | [snippet](#sample-type-20-friendly-captcha) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Friendly.go) |
| 21 | DataDome | DataDome verification | [snippet](#sample-type-21-datadome) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Datadome.go) |
| 23 | Tencent | Tencent CAPTCHA | [snippet](#sample-type-23-tencent) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Tencent.go) |
| 24 | ATB | ATB CAPTCHA | [snippet](#sample-type-24-atb) | [open](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Atb.go) |

<a id="quick-type-snippets"></a>
### 📝 Per-Type Code Snippets

Minimal usage snippet for each supported type. Use these as a starting point and refer to the full example files in `examples/` for complete implementations.

<a id="sample-type-0-standard-image"></a>
#### 🖼️ Sample Type 0: Standard Image
Official description: [Supported CAPTCHAs](https://deathbycaptcha.com/api#supported_captchas)
Full sample: [example.Normal_Captcha.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Normal_Captcha.go)

```go
captcha, err := client.Decode("images/normal.jpg", dbc.DefaultTimeout, map[string]string{})
```

---

<a id="sample-type-4-recaptcha-v2-token"></a>
#### 🤖 Sample Type 4: reCAPTCHA v2 Token
Official description: [reCAPTCHA Token API (v2)](https://deathbycaptcha.com/api/newtokenrecaptcha#token-v2)
Full sample: [example.reCAPTCHA_v2.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.reCAPTCHA_v2.go)

```go
import "encoding/json"

tokenParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "googlekey": "sitekey",
    "pageurl":   "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":         "4",
    "token_params": string(tokenParams),
})
```

---

<a id="sample-type-5-recaptcha-v3-token"></a>
#### 🤖 Sample Type 5: reCAPTCHA v3 Token
Official description: [reCAPTCHA v3](https://deathbycaptcha.com/api/newtokenrecaptcha#reCAPTCHAv3)
Full sample: [example.reCAPTCHA_v3.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.reCAPTCHA_v3.go)

```go
import "encoding/json"

tokenParams, _ := json.Marshal(map[string]interface{}{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "googlekey": "sitekey",
    "pageurl":   "https://target",
    "action":    "verify",
    "min_score": 0.3,
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":         "5",
    "token_params": string(tokenParams),
})
```

---

<a id="sample-type-25-recaptcha-v2-enterprise"></a>
#### 🏢 Sample Type 25: reCAPTCHA v2 Enterprise
Official description: [reCAPTCHA v2 Enterprise](https://deathbycaptcha.com/api/newtokenrecaptcha#reCAPTCHAv2Enterprise)
Full sample: [example.reCAPTCHA_v2_Enterprise.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.reCAPTCHA_v2_Enterprise.go)

```go
import "encoding/json"

enterpriseParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "googlekey": "sitekey",
    "pageurl":   "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":                    "25",
    "token_enterprise_params": string(enterpriseParams),
})
```

---

<a id="sample-type-8-geetest-v3"></a>
#### 🧩 Sample Type 8: GeeTest v3
Official description: [GeeTest](https://deathbycaptcha.com/api/geetest)
Full sample: [example.Geetest_v3.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Geetest_v3.go)

```go
import "encoding/json"

geetestParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "gt":        "gt_value",
    "challenge": "challenge_value",
    "pageurl":   "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":           "8",
    "geetest_params": string(geetestParams),
})
```

---

<a id="sample-type-9-geetest-v4"></a>
#### 🧩 Sample Type 9: GeeTest v4
Official description: [GeeTest](https://deathbycaptcha.com/api/geetest)
Full sample: [example.Geetest_v4.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Geetest_v4.go)

```go
import "encoding/json"

geetestParams, _ := json.Marshal(map[string]string{
    "proxy":      "http://user:pass@127.0.0.1:1234",
    "proxytype":  "HTTP",
    "captcha_id": "captcha_id",
    "pageurl":    "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":           "9",
    "geetest_params": string(geetestParams),
})
```

---

<a id="sample-type-11-text-captcha"></a>
#### 💬 Sample Type 11: Text CAPTCHA
Official description: [Text CAPTCHA](https://deathbycaptcha.com/api/textcaptcha)
Full sample: [example.Textcaptcha.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Textcaptcha.go)

```go
captcha, err := client.Decode(nil, dbc.DefaultTimeout, map[string]string{
    "type":        "11",
    "textcaptcha": "What is two plus two?",
})
```

---

<a id="sample-type-12-cloudflare-turnstile"></a>
#### ☁️ Sample Type 12: Cloudflare Turnstile
Official description: [Cloudflare Turnstile](https://deathbycaptcha.com/api/turnstile)
Full sample: [example.Turnstile.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Turnstile.go)

```go
import "encoding/json"

turnstileParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "sitekey":   "sitekey",
    "pageurl":   "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":             "12",
    "turnstile_params": string(turnstileParams),
})
```

---

<a id="sample-type-13-audio-captcha"></a>
#### 🔊 Sample Type 13: Audio CAPTCHA
Official description: [Audio CAPTCHA](https://deathbycaptcha.com/api/audio)
Full sample: [example.Audio.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Audio.go)

```go
import (
    "encoding/base64"
    "os"
)

audioData, _ := os.ReadFile("images/audio.mp3")
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":     "13",
    "audio":    base64.StdEncoding.EncodeToString(audioData),
    "language": "en",
})
```

---

<a id="sample-type-14-lemin"></a>
#### 🔵 Sample Type 14: Lemin
Official description: [Lemin](https://deathbycaptcha.com/api/lemin)
Full sample: [example.Lemin.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Lemin.go)

```go
import "encoding/json"

leminParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "captchaid": "CROPPED_xxx",
    "pageurl":   "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":         "14",
    "lemin_params": string(leminParams),
})
```

---

<a id="sample-type-15-capy"></a>
#### 🏴 Sample Type 15: Capy
Official description: [Capy](https://deathbycaptcha.com/api/capy)
Full sample: [example.Capy.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Capy.go)

```go
import "encoding/json"

capyParams, _ := json.Marshal(map[string]string{
    "proxy":      "http://user:pass@127.0.0.1:1234",
    "proxytype":  "HTTP",
    "captchakey": "PUZZLE_xxx",
    "api_server": "https://api.capy.me/",
    "pageurl":    "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":        "15",
    "capy_params": string(capyParams),
})
```

---

<a id="sample-type-16-amazon-waf"></a>
#### 🛡️ Sample Type 16: Amazon WAF
Official description: [Amazon WAF](https://deathbycaptcha.com/api/amazonwaf)
Full sample: [example.Amazon_Waf.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Amazon_Waf.go)

```go
import "encoding/json"

wafParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "sitekey":   "sitekey",
    "pageurl":   "https://target",
    "iv":        "iv_value",
    "context":   "context_value",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":       "16",
    "waf_params": string(wafParams),
})
```

---

<a id="sample-type-17-siara"></a>
#### 🔍 Sample Type 17: Siara
Official description: [Siara](https://deathbycaptcha.com/api/siara)
Full sample: [example.Siara.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Siara.go)

```go
import "encoding/json"

siaraParams, _ := json.Marshal(map[string]string{
    "proxy":      "http://user:pass@127.0.0.1:1234",
    "proxytype":  "HTTP",
    "slideurlid": "OXR2LVNvCuXykkZbB8KZIfh162sNT8S2",
    "pageurl":    "https://target",
    "useragent":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":         "17",
    "siara_params": string(siaraParams),
})
```

---

<a id="sample-type-18-mtcaptcha"></a>
#### 🔒 Sample Type 18: MTCaptcha
Official description: [MTCaptcha](https://deathbycaptcha.com/api/mtcaptcha)
Full sample: [example.Mtcaptcha.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Mtcaptcha.go)

```go
import "encoding/json"

mtParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "sitekey":   "MTPublic-xxx",
    "pageurl":   "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":             "18",
    "mtcaptcha_params": string(mtParams),
})
```

---

<a id="sample-type-19-cutcaptcha"></a>
#### ✂️ Sample Type 19: Cutcaptcha
Official description: [Cutcaptcha](https://deathbycaptcha.com/api/cutcaptcha)
Full sample: [example.Cutcaptcha.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Cutcaptcha.go)

```go
import "encoding/json"

cutParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "apikey":    "SAs61IAI",
    "miserykey": "56a9e9b989aa8cf99e0cea28d4b4678b84fa7a4e",
    "pageurl":   "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":              "19",
    "cutcaptcha_params": string(cutParams),
})
```

---

<a id="sample-type-20-friendly-captcha"></a>
#### 💚 Sample Type 20: Friendly Captcha
Official description: [Friendly Captcha](https://deathbycaptcha.com/api/friendly)
Full sample: [example.Friendly.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Friendly.go)

```go
import "encoding/json"

friendlyParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "sitekey":   "FCMG...",
    "pageurl":   "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":            "20",
    "friendly_params": string(friendlyParams),
})
```

---

<a id="sample-type-21-datadome"></a>
#### 🛡️ Sample Type 21: DataDome
Official description: [DataDome](https://deathbycaptcha.com/api/datadome)
Full sample: [example.Datadome.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Datadome.go)

```go
import "encoding/json"

datadomeParams, _ := json.Marshal(map[string]string{
    "proxy":       "http://user:pass@127.0.0.1:1234",
    "proxytype":   "HTTP",
    "pageurl":     "https://target",
    "captcha_url": "https://target/captcha",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":            "21",
    "datadome_params": string(datadomeParams),
})
```

---

<a id="sample-type-23-tencent"></a>
#### 🇨🇳 Sample Type 23: Tencent
Official description: [Tencent](https://deathbycaptcha.com/api/tencent)
Full sample: [example.Tencent.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Tencent.go)

```go
import "encoding/json"

tencentParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "appid":     "2044554539",
    "pageurl":   "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":           "23",
    "tencent_params": string(tencentParams),
})
```

---

<a id="sample-type-24-atb"></a>
#### 🏷️ Sample Type 24: ATB
Official description: [ATB](https://deathbycaptcha.com/api/atb)
Full sample: [example.Atb.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Atb.go)

```go
import "encoding/json"

atbParams, _ := json.Marshal(map[string]string{
    "proxy":     "http://user:pass@127.0.0.1:1234",
    "proxytype": "HTTP",
    "appid":     "af23e041b22d000a11e22a230fa8991c",
    "apiserver": "https://cap.aisecurius.com",
    "pageurl":   "https://target",
})
captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, map[string]string{
    "type":       "24",
    "atb_params": string(atbParams),
})
```

<a id="captcha-types-extended-reference"></a>
## 📚 CAPTCHA Types Extended Reference

Full API-level documentation for selected CAPTCHA types: parameter references, payload schemas, request/response formats, token lifespans, and integration notes.

<a id="recaptcha-image-based-api"></a>
### ⛔ reCAPTCHA Image-Based API — Deprecated (Types 2 & 3)

> ⚠️ **Deprecated.** Types 2 (Coordinates) and 3 (Image Group) are legacy image-based reCAPTCHA challenge methods that are no longer used at captcha solving. Do not use them for new integrations — use the [reCAPTCHA Token API (v2 & v3)](#recaptcha-token-api) instead.

---

<a id="recaptcha-token-api"></a>
### 🔐 reCAPTCHA Token API (v2 & v3)

The Token-based API solves reCAPTCHA challenges by returning a token you inject directly into the page form, rather than clicking images. Given a site URL and site key, DBC solves the challenge on its side and returns a token valid for one submission.

- **Token Image API**: Provided a site URL and site key, the API returns a token that you use to submit the form on the page with the reCAPTCHA challenge.

---

<a id="recaptcha-v2-api-faq"></a>
### ❓ reCAPTCHA v2 API FAQ

**What's the Token Image API URL?**
To use the Token Image API you will have to send a HTTP POST Request to <http://api.dbcapi.me/api/captcha>
**What are the POST parameters for the Token image API?**
-   **`username`**: Your DBC account username
-   **`password`**: Your DBC account password
-   **`type`=4**: Type 4 specifies this is the reCAPTCHA v2 Token API
-   **`token_params`=json(payload)**: the data to access the recaptcha challenge
json payload structure:
    -   **`proxy`**: your proxy url and credentials (if any). Examples:
        -   <http://127.0.0.1:3128>
        -   <http://user:password@127.0.0.1:3128>

    -   **`proxytype`**: your proxy connection protocol. For supported proxy types refer to Which proxy types are supported?. Example:
        -   HTTP

    -   **`googlekey`**: the google recaptcha site key of the website with the recaptcha. For more details about the site key refer to What is a recaptcha site key?. Example:
        -   6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-

    -   **`pageurl`**: the url of the page with the recaptcha challenges. This url has to include the path in which the recaptcha is loaded. Example: if the recaptcha you want to solve is in <http://test.com/path1>, pageurl has to be <http://test.com/path1> and not <http://test.com>.
    -   **`data-s`**: This parameter is only required for solve the google search tokens, the ones visible, while google search trigger the robot protection. Use the data-s value inside the google search response html. For regulars tokens don't use this parameter.

The **`proxy`** parameter is optional, but we strongly recommend to use one to prevent token rejection by the provided page due to inconsistencies between the IP that solved the captcha (ours if no proxy is provided) and the IP that submitted the token for verification (yours).
**Note**: If **`proxy`** is provided, **`proxytype`** is a required parameter.
Full example of **`token_params`**:
```json
{
  "proxy": "http://127.0.0.1:3128",
  "proxytype": "HTTP",
  "googlekey": "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
  "pageurl": "http://test.com/path_with_recaptcha"
}
```
Example of **`token_params`** for google search captchas:
```json
{
  "googlekey": "6Le-wvkSA...",
  "pageurl": "...",
  "data-s": "IUdfh4rh0sd..."
}
```
**What's the response from the Token image API?**
The token image API response has the same structure as regular captchas' response. Refer to Polling for uploaded CAPTCHA status for details about the response. The token will come in the text key of the response. It's valid for one use and has a 2 minute lifespan. It will be a string like the following:
```bash
"03AOPBWq_RPO2vLzyk0h8gH0cA2X4v3tpYCPZR6Y4yxKy1s3Eo7CHZRQntxrdsaD2H0e6S3547xi1FlqJB4rob46J0-wfZMj6YpyVa0WGCfpWzBWcLn7tO_EYsvEC_3kfLNINWa5LnKrnJTDXTOz-JuCKvEXx0EQqzb0OU4z2np4uyu79lc_NdvL0IRFc3Cslu6UFV04CIfqXJBWCE5MY0Ag918r14b43ZdpwHSaVVrUqzCQMCybcGq0yxLQf9eSexFiAWmcWLI5nVNA81meTXhQlyCn5bbbI2IMSEErDqceZjf1mX3M67BhIb4"
```

---

<a id="what-is-recaptcha-v3"></a>
### 🔎 What is reCAPTCHA v3?
This API extends the reCAPTCHA v2 Token API with two additional parameters: `action` and **minimal score (`min_score`)**.
reCAPTCHA v3 returns a score from each user, that evaluate if user is a bot or human. Then the website uses the score value that could range from 0 to 1 to decide if will accept or not the requests. Lower scores near to 0 are identified as bot.
The `action` parameter at reCAPTCHA v3 is an additional data used to separate different captcha validations like for example **login**, **register**, **sales**, **etc**.

---

<a id="recaptcha-v3-api-faq"></a>
### ❓ reCAPTCHA v3 API FAQ

**What is `action` in reCAPTCHA v3?**
Is a new parameter that allows processing user actions on the website differently.
To find this we need to inspect the javascript code of the website looking for call of grecaptcha.execute function. Example:
```javascript
grecaptcha.execute('6Lc2fhwTAAAAAGatXTzFYfvlQMI2T7B6ji8UVV_f', {action: something})
```
Sometimes it's really hard to find it and we need to look through all javascript files. We may also try to find the value of action parameter inside ___grecaptcha_cfg configuration object. Also we can call grecaptcha.execute and inspect javascript code. The API will use "verify" default value it if we won't provide action in our request.
**What is `min-score` in reCAPTCHA v3 API?**
The minimal score needed for the captcha resolution. We recommend using the 0.3 min-score value, scores highers than 0.3 are hard to get.
**What are the POST parameters for the reCAPTCHA v3 API?**
-   **`username`**: Your DBC account username
-   **`password`**: Your DBC account password
-   **`type`=5**: Type 5 specifies this is reCAPTCHA v3 API
-   **`token_params`**=json(payload): the data to access the recaptcha challenge
json payload structure:
    -   **`proxy`**: your proxy url and credentials (if any).Examples:
        -   <http://127.0.0.1:3128>
        -   <http://user:password@127.0.0.1:3128>

    -   **`proxytype`**: your proxy connection protocol. For supported proxy types refer to Which proxy types are supported?. Example:
        -   HTTP

    -   **`googlekey`**: the google recaptcha site key of the website with the recaptcha. For more details about the site key refer to What is a recaptcha site key?. Example:
        -   6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-

    -   **`pageurl`**: the url of the page with the recaptcha challenges. This url has to include the path in which the recaptcha is loaded. Example: if the recaptcha you want to solve is in <http://test.com/path1>, pageurl has to be <http://test.com/path1> and not <http://test.com>.

    -   **`action`**: The action name.

    -   **`min_score`**: The minimal score, usually 0.3

The **`proxy`** parameter is optional, but we strongly recommend to use one to prevent rejection by the provided page due to inconsistencies between the IP that solved the captcha (ours if no proxy is provided) and the IP that submitted the solution for verification (yours).
**Note**: If **`proxy`** is provided, **`proxytype`** is a required parameter.
Full example of **`token_params`**:
```json
{
  "proxy": "http://127.0.0.1:3128",
  "proxytype": "HTTP",
  "googlekey": "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
  "pageurl": "http://test.com/path_with_recaptcha",
  "action": "example/action",
  "min_score": 0.3
}
```
**What's the response from reCAPTCHA v3 API?**
The response has the same structure as regular captcha. Refer to [Polling for uploaded CAPTCHA status](https://deathbycaptcha.com/api#polling-captcha) for details about the response. The solution will come in the **text** key of the response. It's valid for one use and has a 1 minute lifespan.

---

<a id="amazon-waf-api-faq"></a>
### 🛡️ Amazon WAF API (Type 16)

Amazon WAF Captcha (also referred to as AWS WAF Captcha) is part of the Intelligent Threat Mitigation system within Amazon AWS. It presents image-alignment challenges that DBC solves by returning a token you set as the `aws-waf-token` cookie on the target page.

- **Official documentation:** [deathbycaptcha.com/api/amazonwaf](https://deathbycaptcha.com/api/amazonwaf)
- **Go sample:** [examples/example.Amazon_Waf.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Amazon_Waf.go)

**API URL:** Send a HTTP POST request to `http://api.dbcapi.me/api/captcha`

**POST parameters:**

-   **`username`**: Your DBC account username
-   **`password`**: Your DBC account password
-   **`type`=16**: Type 16 specifies this is the Amazon WAF API
-   **`waf_params`=json(payload)**: The data needed to access the Amazon WAF challenge

`waf_params` payload fields:

| Parameter | Required | Description |
|---|---|---|
| `proxy` | Optional\* | Proxy URL with credentials. E.g. `http://user:password@127.0.0.1:3128` |
| `proxytype` | Required if proxy set | Proxy protocol. Currently only `HTTP` is supported. |
| `sitekey` | Required | Amazon WAF site key found in the page's captcha script (value of the `key` parameter) |
| `pageurl` | Required | Full URL of the page showing the Amazon WAF challenge (must include the path) |
| `iv` | Required | Value of the `iv` parameter found in the captcha script on the page |
| `context` | Required | Value of the `context` parameter found in the captcha script on the page |
| `challengejs` | Optional | URL of the `challenge.js` script referenced on the page |
| `captchajs` | Optional | URL of the `captcha.js` script referenced on the page |

> The `proxy` parameter is optional but strongly recommended — using a proxy prevents token rejection caused by IP inconsistencies between the solving machine (DBC) and the submitting machine (yours).
> **📌 Note:** If `proxy` is provided, `proxytype` is required.

Full example of `waf_params`:

```json
{
  "proxy": "http://user:password@127.0.0.1:1234",
  "proxytype": "HTTP",
  "sitekey": "AQIDAHjcYu/GjX+QlghicBgQ/7bFaQZ+m5FKCMDnO+vTbNg96AHDh0IR5vgzHNceHYqZR+GO...",
  "pageurl": "https://efw47fpad9.execute-api.us-east-1.amazonaws.com/latest",
  "iv": "CgAFRjIw2vAAABSM",
  "context": "zPT0jOl1rQlUNaldX6LUpn4D6Tl9bJ8VUQ/NrWFxPii..."
}
```

With optional `challengejs` and `captchajs`:

```json
{
  "proxy": "http://user:password@127.0.0.1:1234",
  "proxytype": "HTTP",
  "sitekey": "AQIDAHjcYu/GjX+QlghicBgQ/7bFaQZ+m5FKCMDnO+vTbNg96AHDh0IR5vgzHNceHYqZR+GO...",
  "pageurl": "https://efw47fpad9.execute-api.us-east-1.amazonaws.com/latest",
  "iv": "CgAFRjIw2vAAABSM",
  "context": "zPT0jOl1rQlUNaldX6LUpn4D6Tl9bJ8VUQ/NrWFxPii...",
  "challengejs": "https://41bcdd4fb3cb.610cd090.us-east-1.token.awswaf.com/41bcdd4fb3cb/0d21de737ccb/cd77baa6c832/challenge.js",
  "captchajs": "https://41bcdd4fb3cb.610cd090.us-east-1.captcha.awswaf.com/41bcdd4fb3cb/0d21de737ccb/cd77baa6c832/captcha.js"
}
```

**Response:** The API returns a token string valid for one use with a 1-minute lifespan. Once received, set it as the `aws-waf-token` cookie on the target page before submitting the form:

```
c3b50e60-d76c-4d13-ae25-159ec7ec3121:EQoAj4x6fnENAAAA:YIvITdQewAaLmaLXo4r6Es783keXM2ahoP...
```

---

<a id="cloudflare-turnstile-api-faq"></a>
### 🌐 Cloudflare Turnstile API (Type 12)

Cloudflare Turnstile is a CAPTCHA alternative that protects pages without requiring user interaction in most cases. DBC solves it by returning a token you inject into the target form or pass to the page's callback.

- **Official documentation:** [deathbycaptcha.com/api/turnstile](https://deathbycaptcha.com/api/turnstile)
- **Go sample:** [examples/example.Turnstile.go](https://github.com/deathbycaptcha/deathbycaptcha-api-client-go/blob/master/examples/example.Turnstile.go)

**API URL:** Send a HTTP POST request to `http://api.dbcapi.me/api/captcha`

| POST Parameter | Description |
|---|---|
| `username` | Your DBC account username |
| `password` | Your DBC account password |
| `type` | `12` — specifies this is a Turnstile API request |
| `turnstile_params` | JSON-encoded payload (see fields below) |

**`turnstile_params` payload fields:**

| Field | Required | Description |
|---|---|---|
| `proxy` | Optional | Proxy URL with optional credentials. E.g. `http://user:password@127.0.0.1:3128` |
| `proxytype` | Required if `proxy` set | Proxy connection protocol. Currently only `HTTP` is supported. |
| `sitekey` | Required | The Turnstile site key found in `data-sitekey` attribute, the captcha iframe URL, or the `turnstile.render` call. E.g. `0x4AAAAAAAGlwMzq_9z6S9Mh` |
| `pageurl` | Required | Full URL of the page hosting the Turnstile challenge, including path. E.g. `https://testsite.com/xxx-test` |
| `action` | Optional | Value of the `data-action` attribute or the `action` option passed to `turnstile.render`. |

> **📌 Note:** The `proxy` parameter is optional but strongly recommended to avoid rejection due to IP inconsistency between the solver and the submitter. If `proxy` is provided, `proxytype` becomes required.

**Example `turnstile_params` (basic):**

```json
{
    "proxy": "http://user:password@127.0.0.1:1234",
    "proxytype": "HTTP",
    "sitekey": "0x4AAAAAAAGlwMzq_9z6S9Mh",
    "pageurl": "https://testsite.com/xxx-test"
}
```

**Example `turnstile_params` with optional `action`:**

```json
{
    "proxy": "http://user:password@127.0.0.1:1234",
    "proxytype": "HTTP",
    "sitekey": "0x4AAAAAAAGlwMzq_9z6S9Mh",
    "pageurl": "https://testsite.com/xxx-test",
    "action": "login"
}
```

**Response:** The API returns a token string valid for one use with a 2-minute lifespan. Submit it via the `input[name="cf-turnstile-response"]` field (or `input[name="g-recaptcha-response"]` when reCAPTCHA compatibility mode is enabled), or pass it to the callback defined in `turnstile.render` / `data-callback`:

```
0.Ja5ditqqMhsYE6r9okpFeZVg6ixDXZcbSTzxypXJkAeN-D-VxmNaEBR_XfsPGk-nxhJFUwMERSIXk6npwAifIYfKuP5AZHeLgCAm0W6CAyJlc9WvO_t7pGYnR_wwbyUyroooPkOI9mfHeeXb1urRmsTF_kP5pU5kQ05OVx3EyuXK3nl0fd4y1u7SyYThi...
```

## 🧪 Testing & CI

### Running Unit Tests Locally

```bash
go test ./deathbycaptcha/... -v -timeout 120s
```

### Running with Coverage

```bash
go test ./deathbycaptcha/... -coverprofile=coverage.out -covermode=atomic -timeout 120s
go tool cover -func=coverage.out
```

### Running Integration Tests (Requires Credentials)

```bash
DBC_USERNAME=your_user DBC_PASSWORD=your_pass \
  go test -tags=integration ./tests/integration/... -v -timeout 300s
```

Without credentials, integration tests are automatically skipped.

### CI Workflows

| Workflow | File | Purpose |
|----------|------|---------|
| Unit Tests Go 1.25 | `.github/workflows/unit-tests-go125.yml` | Runs `go test ./deathbycaptcha/...` on Go 1.25 |
| Unit Tests Go 1.26 | `.github/workflows/unit-tests-go126.yml` | Runs `go test ./deathbycaptcha/...` on Go 1.26 |
| Coverage | `.github/workflows/coverage.yml` | Measures line coverage; deploys badge to GitHub Pages |
| API Integration | `.github/workflows/integration.yml` | Live API tests; skips when secrets are absent |
| Publish | `.github/workflows/publish.yml` | Triggers pkg.go.dev index on tagged releases |

**Coverage badge** is served from GitHub Pages at:
```
https://img.shields.io/endpoint?url=https://deathbycaptcha.github.io/deathbycaptcha-api-client-go/coverage-badge.json
```

Color thresholds: `brightgreen` ≥ 85 %, `yellow` 70–84 %, `red` < 70 %.

**Badge interpretation:**
- `Unit Tests Go 1.25` / `Unit Tests Go 1.26` — main unit test pipeline on Go 1.25 and Go 1.26 passed. Covers build, vet, and all unit tests without network access.
- `API Integration` — live API test pipeline passed (or skipped when secrets are absent — the badge shows `passing` in both cases).
- `Coverage` — latest reported line coverage, updated on each push to `main`. Accepted threshold: ≥ 80 %.

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

## ⚠️ Responsible Use

See [RESPONSIBLE_USE.md](RESPONSIBLE_USE.md). This library is intended for legitimate automation, testing, and research only. Do not use it to bypass protections on third-party sites without explicit authorization.
