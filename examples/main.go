package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	kitesession "github.com/nsvirk/gokitesession"
	kiteconnect "github.com/nsvirk/gokiteticker/kiteconnect"
	kitemodels "github.com/nsvirk/gokiteticker/models"
	kiteticker "github.com/nsvirk/gokiteticker/ticker"
)

var (
	// kite user credentials
	userId     string = os.Getenv("KITE_USER_ID")
	password   string = os.Getenv("KITE_PASSWORD")
	totpSecret string = os.Getenv("KITE_TOTP_SECRET")

	// kiteticker instrument tokens to subscribe
	instTokens []uint32 = append([]uint32{}, 256265, 264969, 5633, 779521, 408065, 738561, 895745)
)

var (
	ticker *kiteticker.Ticker
)

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	fmt.Println("Connected")
	fmt.Println("Subscribing to", instTokens)
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("")
	err := ticker.Subscribe(instTokens)
	if err != nil {
		fmt.Println("err: ", err)
	}
	// Set subscription mode for given list of tokens
	// Default mode is Quote
	err = ticker.SetMode(kiteticker.ModeLTP, instTokens)
	// err = ticker.SetMode(kiteticker.ModeFull, instTokens)
	if err != nil {
		fmt.Println("err: ", err)
	}
}

// Triggered when tick is recevived
func onTick(tick kitemodels.Tick) {
	// fmt.Println("Tick: ", tick)
	// tickJson, err := json.MarshalIndent(tick, "", " \t")
	tickJson, err := json.Marshal(tick)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(tickJson))
	fmt.Println("--------------------------------------------------------------")
}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	fmt.Printf("Order: %s", order.OrderID)
}

func main() {

	// Create a new Kite session instance
	ks := kitesession.New(userId)

	// Set debug mode
	ks.SetDebug(true)

	// Generate totp value
	totpValue, err := ks.GenerateTotpValue(totpSecret)
	if err != nil {
		fmt.Printf("Error generating totp value: %v", err)
		return
	}
	// Check the inputs values
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("Kite User")
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("User ID     	: ", userId)
	fmt.Println("Password     	: ", password)
	fmt.Println("Totp Value  	: ", totpValue)
	fmt.Println("")

	// Get kite session data
	session, err := ks.GenerateSession(password, totpValue)
	if err != nil {
		fmt.Printf("Error generating session: %v", err)
		return
	}

	fmt.Println("--------------------------------------------------------------")
	fmt.Println("Kite Session")
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("user_id     	: ", session.UserId)
	fmt.Println("username   	: ", session.Username)
	fmt.Println("enctoken    	: ", session.Enctoken[0:36], "...")
	fmt.Println("login_time  	: ", session.LoginTime)
	fmt.Println("")
	// fmt.Println(session)

	fmt.Println("--------------------------------------------------------------")
	fmt.Println("Kite Ticker")
	fmt.Println("--------------------------------------------------------------")

	// Create new Kite ticker instance
	ticker = kiteticker.New(userId, session.Enctoken)

	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)
	ticker.OnOrderUpdate(onOrderUpdate)

	// Start the connection
	ticker.Serve()
}
