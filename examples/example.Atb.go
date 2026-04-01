// example.Atb — DeathByCaptcha Go client example.
// Solves an ATB CAPTCHA challenge (type 24) and prints the token.
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

	// Put your ATB data here.
	captchaDict := map[string]string{
		"proxy":     "http://user:password@127.0.0.1:1234",
		"proxytype": "HTTP",
		"appid":     "your_atb_appid_here",
		"apiserver": "https://cap.aisecurius.com",
		"pageurl":   "https://your-atb-site.example.com/page",
	}
	jsonCaptcha, _ := json.Marshal(captchaDict)

	params := map[string]string{
		"type":       "24",
		"atb_params": string(jsonCaptcha),
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
