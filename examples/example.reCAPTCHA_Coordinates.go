// example.reCAPTCHA_Coordinates — DeathByCaptcha Go client example.
// Submits a reCAPTCHA Coordinates challenge (type 2) using a screenshot.
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

	// Path to the screenshot of the reCAPTCHA grid.
	captchaFile := "images/recaptcha_coordinates.jpg"

	params := map[string]string{
		"type": "2",
	}

	captcha, err := client.Decode(captchaFile, dbc.DefaultTimeout, params)
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
