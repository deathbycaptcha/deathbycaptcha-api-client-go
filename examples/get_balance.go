// get_balance — DeathByCaptcha Go client example.
// Retrieves the account balance using both HTTP and Socket transports.
package main

import (
	"fmt"
	"log"

	dbc "github.com/deathbycaptcha/deathbycaptcha-api-client-go/deathbycaptcha"
)

func main() {
	// Put your DBC account username and password here.
	username := "your_username"
	password := "your_password"

	// -- HTTP client --
	httpClient := dbc.NewHttpClient(username, password)
	defer httpClient.Close()

	bal, err := httpClient.GetBalance()
	if err != nil {
		log.Fatalf("HTTP GetBalance error: %v", err)
	}
	fmt.Printf("HTTP balance: %.4f\n", bal)

	// -- Socket client --
	socketClient := dbc.NewSocketClient(username, password)
	defer socketClient.Close()

	bal2, err := socketClient.GetBalance()
	if err != nil {
		log.Fatalf("Socket GetBalance error: %v", err)
	}
	fmt.Printf("Socket balance: %.4f\n", bal2)
}
