//go:build integration

// Package integration contains live API tests that require valid DBC credentials.
// They are skipped automatically when DBC_USERNAME / DBC_PASSWORD env vars are absent.
//
// Run with:
//
//	DBC_USERNAME=user DBC_PASSWORD=pass go test -tags=integration ./tests/integration/... -v
package integration

import (
	"os"
	"testing"

	dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/deathbycaptcha"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func credentials(t *testing.T) (string, string) {
	t.Helper()
	username := os.Getenv("DBC_USERNAME")
	password := os.Getenv("DBC_PASSWORD")
	if username == "" || password == "" {
		t.Skip("DBC_USERNAME and/or DBC_PASSWORD not set; skipping integration test")
	}
	return username, password
}

// TestIntegration_Balance checks that get_balance() returns a positive value.
func TestIntegration_Balance(t *testing.T) {
	user, pass := credentials(t)
	c := dbc.NewHttpClient(user, pass)
	defer c.Close()

	bal, err := c.GetBalance()
	require.NoError(t, err)
	assert.Greater(t, bal, 0.0, "balance should be > 0")
}

// TestIntegration_UserInfo checks that GetUser returns a non-zero user_id.
func TestIntegration_UserInfo(t *testing.T) {
	user, pass := credentials(t)
	c := dbc.NewHttpClient(user, pass)
	defer c.Close()

	u, err := c.GetUser()
	require.NoError(t, err)
	assert.NotZero(t, u.UserID, "user_id should be non-zero")
	assert.False(t, u.IsBanned, "account should not be banned")
}

// TestIntegration_Status verifies the service is not overloaded.
func TestIntegration_Status(t *testing.T) {
	user, pass := credentials(t)
	c := dbc.NewHttpClient(user, pass)
	defer c.Close()

	overloaded, err := c.GetStatus()
	require.NoError(t, err)
	assert.False(t, overloaded, "service should not be overloaded")
}

// TestIntegration_DecodeImage uploads the fixture image and waits for a solution.
func TestIntegration_DecodeImage(t *testing.T) {
	user, pass := credentials(t)
	c := dbc.NewHttpClient(user, pass)
	defer c.Close()

	cap, err := c.Decode("../../tests/fixtures/test.jpg", 60, map[string]string{})
	require.NoError(t, err)
	require.NotNil(t, cap, "expected a solved captcha, got nil (timeout)")
	assert.NotEmpty(t, *cap.Text, "solved text should not be empty")
}
