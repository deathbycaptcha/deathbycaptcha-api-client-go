// Package deathbycaptcha provides HTTP and Socket API clients for the
// DeathByCaptcha CAPTCHA-solving service (https://deathbycaptcha.com).
//
// Two client implementations are provided:
//   - HttpClient  — REST over HTTPS (recommended for most use cases).
//   - SocketClient — persistent TCP socket (lower latency, higher throughput).
//
// Both implement the Client interface and are safe for concurrent use by
// multiple goroutines.
package deathbycaptcha

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const (
	// APIVersion is sent as the User-Agent header and as the socket "version"
	// field so the server can identify this client.
	APIVersion = "DBC/Go v4.7.0"

	// HTTPBaseURL is the canonical HTTPS endpoint.
	HTTPBaseURL = "https://api.dbcapi.me/api"

	// SocketHost is the TCP host for the socket API.
	SocketHost = "api.dbcapi.me"

	// DefaultTimeout is the default image-CAPTCHA solve timeout in seconds.
	DefaultTimeout = 60

	// DefaultTokenTimeout is the default token-CAPTCHA solve timeout in seconds.
	DefaultTokenTimeout = 120

	// DefaultPollInterval is used once the POLLS_INTERVAL schedule is exhausted.
	DefaultPollInterval = 3
)

// SocketPorts lists the available TCP ports for the socket API.
var SocketPorts = []int{8123, 8124, 8125, 8126, 8127, 8128, 8129, 8130}

// PollsInterval is the schedule of sleep durations (seconds) between each
// successive poll before falling back to DefaultPollInterval.
var PollsInterval = []int{1, 1, 2, 3, 2, 2, 3, 2, 2}

// ---------------------------------------------------------------------------
// Error types
// ---------------------------------------------------------------------------

// AccessDeniedException is returned when credentials are invalid or balance is insufficient.
type AccessDeniedException struct{ Msg string }

func (e *AccessDeniedException) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "access denied: check credentials or balance"
}

// ServiceOverloadException is returned when the server responds with HTTP 503.
type ServiceOverloadException struct{ Msg string }

func (e *ServiceOverloadException) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "service overloaded, try again later"
}

// ---------------------------------------------------------------------------
// Data models
// ---------------------------------------------------------------------------

// User holds account information returned by the API.
type User struct {
	UserID   int     `json:"user"`
	Rate     float64 `json:"rate"`
	Balance  float64 `json:"balance"`
	IsBanned bool    `json:"is_banned"`
}

// Captcha holds captcha details returned by the API.
type Captcha struct {
	CaptchaID int     `json:"captcha"`
	Text      *string `json:"text"`
	IsCorrect bool    `json:"is_correct"`
}

// IsSolved reports whether the captcha has been solved.
func (c *Captcha) IsSolved() bool {
	return c.Text != nil && *c.Text != ""
}

// serverStatus is used internally to unmarshal /status responses.
type serverStatus struct {
	IsServiceOverloaded bool `json:"is_service_overloaded"`
	Status              int  `json:"status"`
}

// ---------------------------------------------------------------------------
// Client interface
// ---------------------------------------------------------------------------

