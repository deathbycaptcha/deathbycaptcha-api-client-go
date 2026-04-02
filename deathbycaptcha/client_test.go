package deathbycaptcha_test

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/v4/deathbycaptcha"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func userJSON() string {
	return `{"user":12345,"rate":0.139,"balance":250.0,"is_banned":false}`
}

func captchaPendingJSON(id int) string {
	return `{"captcha":` + strconv.Itoa(id) + `,"text":null,"is_correct":false}`
}

func captchaSolvedJSON(id int, text string) string {
	return `{"captcha":` + strconv.Itoa(id) + `,"text":"` + text + `","is_correct":true}`
}

// newMockServer starts an httptest.Server whose handler is built from a map
// of "METHOD /path" -> response body, all returning 200 with JSON content type.
func newMockServer(routes map[string]string) *httptest.Server {
	mux := http.NewServeMux()
	for pattern, body := range routes {
		parts := strings.SplitN(pattern, " ", 2)
		method := parts[0]
		path := parts[1]
		b := body
		m := method
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, b)
		})
	}
	return httptest.NewServer(mux)
}

func newHttpClientWithServer(server *httptest.Server) *dbc.HttpClient {
	c := dbc.NewHttpClient("testuser", "testpass")
	c.SetBaseURL(server.URL)
	return c
}

// ---------------------------------------------------------------------------
// Model tests
// ---------------------------------------------------------------------------

func TestCaptchaIsSolved_WithText(t *testing.T) {
	text := "ABCDEF"
	cap := &dbc.Captcha{CaptchaID: 1, Text: &text, IsCorrect: true}
	assert.True(t, cap.IsSolved())
}

func TestCaptchaIsSolved_NilText(t *testing.T) {
	cap := &dbc.Captcha{CaptchaID: 1, Text: nil, IsCorrect: false}
	assert.False(t, cap.IsSolved())
}

func TestCaptchaIsSolved_EmptyText(t *testing.T) {
	empty := ""
	cap := &dbc.Captcha{CaptchaID: 1, Text: &empty, IsCorrect: false}
	assert.False(t, cap.IsSolved())
}

func TestUserModel(t *testing.T) {
	u := &dbc.User{UserID: 99, Rate: 0.5, Balance: 10.0, IsBanned: true}
	assert.Equal(t, 99, u.UserID)
	assert.True(t, u.IsBanned)
}

// ---------------------------------------------------------------------------
// AccessDeniedException
// ---------------------------------------------------------------------------

func TestAccessDeniedException_Error(t *testing.T) {
	e := &dbc.AccessDeniedException{}
	assert.NotEmpty(t, e.Error())
}

func TestServiceOverloadException_Error(t *testing.T) {
	e := &dbc.ServiceOverloadException{}
	assert.NotEmpty(t, e.Error())
}

// ---------------------------------------------------------------------------
// HttpClient — construction
// ---------------------------------------------------------------------------

func TestNewHttpClient_UsernamePassword(t *testing.T) {
	c := dbc.NewHttpClient("user", "pass")
	require.NotNil(t, c)
}

func TestNewHttpClientWithToken(t *testing.T) {
	c := dbc.NewHttpClientWithToken("mytoken")
	require.NotNil(t, c)
}

// ---------------------------------------------------------------------------
// HttpClient — GetUser
// ---------------------------------------------------------------------------

func TestHttpClient_GetUser(t *testing.T) {
	srv := newMockServer(map[string]string{"POST /": userJSON()})
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	u, err := c.GetUser()
	require.NoError(t, err)
	assert.Equal(t, 12345, u.UserID)
	assert.Equal(t, 250.0, u.Balance)
	assert.False(t, u.IsBanned)
}

// ---------------------------------------------------------------------------
// HttpClient — GetBalance
// ---------------------------------------------------------------------------

func TestHttpClient_GetBalance(t *testing.T) {
	srv := newMockServer(map[string]string{"POST /": userJSON()})
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	bal, err := c.GetBalance()
	require.NoError(t, err)
	assert.Equal(t, 250.0, bal)
}

// ---------------------------------------------------------------------------
// HttpClient — GetCaptcha
// ---------------------------------------------------------------------------

func TestHttpClient_GetCaptcha_Pending(t *testing.T) {
	srv := newMockServer(map[string]string{"GET /captcha/555": captchaPendingJSON(555)})
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	cap, err := c.GetCaptcha(555)
	require.NoError(t, err)
	assert.Equal(t, 555, cap.CaptchaID)
	assert.False(t, cap.IsSolved())
}

