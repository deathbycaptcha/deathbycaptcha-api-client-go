// example.Audio — DeathByCaptcha Go client example.
// Solves an Audio CAPTCHA (type 13) by submitting a base64-encoded audio file.
package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/deathbycaptcha"
)

func main() {
	// Put your DBC account username and password here.
	username := "your_username"
	password := "your_password"

	client := dbc.NewSocketClient(username, password)
	defer client.Close()

	// Read the audio file and encode it as base64.
	audioData, err := os.ReadFile("images/audio.mp3")
	if err != nil {
		log.Fatalf("Failed to read audio file: %v", err)
	}
	audioBase64 := base64.StdEncoding.EncodeToString(audioData)

	params := map[string]string{
		"type":     "13",
		"audio":    audioBase64,
		"language": "en",
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
