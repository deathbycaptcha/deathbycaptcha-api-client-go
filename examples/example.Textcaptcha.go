// example.Textcaptcha — DeathByCaptcha Go client example.
// Solves a Text CAPTCHA challenge (type 11) and prints the answer.
package main

import (
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

	// Put your text CAPTCHA question here.
	params := map[string]string{
		"type":        "11",
		"textcaptcha": "What is 2 + 2?",
	}

	captcha, err := client.Decode(nil, dbc.DefaultTimeout, params)
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