func TestHttpClient_GetCaptcha_Solved(t *testing.T) {
	srv := newMockServer(map[string]string{"GET /captcha/777": captchaSolvedJSON(777, "HELLO")})
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	cap, err := c.GetCaptcha(777)
	require.NoError(t, err)
	assert.True(t, cap.IsSolved())
	assert.Equal(t, "HELLO", *cap.Text)
}

// ---------------------------------------------------------------------------
// HttpClient — GetText
// ---------------------------------------------------------------------------

func TestHttpClient_GetText_Nil(t *testing.T) {
	srv := newMockServer(map[string]string{"GET /captcha/1": captchaPendingJSON(1)})
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	text, err := c.GetText(1)
	require.NoError(t, err)
	assert.Empty(t, text)
}

func TestHttpClient_GetText_Solved(t *testing.T) {
	srv := newMockServer(map[string]string{"GET /captcha/2": captchaSolvedJSON(2, "SOLVED")})
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	text, err := c.GetText(2)
	require.NoError(t, err)
	assert.Equal(t, "SOLVED", text)
}

// ---------------------------------------------------------------------------
// HttpClient — Report
// ---------------------------------------------------------------------------

func TestHttpClient_Report_Success(t *testing.T) {
	body := `{"captcha":10,"text":"wrong","is_correct":false}`
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha/10/report", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	ok, err := c.Report(10)
	require.NoError(t, err)
	assert.True(t, ok)
}

// ---------------------------------------------------------------------------
// HttpClient — Upload image (base64)
// ---------------------------------------------------------------------------

func TestHttpClient_Upload_Image(t *testing.T) {
	uploadBody := `{"captcha":999,"text":null,"is_correct":false}`
	var gotBody string
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, uploadBody)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	// Use tiny PNG bytes (1x1 truecolor minimal PNG, no error for type check)
	imgBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG header
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x00, 0x02, 0x00, 0x01, 0xE2, 0x21, 0xBC,
		0x33, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, // IEND chunk
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}
	cap, err := c.Upload(imgBytes, map[string]string{})
	require.NoError(t, err)
	assert.Equal(t, 999, cap.CaptchaID)
	// The body is application/x-www-form-urlencoded; "base64:" becomes "base64%3A"
	assert.True(t,
		strings.Contains(gotBody, "base64:") || strings.Contains(gotBody, "base64%3A"),
		"expected base64 prefix in body")
}

// ---------------------------------------------------------------------------
// HttpClient — Upload token type 4
// ---------------------------------------------------------------------------

func TestHttpClient_Upload_TokenType4(t *testing.T) {
	uploadBody := `{"captcha":888,"text":null,"is_correct":false}`
	var gotBody string
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, uploadBody)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	tokenJSON := `{"googlekey":"sitekey","pageurl":"https://example.com"}`
	params := map[string]string{"type": "4", "token_params": tokenJSON}
	cap, err := c.Upload(nil, params)
	require.NoError(t, err)
	assert.Equal(t, 888, cap.CaptchaID)
	assert.Contains(t, gotBody, "token_params")
	assert.Contains(t, gotBody, "type=4")
}

// ---------------------------------------------------------------------------
// HttpClient — 403 raises AccessDeniedException
// ---------------------------------------------------------------------------

func TestHttpClient_AccessDenied(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	_, err := c.GetUser()
	require.Error(t, err)
	var ade *dbc.AccessDeniedException
	assert.ErrorAs(t, err, &ade)
}

// ---------------------------------------------------------------------------
// HttpClient — 503 raises ServiceOverloadException
// ---------------------------------------------------------------------------

func TestHttpClient_ServiceOverload(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	_, err := c.GetUser()
	require.Error(t, err)
	var soe *dbc.ServiceOverloadException
	assert.ErrorAs(t, err, &soe)
}

// ---------------------------------------------------------------------------
// HttpClient — Decode timeout (mock always returns pending)
// ---------------------------------------------------------------------------

