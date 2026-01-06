package main

import (
	"log"

	"github.com/Tran-Nam/periodic_remind/pkg/telegram"
)

func main() {
	// load .env (optional) and send a test message
	_ = telegram.LoadDotEnv(".env")
	if err := telegram.SendMessageFromEnv("Hello from periodic_remind"); err != nil {
		log.Fatal(err)
	}
}