// Client is the common interface implemented by HttpClient and SocketClient.
type Client interface {
	// GetUser returns account information (ID, rate, balance, banned status).
	GetUser() (*User, error)

	// GetBalance returns the account balance in US cents.
	GetBalance() (float64, error)

	// GetCaptcha returns the details for an uploaded CAPTCHA.
	GetCaptcha(id int) (*Captcha, error)

	// GetText returns the solved text for a CAPTCHA, or empty string if pending.
	GetText(id int) (string, error)

	// Report reports an incorrectly solved CAPTCHA. Returns true on success.
	Report(id int) (bool, error)

	// Upload submits a CAPTCHA without waiting for the solution.
	// For image CAPTCHAs, pass a file path (string), []byte, or io.Reader.
	// For token CAPTCHAs, pass nil and set type/params via kwargs.
	Upload(captcha interface{}, params map[string]string) (*Captcha, error)

	// Decode submits a CAPTCHA and polls until solved or timeout expires.
	// Returns nil when the timeout is reached without a solution.
	Decode(captcha interface{}, timeout int, params map[string]string) (*Captcha, error)

	// Close releases any resources held by the client (e.g. socket connections).
	Close()
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// pollInterval returns the sleep duration for the given poll index and the
// next index to use.
func pollInterval(idx int) (time.Duration, int) {
	if idx < len(PollsInterval) {
		return time.Duration(PollsInterval[idx]) * time.Second, idx + 1
	}
	return time.Duration(DefaultPollInterval) * time.Second, idx + 1
}

// loadImage reads image bytes from a file path, []byte, or io.Reader.
func loadImage(captcha interface{}) ([]byte, error) {
	switch v := captcha.(type) {
	case string:
		data, err := os.ReadFile(v)
		if err != nil {
			return nil, fmt.Errorf("deathbycaptcha: cannot read file %q: %w", v, err)
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("deathbycaptcha: file %q is empty", v)
		}
		return data, nil
	case []byte:
		if len(v) == 0 {
			return nil, fmt.Errorf("deathbycaptcha: image bytes are empty")
		}
		return v, nil
	case io.Reader:
		data, err := io.ReadAll(v)
		if err != nil {
			return nil, fmt.Errorf("deathbycaptcha: cannot read image: %w", err)
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("deathbycaptcha: image reader produced no bytes")
		}
		return data, nil
	default:
		return nil, fmt.Errorf("deathbycaptcha: unsupported captcha type %T", captcha)
	}
}

// effectiveTimeout resolves the timeout to use: if t <= 0 and captcha is nil
// the token timeout is used; otherwise the image timeout is used.
func effectiveTimeout(t int, isToken bool) int {
	if t > 0 {
		return t
	}
	if isToken {
		return DefaultTokenTimeout
	}
	return DefaultTimeout
}

// ---------------------------------------------------------------------------
// HttpClient
// ---------------------------------------------------------------------------

// HttpClient calls the DeathByCaptcha REST API over HTTPS.
// It is safe for concurrent use.
type HttpClient struct {
	username  string
	password  string
	authtoken string
	baseURL   string
	http      *http.Client
}

// NewHttpClient creates a new HttpClient authenticated with username+password.
func NewHttpClient(username, password string) *HttpClient {
	return &HttpClient{
		username: username,
		password: password,
		baseURL:  HTTPBaseURL,
		http:     &http.Client{},
	}
}

// NewHttpClientWithToken creates a new HttpClient authenticated with an authtoken.
func NewHttpClientWithToken(authtoken string) *HttpClient {
	return &HttpClient{
		authtoken: authtoken,
		baseURL:   HTTPBaseURL,
		http:      &http.Client{},
	}
}

// SetBaseURL overrides the base URL (useful for testing with a mock server).
func (c *HttpClient) SetBaseURL(u string) { c.baseURL = u }

func (c *HttpClient) authFields() url.Values {
	v := url.Values{}
	if c.authtoken != "" {
		v.Set("authtoken", c.authtoken)
	} else {
		v.Set("username", c.username)
		v.Set("password", c.password)
	}
	return v
}

func (c *HttpClient) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", APIVersion)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.handleResponse(resp)
}

func (c *HttpClient) post(endpoint string, data url.Values) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, c.baseURL+endpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", APIVersion)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.handleResponse(resp)
}