func TestHttpClient_Decode_Timeout(t *testing.T) {
	callCount := 0
	mux := http.NewServeMux()
	// Upload endpoint
	mux.HandleFunc("/captcha", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"captcha":11,"text":null,"is_correct":false}`)
		}
	})
	// Poll endpoint
	mux.HandleFunc("/captcha/11", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"captcha":11,"text":null,"is_correct":false}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	// Use timeout=2 (very short, will exhaust)
	cap, err := c.Decode(nil, 2, map[string]string{"type": "4", "token_params": `{}`})
	require.NoError(t, err)
	assert.Nil(t, cap)
}

// ---------------------------------------------------------------------------
// HttpClient — Decode solved on second poll
// ---------------------------------------------------------------------------

func TestHttpClient_Decode_SolvedOnSecondPoll(t *testing.T) {
	pollCount := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"captcha":22,"text":null,"is_correct":false}`)
		}
	})
	mux.HandleFunc("/captcha/22", func(w http.ResponseWriter, r *http.Request) {
		pollCount++
		w.Header().Set("Content-Type", "application/json")
		if pollCount >= 2 {
			io.WriteString(w, captchaSolvedJSON(22, "RESULT"))
		} else {
			io.WriteString(w, captchaPendingJSON(22))
		}
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	cap, err := c.Decode(nil, 30, map[string]string{"type": "4", "token_params": `{}`})
	require.NoError(t, err)
	require.NotNil(t, cap)
	assert.Equal(t, "RESULT", *cap.Text)
}

// ---------------------------------------------------------------------------
// HttpClient — Close is a no-op
// ---------------------------------------------------------------------------

func TestHttpClient_Close(t *testing.T) {
	c := dbc.NewHttpClient("u", "p")
	c.Close() // must not panic
	c.Close() // safe to call twice
}

// ---------------------------------------------------------------------------
// HttpClient — GetStatus
// ---------------------------------------------------------------------------

func TestHttpClient_GetStatus(t *testing.T) {
	srv := newMockServer(map[string]string{
		"GET /status": `{"is_service_overloaded":false,"status":0}`,
	})
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	overloaded, err := c.GetStatus()
	require.NoError(t, err)
	assert.False(t, overloaded)
}

// ---------------------------------------------------------------------------
// HttpClient — Upload from file path
// ---------------------------------------------------------------------------

func TestHttpClient_Upload_FilePath(t *testing.T) {
	// Write a temp PNG file
	f, err := os.CreateTemp("", "test*.png")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	pngBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x00, 0x02, 0x00, 0x01, 0xE2, 0x21, 0xBC,
		0x33, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}
	f.Write(pngBytes)
	f.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"captcha":300,"text":null,"is_correct":false}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	cap, err := c.Upload(f.Name(), map[string]string{})
	require.NoError(t, err)
	assert.Equal(t, 300, cap.CaptchaID)
}

// ---------------------------------------------------------------------------
// HttpClient — Upload from io.Reader
// ---------------------------------------------------------------------------

func TestHttpClient_Upload_IOReader(t *testing.T) {
	pngBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x00, 0x02, 0x00, 0x01, 0xE2, 0x21, 0xBC,
		0x33, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"captcha":301,"text":null,"is_correct":false}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	cap, err := c.Upload(strings.NewReader(string(pngBytes)), map[string]string{})
	require.NoError(t, err)
	assert.Equal(t, 301, cap.CaptchaID)
}

// ---------------------------------------------------------------------------
// SocketClient — construction (no connection during build)
// ---------------------------------------------------------------------------

func TestNewSocketClient(t *testing.T) {
	c := dbc.NewSocketClient("user", "pass")
	require.NotNil(t, c)
}

func TestNewSocketClientWithToken(t *testing.T) {
	c := dbc.NewSocketClientWithToken("tok")
	require.NotNil(t, c)
}

func TestSocketClient_Close_Idempotent(t *testing.T) {
	c := dbc.NewSocketClient("u", "p")
	c.Close()
	c.Close() // must not panic
}

// ---------------------------------------------------------------------------
// SocketClient — mock TCP server helpers
// ---------------------------------------------------------------------------

type mockSocketServer struct {
	ln        net.Listener
	responses []map[string]interface{}
	idx       int
	mu        sync.Mutex
}

func newMockSocketServer(t *testing.T, responses []map[string]interface{}) *mockSocketServer {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	s := &mockSocketServer{ln: ln, responses: responses}
	go s.serve()
	return s
}

func (s *mockSocketServer) Addr() string { return s.ln.Addr().String() }

func (s *mockSocketServer) serve() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return
		}
		go s.handleConn(conn)
	}
}

func (s *mockSocketServer) handleConn(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// scanner.Scan() returns one line (strips the \n delimiter).
		// The client sends JSON + \r\n; Scanner strips \n, leaving possible \r.
		s.mu.Lock()
		if s.idx >= len(s.responses) {
			s.mu.Unlock()
			return
		}
		resp := s.responses[s.idx]
		s.idx++
		s.mu.Unlock()

		b, _ := json.Marshal(resp)
		b = append(b, '\r', '\n') // CRLF frame terminator (matches socketTerminator)
		conn.Write(b)             //nolint:errcheck
	}
}

