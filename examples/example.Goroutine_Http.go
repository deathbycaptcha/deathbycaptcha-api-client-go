// example.Goroutine_Http — DeathByCaptcha Go client, concurrent HTTP usage.
//
// Demonstrates the recommended pattern for concurrent use of HttpClient:
// each goroutine creates its own client instance, which avoids any shared
// connection state. Go's net/http.Client is also safe for shared use, but
// one-instance-per-goroutine gives the clearest ownership and is idiomatic
// in high-throughput producer/consumer pipelines.
//
// Usage:
//
//	DBC_USERNAME=user DBC_PASSWORD=pass go run example.Goroutine_Http.go
//	DBC_USERNAME=user DBC_PASSWORD=pass DBC_THREADS=4 DBC_IMAGE_PATH=images/normal.jpg \
//	  go run example.Goroutine_Http.go
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

// httpWorker runs in its own goroutine and owns its own HttpClient instance.
func httpWorker(username, password, imagePath string, decodeMode bool, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Each goroutine creates its own client — no shared connection state.
	client := dbc.NewHttpClient(username, password)
	defer client.Close()

	if !decodeMode {
		bal, err := client.GetBalance()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[HTTP goroutine %d] error: %v\n", workerID, err)
			return
		}
		fmt.Printf("[HTTP goroutine %d] balance: %.4f US cents\n", workerID, bal)
		return
	}

	captcha, err := client.Decode(imagePath, dbc.DefaultTimeout, map[string]string{})
	if err != nil {
		var ade *dbc.AccessDeniedException
		if errors.As(err, &ade) {
			fmt.Fprintf(os.Stderr, "[HTTP goroutine %d] access denied: %v\n", workerID, err)
		} else {
			fmt.Fprintf(os.Stderr, "[HTTP goroutine %d] error: %v\n", workerID, err)
		}
		return
	}
	if captcha == nil {
		fmt.Printf("[HTTP goroutine %d] timeout — no solution received\n", workerID)
		return
	}

	fmt.Printf("[HTTP goroutine %d] solved CAPTCHA %d: %s\n", workerID, captcha.CaptchaID, *captcha.Text)

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
	sanity := dbc.NewHttpClient(username, password)
	bal, err := sanity.GetBalance()
	sanity.Close()
	if err != nil {
		log.Fatalf("[HTTP main] credential check failed: %v", err)
	}
	fmt.Printf("[HTTP main] initial balance: %.4f US cents\n", bal)

	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	for i := 1; i <= goroutineCount; i++ {
		go httpWorker(username, password, imagePath, decodeMode, i, &wg)
	}

	wg.Wait()
	fmt.Println("[HTTP main] all goroutines finished.")
}