func (c *HttpClient) handleResponse(resp *http.Response) ([]byte, error) {
	switch resp.StatusCode {
	case http.StatusForbidden:
		return nil, &AccessDeniedException{Msg: "access denied: check credentials or balance"}
	case http.StatusBadRequest, 413:
		return nil, fmt.Errorf("deathbycaptcha: captcha rejected (bad request)")
	case http.StatusServiceUnavailable:
		return nil, &ServiceOverloadException{Msg: "service overloaded, try again later"}
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("deathbycaptcha: unexpected HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// GetUser implements Client.
func (c *HttpClient) GetUser() (*User, error) {
	data := c.authFields()
	body, err := c.post("/", data)
	if err != nil {
		return nil, err
	}
	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("deathbycaptcha: cannot parse user: %w", err)
	}
	return &u, nil
}

// GetBalance implements Client.
func (c *HttpClient) GetBalance() (float64, error) {
	u, err := c.GetUser()
	if err != nil {
		return 0, err
	}
	return u.Balance, nil
}

// GetCaptcha implements Client.
func (c *HttpClient) GetCaptcha(id int) (*Captcha, error) {
	body, err := c.get(fmt.Sprintf("/captcha/%d", id))
	if err != nil {
		return nil, err
	}
	var cap Captcha
	if err := json.Unmarshal(body, &cap); err != nil {
		return nil, fmt.Errorf("deathbycaptcha: cannot parse captcha: %w", err)
	}
	return &cap, nil
}

// GetText implements Client.
func (c *HttpClient) GetText(id int) (string, error) {
	cap, err := c.GetCaptcha(id)
	if err != nil {
		return "", err
	}
	if cap.Text != nil {
		return *cap.Text, nil
	}
	return "", nil
}

// Report implements Client.
func (c *HttpClient) Report(id int) (bool, error) {
	data := c.authFields()
	body, err := c.post(fmt.Sprintf("/captcha/%d/report", id), data)
	if err != nil {
		return false, err
	}
	var cap Captcha
	if err := json.Unmarshal(body, &cap); err != nil {
		return false, nil
	}
	return !cap.IsCorrect, nil
}

// Upload implements Client.
// For image CAPTCHAs (type 0), pass the image as a file path, []byte, or io.Reader.
// For token/parameter CAPTCHAs, pass nil and include type + params fields in params.
func (c *HttpClient) Upload(captcha interface{}, params map[string]string) (*Captcha, error) {
	data := c.authFields()
	for k, v := range params {
		data.Set(k, v)
	}

	if captcha != nil {
		imgBytes, err := loadImage(captcha)
		if err != nil {
			return nil, err
		}
		encoded := "base64:" + base64.StdEncoding.EncodeToString(imgBytes)
		data.Set("captchafile", encoded)
	}

	body, err := c.post("/captcha", data)
	if err != nil {
		return nil, err
	}

	var cap Captcha
	if err := json.Unmarshal(body, &cap); err != nil {
		return nil, fmt.Errorf("deathbycaptcha: cannot parse upload response: %w", err)
	}
	if cap.CaptchaID == 0 {
		return nil, fmt.Errorf("deathbycaptcha: upload failed (captcha id = 0)")
	}
	return &cap, nil
}

// Decode implements Client.
func (c *HttpClient) Decode(captcha interface{}, timeout int, params map[string]string) (*Captcha, error) {
	return decode(c, captcha, timeout, params)
}

// Close implements Client (no-op for HttpClient).
func (c *HttpClient) Close() {}

// GetStatus queries the /status endpoint and returns a bool indicating whether
// the service is currently overloaded.
func (c *HttpClient) GetStatus() (bool, error) {
	body, err := c.get("/status")
	if err != nil {
		return false, err
	}
	var s serverStatus
	if err := json.Unmarshal(body, &s); err != nil {
		return false, fmt.Errorf("deathbycaptcha: cannot parse status: %w", err)
	}
	return s.IsServiceOverloaded, nil
}

// ---------------------------------------------------------------------------
// SocketClient
// ---------------------------------------------------------------------------

// SocketClient calls the DeathByCaptcha TCP socket API.
// It maintains a persistent connection and is safe for concurrent use.
type SocketClient struct {
	username   string
	password   string
	authtoken  string
	overrideAddr string // non-empty overrides SocketHost+randomPort (for testing)
	mu         sync.Mutex
	conn       net.Conn
	loggedIn   bool
}

// NewSocketClient creates a new SocketClient authenticated with username+password.
func NewSocketClient(username, password string) *SocketClient {
	return &SocketClient{username: username, password: password}
}

// NewSocketClientWithToken creates a new SocketClient authenticated with an authtoken.
func NewSocketClientWithToken(authtoken string) *SocketClient {
	return &SocketClient{authtoken: authtoken}
}

// NewSocketClientForAddr creates a SocketClient that connects to a specific
// address (host:port). Intended for testing with mock TCP servers.
func NewSocketClientForAddr(username, password, addr string) *SocketClient {
	return &SocketClient{username: username, password: password, overrideAddr: addr}
}

// NewSocketClientWithTokenForAddr creates a SocketClient with authtoken that
// connects to a specific address. Intended for testing.
func NewSocketClientWithTokenForAddr(authtoken, addr string) *SocketClient {
	return &SocketClient{authtoken: authtoken, overrideAddr: addr}
}

func (c *SocketClient) authFields() map[string]interface{} {
	if c.authtoken != "" {
		return map[string]interface{}{"authtoken": c.authtoken}
	}
	return map[string]interface{}{"username": c.username, "password": c.password}
}

func (c *SocketClient) connect() error {
	var addr string
	if c.overrideAddr != "" {
		addr = c.overrideAddr
	} else {
		port := SocketPorts[rand.Intn(len(SocketPorts))]
		addr = fmt.Sprintf("%s:%d", SocketHost, port)
	}
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("deathbycaptcha: cannot connect to socket %s: %w", addr, err)
	}
	c.conn = conn
	c.loggedIn = false
	return nil
}

// send serializes data as JSON, appends a null byte, and writes to the socket.
func (c *SocketClient) send(data map[string]interface{}) error {
	data["version"] = APIVersion
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	b = append(b, 0x00) // null-byte frame terminator
	_, err = c.conn.Write(b)
	return err
}

// recv reads from the socket until a null byte is received, then unmarshals JSON.
func (c *SocketClient) recv() (map[string]interface{}, error) {
	var buf []byte
	tmp := make([]byte, 256)
	for {
		n, err := c.conn.Read(tmp)
		if err != nil {
			return nil, err
		}
		buf = append(buf, tmp[:n]...)
		if idx := strings.IndexByte(string(buf), 0x00); idx >= 0 {
			buf = buf[:idx]
			break
		}
	}
	var result map[string]interface{}
	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, fmt.Errorf("deathbycaptcha: cannot parse socket response: %w", err)
	}
	return result, nil
}

