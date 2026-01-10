package main

import (
	"log"

	"github.com/Tran-Nam/periodic_remind/pkg/telegram"
)

func main() {
	if err := telegram.SendMessageFromEnv("Good afternoon!"); err != nil {
		log.Fatal(err)
	}
}
