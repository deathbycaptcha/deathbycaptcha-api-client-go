package deathbycaptcha_test

// coverage_boost_test.go — additional tests to push coverage to >= 80%.
// Covers: SocketClient.Decode, SocketClient.GetText, SocketClient.Close,
// HttpClient.GetStatus error path, effectiveTimeout branches, authtoken path,
// mapToCaptcha is_correct variants.

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/deathbycaptcha"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// SocketClient — GetText via mock
// ---------------------------------------------------------------------------

func TestSocketClient_GetText_Pending(t *testing.T) {
	loginResp := map[string]interface{}{"user": float64(1), "status": float64(0)}
	capResp := map[string]interface{}{"captcha": float64(77), "text": nil, "is_correct": false, "status": float64(0)}
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, capResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	text, err := c.GetText(77)
	require.NoError(t, err)
	assert.Empty(t, text)
}

func TestSocketClient_GetText_Solved(t *testing.T) {
	loginResp := map[string]interface{}{"user": float64(1), "status": float64(0)}
	solved := "done"
	capResp := map[string]interface{}{"captcha": float64(78), "text": solved, "is_correct": true, "status": float64(0)}
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, capResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	text, err := c.GetText(78)
	require.NoError(t, err)
	assert.Equal(t, "done", text)
}

// ---------------------------------------------------------------------------
// SocketClient — Decode solved
// ---------------------------------------------------------------------------

func TestSocketClient_Decode_Solved(t *testing.T) {
	loginResp := map[string]interface{}{"user": float64(1), "status": float64(0)}
	uploadResp := map[string]interface{}{"captcha": float64(200), "text": nil, "is_correct": false, "status": float64(0)}
	solved := "captchatext"
	pollResp := map[string]interface{}{"captcha": float64(200), "text": solved, "is_correct": true, "status": float64(0)}
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, uploadResp, pollResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	pngBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE,
	}
	cap, err := c.Decode(pngBytes, 30, map[string]string{})
	require.NoError(t, err)
	require.NotNil(t, cap)
	assert.Equal(t, "captchatext", *cap.Text)
}

// ---------------------------------------------------------------------------
// SocketClient — Decode timeout (never solved)
// ---------------------------------------------------------------------------

func TestSocketClient_Decode_Timeout(t *testing.T) {
	// Build a mock that always returns pending for every poll
	loginResp := map[string]interface{}{"user": float64(1), "status": float64(0)}
	uploadResp := map[string]interface{}{"captcha": float64(201), "text": nil, "is_correct": false, "status": float64(0)}
	pending := map[string]interface{}{"captcha": float64(201), "text": nil, "is_correct": false, "status": float64(0)}

	// Provide enough pending responses for 2s timeout
	responses := []map[string]interface{}{loginResp, uploadResp}
	for i := 0; i < 10; i++ {
		responses = append(responses, pending)
	}
	mock := newMockSocketServer(t, responses)
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	pngBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE,
	}
	cap, err := c.Decode(pngBytes, 2, map[string]string{})
	require.NoError(t, err)
	assert.Nil(t, cap)
}

// ---------------------------------------------------------------------------
// SocketClient — Close with active connection
// ---------------------------------------------------------------------------

func TestSocketClient_Close_WithConnection(t *testing.T) {
	loginResp := map[string]interface{}{"user": float64(1), "status": float64(0)}
	userResp := map[string]interface{}{"user": float64(1), "rate": 0.1, "balance": 5.0, "is_banned": false, "status": float64(0)}
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, userResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	_, err := c.GetUser()
	require.NoError(t, err)

	c.Close() // should not panic with active conn
	c.Close() // idempotent
}

// ---------------------------------------------------------------------------
// HttpClient — authtoken path (GetBalance uses POST /)
// ---------------------------------------------------------------------------

func TestHttpClient_Authtoken_GetBalance(t *testing.T) {
	var gotBody string
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"user":1,"rate":0.1,"balance":77.0,"is_banned":false}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := dbc.NewHttpClientWithToken("my-auth-token")
	c.SetBaseURL(srv.URL)

	bal, err := c.GetBalance()
	require.NoError(t, err)
	assert.Equal(t, 77.0, bal)
	assert.Contains(t, gotBody, "authtoken=my-auth-token")
}

// ---------------------------------------------------------------------------
// HttpClient — 400 error path
// ---------------------------------------------------------------------------

func TestHttpClient_BadRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	_, err := c.Upload(nil, map[string]string{"type": "4", "token_params": "{}"})
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// effectiveTimeout — isToken=true branch
// ---------------------------------------------------------------------------

func TestHttpClient_Decode_DefaultTokenTimeout(t *testing.T) {
	// Pass timeout=0 with nil captcha => should use DefaultTokenTimeout (120)
	// Use a mock that solves immediately so the test is fast.
	pollCount := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/captcha", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"captcha":50,"text":null,"is_correct":false}`)
		}
	})
	mux.HandleFunc("/captcha/50", func(w http.ResponseWriter, r *http.Request) {
		pollCount++
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"captcha":50,"text":"fast","is_correct":true}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	c := newHttpClientWithServer(srv)

	// timeout=0, captcha=nil => token path => effectiveTimeout = 120
	// But mock solves immediately so the test is instant.
	cap, err := c.Decode(nil, 0, map[string]string{"type": "4", "token_params": `{}`})
	require.NoError(t, err)
	require.NotNil(t, cap)
	assert.Equal(t, "fast", *cap.Text)
}

// ---------------------------------------------------------------------------
// mapToCaptcha — is_correct as numeric string "1"
// ---------------------------------------------------------------------------

func TestSocketClient_IsCorrect_StringTrue(t *testing.T) {
	loginResp := map[string]interface{}{"user": float64(1), "status": float64(0)}
	capResp := map[string]interface{}{"captcha": float64(9), "text": "ok", "is_correct": "1", "status": float64(0)}
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, capResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	cap, err := c.GetCaptcha(9)
	require.NoError(t, err)
	assert.True(t, cap.IsCorrect)
}

// ---------------------------------------------------------------------------
// SocketClient — Upload token type 4
// ---------------------------------------------------------------------------

func TestSocketClient_Upload_TokenType4(t *testing.T) {
	loginResp := map[string]interface{}{"user": float64(1), "status": float64(0)}
	uploadResp := map[string]interface{}{"captcha": float64(600), "text": nil, "is_correct": false, "status": float64(0)}
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, uploadResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("u", "p", mock.Addr())
	tokenJSON := `{"googlekey":"key","pageurl":"https://example.com"}`
	cap, err := c.Upload(nil, map[string]string{"type": "4", "token_params": tokenJSON})
	require.NoError(t, err)
	assert.Equal(t, 600, cap.CaptchaID)
}

// ---------------------------------------------------------------------------
// SocketClient — authtoken path
// ---------------------------------------------------------------------------

func TestSocketClient_Authtoken_GetBalance(t *testing.T) {
	loginResp := map[string]interface{}{"user": float64(5), "balance": 22.5, "status": float64(0)}
	userResp := map[string]interface{}{"user": float64(5), "rate": 0.1, "balance": 22.5, "is_banned": false, "status": float64(0)}
	mock := newMockSocketServer(t, []map[string]interface{}{loginResp, userResp})
	defer mock.Close()

	c := dbc.NewSocketClientForAddr("", "", mock.Addr())
	// Use the exported constructor to set authtoken
	_ = c
	c2 := dbc.NewSocketClientWithTokenForAddr("mytoken", mock.Addr())
	bal, err := c2.GetBalance()
	require.NoError(t, err)
	assert.Equal(t, 22.5, bal)
}
