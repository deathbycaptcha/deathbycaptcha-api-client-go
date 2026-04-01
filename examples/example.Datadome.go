// example.Datadome — DeathByCaptcha Go client example.
// Solves a DataDome CAPTCHA challenge (type 21) and prints the token.
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

	// Put your DataDome data here.
	captchaDict := map[string]string{
		"proxy":           "http://user:password@127.0.0.1:1234",
		"proxytype":       "HTTP",
		"captcha_url":     "https://geo.captcha-delivery.com/captcha/?initialCid=...",
		"userAgent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) ...",
		"pageurl":         "https://your-datadome-site.example.com/page",
		"datadome_cookie": "datadome=your_cookie_value_here",
	}
	jsonCaptcha, _ := json.Marshal(captchaDict)

	params := map[string]string{
		"type":             "21",
		"datadome_params":  string(jsonCaptcha),
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
