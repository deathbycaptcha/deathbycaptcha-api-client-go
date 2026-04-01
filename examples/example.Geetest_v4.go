// example.Geetest_v4 — DeathByCaptcha Go client example.
// Solves a GeeTest v4 challenge (type 9) and prints the token.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/v4/deathbycaptcha"
)

func main() {
	// Put your DBC account username and password here.
	username := "your_username"
	password := "your_password"

	client := dbc.NewSocketClient(username, password)
	defer client.Close()

	// Put your GeeTest v4 data here.
	captchaDict := map[string]string{
		"proxy":      "http://user:password@127.0.0.1:1234",
		"proxytype":  "HTTP",
		"captcha_id": "your_captcha_id_here",
		"pageurl":    "https://your-geetest4-site.example.com/page",
	}
	jsonCaptcha, _ := json.Marshal(captchaDict)

	params := map[string]string{
		"type":           "9",
		"geetest_params": string(jsonCaptcha),
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
