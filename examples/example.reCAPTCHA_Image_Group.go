// example.reCAPTCHA_Image_Group — DeathByCaptcha Go client example.
// Submits a reCAPTCHA Image Group challenge (type 3).
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

	// Path to the CAPTCHA screenshot (image group challenge).
	captchaFile := "images/recaptcha_image_group.jpg"

	params := map[string]string{
		"type":        "3",
		"banner":      "images/banner.jpg",
		"banner_text": "Select all images with buses",
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
