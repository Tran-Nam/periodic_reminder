package main

import (
	"log"

	"github.com/Tran-Nam/periodic_remind/pkg/telegram"
	"github.com/Tran-Nam/periodic_remind/pkg/common"
)

func main() {
	// load .env (optional) and send a test message
	_ = common.LoadDotEnv(".env")
	if err := telegram.SendMessageFromEnv("Good afternoon!"); err != nil {
		log.Fatal(err)
	}
}
