// example.Amazon_Waf — DeathByCaptcha Go client example.
// Solves an Amazon WAF challenge (type 16) and prints the token.
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

	// Put your Amazon WAF data here.
	captchaDict := map[string]string{
		"proxy":     "http://user:password@127.0.0.1:1234",
		"proxytype": "HTTP",
		"sitekey":   "your_sitekey_here",
		"pageurl":   "https://efw47fpad9.execute-api.us-east-1.amazonaws.com/latest",
		"iv":        "your_iv_here",
		"context":   "your_context_here",
	}
	jsonCaptcha, _ := json.Marshal(captchaDict)

	params := map[string]string{
		"type":       "16",
		"waf_params": string(jsonCaptcha),
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
