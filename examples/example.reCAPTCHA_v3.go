// example.reCAPTCHA_v3 — DeathByCaptcha Go client example.
// Solves a reCAPTCHA v3 challenge (type 5) and prints the token.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/deathbycaptcha"
)

func main() {
	// Put your DBC account username and password here.
	username := "your_username"
	password := "your_password"

	client := dbc.NewSocketClient(username, password)
	defer client.Close()

	// reCAPTCHA v3 requires 'action' — the name of the action that triggers
	// the reCAPTCHA v3 validation — and 'min_score' (0.1 – 0.9).
	captchaDict := map[string]interface{}{
		"proxy":     "http://user:password@127.0.0.1:1234",
		"proxytype": "HTTP",
		"googlekey": "6LdyC2cUAAAAACGuDKpXeDorzUDWXmdqeg-xy696",
		"pageurl":   "https://recaptchav3.demo.com/scores.php",
		"action":    "examples/v3scores",
		"min_score": 0.3,
	}
	jsonCaptcha, _ := json.Marshal(captchaDict)

	params := map[string]string{
		"type":         "5",
		"token_params": string(jsonCaptcha),
	}

	captcha, err := client.Decode(nil, dbc.DefaultTokenTimeout, params)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if captcha == nil {
		fmt.Println("Failed to solve CAPTCHA (timeout).")
		return
	}

	fmt.Printf("CAPTCHA %d solved: %s\n", captcha.CaptchaID, *captcha.Text)

	// Report if the solution is incorrect.
	// client.Report(captcha.CaptchaID)
}
