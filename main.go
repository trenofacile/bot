package main

import (
	"log"
	"os"

	"github.com/trenofacile/bot/plugins"
)

func main() {
	b, err := NewBot(
		os.Getenv("HTTP_PORT"),
		os.Getenv("MESSENGER_ACCESSTOKEN"),
		os.Getenv("MESSENGER_VERIFYTOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	b.AddPlugin(&plugins.EchoPlugin{})

	err = b.Start()
	if err != nil {
		log.Fatal(err)
	}
}
