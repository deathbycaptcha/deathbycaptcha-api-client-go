// example.Goroutine_Socket — DeathByCaptcha Go client, concurrent Socket usage.
//
// Demonstrates the recommended pattern for concurrent use of SocketClient:
// each goroutine creates its own SocketClient instance backed by its own TCP
// connection. This avoids head-of-line blocking on a single socket and keeps
// goroutine ownership clear.
//
// Note: sharing a single SocketClient across goroutines is safe (the client
// uses sync.Mutex internally), but one-per-goroutine is preferred for
// throughput in parallel-solving scenarios.
//
// Usage:
//
//	DBC_USERNAME=user DBC_PASSWORD=pass go run example.Goroutine_Socket.go
//	DBC_USERNAME=user DBC_PASSWORD=pass DBC_THREADS=4 DBC_IMAGE_PATH=images/normal.jpg \
//	  go run example.Goroutine_Socket.go
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/v4/deathbycaptcha"
)

// socketWorker runs in its own goroutine and owns its own SocketClient instance.
func socketWorker(username, password, imagePath string, decodeMode bool, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Each goroutine creates its own client — independent TCP connection.
	client := dbc.NewSocketClient(username, password)
	defer client.Close()

	if !decodeMode {
		bal, err := client.GetBalance()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[SOCKET goroutine %d] error: %v\n", workerID, err)
			return
		}
		fmt.Printf("[SOCKET goroutine %d] balance: %.4f US cents\n", workerID, bal)
		return
	}

	captcha, err := client.Decode(imagePath, dbc.DefaultTimeout, map[string]string{})
	if err != nil {
		var ade *dbc.AccessDeniedException
		if errors.As(err, &ade) {
			fmt.Fprintf(os.Stderr, "[SOCKET goroutine %d] access denied: %v\n", workerID, err)
		} else {
			fmt.Fprintf(os.Stderr, "[SOCKET goroutine %d] error: %v\n", workerID, err)
		}
		return
	}
	if captcha == nil {
		fmt.Printf("[SOCKET goroutine %d] timeout — no solution received\n", workerID)
		return
	}

	fmt.Printf("[SOCKET goroutine %d] solved CAPTCHA %d: %s\n", workerID, captcha.CaptchaID, *captcha.Text)

	// Uncomment to report an incorrect solution:
	// client.Report(captcha.CaptchaID)
}

func main() {
	username := os.Getenv("DBC_USERNAME")
	password := os.Getenv("DBC_PASSWORD")
	imagePath := os.Getenv("DBC_IMAGE_PATH")
	threadsEnv := os.Getenv("DBC_THREADS")

	if username == "" || password == "" {
		log.Fatal("Set DBC_USERNAME and DBC_PASSWORD environment variables before running.")
	}

	goroutineCount := 2
	if threadsEnv != "" {
		n, err := strconv.Atoi(threadsEnv)
		if err != nil || n <= 0 {
			log.Fatal("DBC_THREADS must be a positive integer.")
		}
		goroutineCount = n
	}

	decodeMode := imagePath != ""

	// Sanity-check credentials on the main goroutine before spawning workers.
	sanity := dbc.NewSocketClient(username, password)
	bal, err := sanity.GetBalance()
	sanity.Close()
	if err != nil {
		log.Fatalf("[SOCKET main] credential check failed: %v", err)
	}
	fmt.Printf("[SOCKET main] initial balance: %.4f US cents\n", bal)

	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	for i := 1; i <= goroutineCount; i++ {
		go socketWorker(username, password, imagePath, decodeMode, i, &wg)
	}

	wg.Wait()
	fmt.Println("[SOCKET main] all goroutines finished.")
}
