// example.Normal_Captcha — DeathByCaptcha Go client example.
// Submits a standard image CAPTCHA (type 0) and prints the solution.
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

	// Use SocketClient for lower latency, or HttpClient for encrypted HTTPS.
	client := dbc.NewSocketClient(username, password)
	defer client.Close()

	// Path to the image CAPTCHA file.
	captchaFile := "images/normal.jpg"

	captcha, err := client.Decode(captchaFile, dbc.DefaultTimeout, map[string]string{})
	if err != nil {
		var ade *dbc.AccessDeniedException
		if ok := fmt.Sprintf("%v", err); ok != "" {
			_ = ade
		}
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