func (s *mockSocketServer) Close() { s.ln.Close() }

// ---------------------------------------------------------------------------
// SocketClient — serialization: JSON + CRLF terminator
// ---------------------------------------------------------------------------

func TestSocketClient_FrameHasCRLF(t *testing.T) {
	// The mock server simply echoes the login response back
	loginResp := map[string]interface{}{
		"user":      float64(99),
		"rate":      0.139,
		"balance":   50.0,
		"is_banned": false,
		"status":    float64(0),
	}
	// first call is login, second is user command
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, loginResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("user", "pass", mock.Addr())
	u, err := c.GetUser()
	require.NoError(t, err)
	assert.Equal(t, 99, u.UserID)
}

// ---------------------------------------------------------------------------
// SocketClient — GetBalance via mock
// ---------------------------------------------------------------------------

func TestSocketClient_GetBalance(t *testing.T) {
	resp := map[string]interface{}{
		"user":      float64(1),
		"rate":      0.1,
		"balance":   99.9,
		"is_banned": false,
		"status":    float64(0),
	}
	mock := newMockSocketServer(t, []map[string]interface{}{resp, resp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	bal, err := c.GetBalance()
	require.NoError(t, err)
	assert.Equal(t, 99.9, bal)
}

// ---------------------------------------------------------------------------
// SocketClient — GetCaptcha pending via mock
// ---------------------------------------------------------------------------

func TestSocketClient_GetCaptcha_Pending(t *testing.T) {
	loginResp := map[string]interface{}{"user": float64(1), "status": float64(0)}
	capResp := map[string]interface{}{"captcha": float64(55), "text": nil, "is_correct": false, "status": float64(0)}
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, capResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	cap, err := c.GetCaptcha(55)
	require.NoError(t, err)
	assert.Equal(t, 55, cap.CaptchaID)
	assert.False(t, cap.IsSolved())
}

// ---------------------------------------------------------------------------
// SocketClient — Report via mock
// ---------------------------------------------------------------------------

func TestSocketClient_Report(t *testing.T) {
	loginResp := map[string]interface{}{"user": float64(1), "status": float64(0)}
	reportResp := map[string]interface{}{"captcha": float64(10), "text": "wrong", "is_correct": false, "status": float64(0)}
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, reportResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	ok, err := c.Report(10)
	require.NoError(t, err)
	assert.True(t, ok)
}

// ---------------------------------------------------------------------------
// SocketClient — Upload via mock
// ---------------------------------------------------------------------------

func TestSocketClient_Upload(t *testing.T) {
	loginResp := map[string]interface{}{"user": float64(1), "status": float64(0)}
	uploadResp := map[string]interface{}{"captcha": float64(500), "text": nil, "is_correct": false, "status": float64(0)}
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, uploadResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	pngBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE,
	}
	cap, err := c.Upload(pngBytes, map[string]string{})
	require.NoError(t, err)
	assert.Equal(t, 500, cap.CaptchaID)
}

// ---------------------------------------------------------------------------
// TokenParams helper
// ---------------------------------------------------------------------------

func TestTokenParams_Type4(t *testing.T) {
	p := dbc.TokenParams(4, "token_params", `{"googlekey":"k","pageurl":"u"}`)
	assert.Equal(t, "4", p["type"])
	assert.Contains(t, p["token_params"], "googlekey")
}

func TestTokenParams_Type12(t *testing.T) {
	p := dbc.TokenParams(12, "turnstile_params", `{"sitekey":"s","pageurl":"u"}`)
	assert.Equal(t, "12", p["type"])
}

// ---------------------------------------------------------------------------
// loadImage error paths
// ---------------------------------------------------------------------------

func TestHttpClient_Upload_EmptyBytes(t *testing.T) {
	c := dbc.NewHttpClient("u", "p")
	_, err := c.Upload([]byte{}, map[string]string{})
	require.Error(t, err)
}

func TestHttpClient_Upload_NonExistentFile(t *testing.T) {
	c := dbc.NewHttpClient("u", "p")
	_, err := c.Upload("/non/existent/path.jpg", map[string]string{})
	require.Error(t, err)
}

func TestHttpClient_Upload_UnsupportedType(t *testing.T) {
	c := dbc.NewHttpClient("u", "p")
	_, err := c.Upload(12345, map[string]string{})
	require.Error(t, err)
}