func (c *SocketClient) login() error {
	data := c.authFields()
	data["cmd"] = "login"
	if err := c.send(data); err != nil {
		return err
	}
	resp, err := c.recv()
	if err != nil {
		return err
	}
	if errMsg, ok := resp["error"].(string); ok && errMsg != "" {
		return &AccessDeniedException{Msg: errMsg}
	}
	c.loggedIn = true
	return nil
}

// call sends a command over the socket, reconnecting and re-logging-in on failure.
func (c *SocketClient) call(cmd string, data map[string]interface{}) (map[string]interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for attempt := 0; attempt < 2; attempt++ {
		if c.conn == nil {
			if err := c.connect(); err != nil {
				return nil, err
			}
		}
		if !c.loggedIn && cmd != "login" {
			if err := c.login(); err != nil {
				c.conn.Close()
				c.conn = nil
				return nil, err
			}
		}

		d := make(map[string]interface{})
		for k, v := range data {
			d[k] = v
		}
		d["cmd"] = cmd

		if err := c.send(d); err != nil {
			c.conn.Close()
			c.conn = nil
			continue
		}
		resp, err := c.recv()
		if err != nil {
			c.conn.Close()
			c.conn = nil
			continue
		}
		if errMsg, ok := resp["error"].(string); ok && errMsg != "" {
			return nil, fmt.Errorf("deathbycaptcha: socket error: %s", errMsg)
		}
		return resp, nil
	}
	return nil, fmt.Errorf("deathbycaptcha: socket call failed after retries")
}

func mapToCaptcha(m map[string]interface{}) *Captcha {
	cap := &Captcha{}
	if id, ok := m["captcha"]; ok {
		switch v := id.(type) {
		case float64:
			cap.CaptchaID = int(v)
		case int:
			cap.CaptchaID = v
		}
	}
	if text, ok := m["text"]; ok && text != nil {
		s := fmt.Sprintf("%v", text)
		cap.Text = &s
	}
	if ic, ok := m["is_correct"]; ok {
		switch v := ic.(type) {
		case bool:
			cap.IsCorrect = v
		case float64:
			cap.IsCorrect = v != 0
		case string:
			cap.IsCorrect = v == "1" || strings.EqualFold(v, "true")
		}
	}
	return cap
}

