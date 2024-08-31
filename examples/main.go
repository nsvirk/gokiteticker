// Package main implements a Kite ticker client that subscribes to market data
// and processes incoming ticks.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	kitesession "github.com/nsvirk/gokitesession"
	kiteticker "github.com/nsvirk/gokiteticker"
	kitemodels "github.com/nsvirk/gokiteticker/models"
)

// Global variables
var (
	// instTokens is a slice of instrument tokens to subscribe to
	instTokens = []uint32{256265, 264969, 5633, 779521, 408065, 738561, 895745}

	// ticker is the main Kite ticker instance
	ticker *kiteticker.Ticker
)

// Config holds the configuration data for the Kite session.
type Config struct {
	UserID     string
	Password   string
	TOTPSecret string
}

// main is the entry point of the application
func main() {
	// Get the Kite session
	session, err := getSession()
	if err != nil {
		log.Fatalf("Failed to get session: %v", err)
	}

	// Get the userId and enctoken from the session
	userId := session.UserID
	enctoken := session.Enctoken

	// Create new Kite ticker instance
	ticker = kiteticker.New(userId, enctoken)

	// Set up callbacks
	setupCallbacks()

	// Start the connection
	// Serve starts the connection to ticker server. Since its blocking its
	// recommended to use it in a go routine.
	ticker.Serve()

	// Alternatively, you can use ServeWithContext to start the connection
	// go ticker.ServeWithContext(context.Background())
	// select {}
}

// setupCallbacks assigns all the necessary callbacks to the ticker
func setupCallbacks() {
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)
	ticker.OnOrderUpdate(onOrderUpdate)
}

// onError is triggered when any error is raised
func onError(err error) {
	log.Printf("Error: %v", err)
}

// onClose is triggered when websocket connection is closed
func onClose(code int, reason string) {
	log.Printf("Connection closed: code=%d, reason=%s", code, reason)
}

// onConnect is triggered when connection is established and ready to send and accept data
func onConnect() {
	log.Println("Connected")
	log.Printf("Subscribing to %v", instTokens)
	log.Println("--------------------------------------------------------------")

	if err := ticker.Subscribe(instTokens); err != nil {
		log.Printf("Subscription error: %v", err)
		return
	}

	// Set subscription mode for given list of tokens
	// Default mode is Quote
	if err := ticker.SetMode(kiteticker.ModeLTP, instTokens); err != nil {
		log.Printf("SetMode error: %v", err)
	}
}

// onTick is triggered when a tick is received
func onTick(tick kitemodels.Tick) {
	tickJSON, err := json.Marshal(tick)
	if err != nil {
		log.Printf("Error marshalling tick: %v", err)
		return
	}

	printTicks(tickJSON)
}

// onReconnect is triggered when reconnection is attempted
func onReconnect(attempt int, delay time.Duration) {
	log.Printf("Reconnect attempt %d in %.2fs", attempt, delay.Seconds())
}

// onNoReconnect is triggered when maximum number of reconnect attempts is reached
func onNoReconnect(attempt int) {
	log.Printf("Maximum number of reconnect attempts reached: %d", attempt)
}

// onOrderUpdate is triggered when an order update is received
func onOrderUpdate(order kiteticker.Order) {
	log.Printf("Order update received: OrderID=%s", order.OrderID)
}

// getSession retrieves the Kite session
func getSession() (*kitesession.Session, error) {
	config, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	ks := kitesession.New()
	ks.SetDebug(false)

	totpValue, err := kitesession.GenerateTOTPValue(config.TOTPSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP value: %w", err)
	}

	printInputValues(config.UserID, config.Password, totpValue)

	session, err := ks.GenerateSession(config.UserID, config.Password, totpValue)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session: %w", err)
	}

	printSessionInfo(session)

	return session, nil
}

// getConfig retrieves the configuration from environment variables
func getConfig() (*Config, error) {
	userID := os.Getenv("KITE_USER_ID")
	password := os.Getenv("KITE_PASSWORD")
	totpSecret := os.Getenv("KITE_TOTP_SECRET")

	if userID == "" || password == "" || totpSecret == "" {
		return nil, fmt.Errorf("KITE_USER_ID, KITE_PASSWORD, and KITE_TOTP_SECRET environment variables must be set")
	}

	return &Config{
		UserID:     userID,
		Password:   password,
		TOTPSecret: totpSecret,
	}, nil
}

// printInputValues prints the input values used for authentication
func printInputValues(userID, password, totpValue string) {
	log.Println("--------------------------------------------------------------")
	log.Println("User Inputs")
	log.Println("--------------------------------------------------------------")
	log.Printf("User ID      : %s", userID)
	log.Printf("Password     : %s", password)
	log.Printf("TOTP Value   : %s\n", totpValue)
}

// printSessionInfo prints the session information
func printSessionInfo(session *kitesession.Session) {
	log.Println("--------------------------------------------------------------")
	log.Println("Kite Session")
	log.Println("--------------------------------------------------------------")
	log.Printf("user_id      : %s", session.UserID)
	log.Printf("enctoken     : %s", session.Enctoken)
	log.Printf("login_time   : %s", session.LoginTime)
}

// printTicks prints the ticker information
func printTicks(tickJSON []byte) {
	log.Println("--------------------------------------------------------------")
	log.Println("Tick Data")
	log.Println("--------------------------------------------------------------")

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, tickJSON, "", "\t"); err != nil {
		log.Printf("Error processing tick data: %v", err)
		return
	}

	log.Println(prettyJSON.String())
	log.Println("--------------------------------------------------------------")
}