func mapToUser(m map[string]interface{}) *User {
	u := &User{}
	if v, ok := m["user"]; ok {
		switch val := v.(type) {
		case float64:
			u.UserID = int(val)
		}
	}
	if v, ok := m["rate"]; ok {
		if val, ok := v.(float64); ok {
			u.Rate = val
		}
	}
	if v, ok := m["balance"]; ok {
		if val, ok := v.(float64); ok {
			u.Balance = val
		}
	}
	if v, ok := m["is_banned"]; ok {
		switch val := v.(type) {
		case bool:
			u.IsBanned = val
		case float64:
			u.IsBanned = val != 0
		}
	}
	return u
}

// GetUser implements Client.
func (c *SocketClient) GetUser() (*User, error) {
	resp, err := c.call("user", c.authFields())
	if err != nil {
		return nil, err
	}
	return mapToUser(resp), nil
}

// GetBalance implements Client.
func (c *SocketClient) GetBalance() (float64, error) {
	u, err := c.GetUser()
	if err != nil {
		return 0, err
	}
	return u.Balance, nil
}

// GetCaptcha implements Client.
func (c *SocketClient) GetCaptcha(id int) (*Captcha, error) {
	resp, err := c.call("captcha", map[string]interface{}{"captcha": id})
	if err != nil {
		return nil, err
	}
	return mapToCaptcha(resp), nil
}

// GetText implements Client.
func (c *SocketClient) GetText(id int) (string, error) {
	cap, err := c.GetCaptcha(id)
	if err != nil {
		return "", err
	}
	if cap.Text != nil {
		return *cap.Text, nil
	}
	return "", nil
}

// Report implements Client.
func (c *SocketClient) Report(id int) (bool, error) {
	resp, err := c.call("report", map[string]interface{}{"captcha": id})
	if err != nil {
		return false, err
	}
	cap := mapToCaptcha(resp)
	return !cap.IsCorrect, nil
}

// Upload implements Client.
func (c *SocketClient) Upload(captcha interface{}, params map[string]string) (*Captcha, error) {
	data := make(map[string]interface{})
	for k, v := range params {
		data[k] = v
	}

	if captcha != nil {
		imgBytes, err := loadImage(captcha)
		if err != nil {
			return nil, err
		}
		data["captcha_base64"] = base64.StdEncoding.EncodeToString(imgBytes)
	}

	// If type is a token captcha and params are a JSON string, parse them into map
	if tp, ok := params["token_params"]; ok {
		var parsed interface{}
		if err := json.Unmarshal([]byte(tp), &parsed); err == nil {
			data["token_params"] = parsed
		}
	}

	resp, err := c.call("upload", data)
	if err != nil {
		return nil, err
	}
	cap := mapToCaptcha(resp)
	if cap.CaptchaID == 0 {
		return nil, fmt.Errorf("deathbycaptcha: upload failed")
	}
	return cap, nil
}

// Decode implements Client.
func (c *SocketClient) Decode(captcha interface{}, timeout int, params map[string]string) (*Captcha, error) {
	return decode(c, captcha, timeout, params)
}

// Close implements Client. Closes the underlying TCP connection.
func (c *SocketClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

// ---------------------------------------------------------------------------
// Shared decode loop
// ---------------------------------------------------------------------------

func decode(cl Client, captcha interface{}, timeout int, params map[string]string) (*Captcha, error) {
	isToken := captcha == nil
	t := effectiveTimeout(timeout, isToken)

	uploaded, err := cl.Upload(captcha, params)
	if err != nil {
		return nil, err
	}

	deadline := time.Now().Add(time.Duration(t) * time.Second)
	intvlIdx := 0

	for time.Now().Before(deadline) {
		if uploaded.IsSolved() {
			return uploaded, nil
		}
		d, next := pollInterval(intvlIdx)
		intvlIdx = next
		time.Sleep(d)

		updated, err := cl.GetCaptcha(uploaded.CaptchaID)
		if err != nil {
			return nil, err
		}
		uploaded = updated
	}

	if uploaded.IsSolved() {
		return uploaded, nil
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Convenience helpers for token CAPTCHAs
// ---------------------------------------------------------------------------

// TokenParams builds a params map for token-based CAPTCHAs (type 4, 5, etc.)
// where jsonParams is already serialized as JSON.
func TokenParams(captchaType int, fieldName, jsonParams string) map[string]string {
	return map[string]string{
		"type":    strconv.Itoa(captchaType),
		fieldName: jsonParams,
	}
}
