package main

import (
	"log"
	"os"

	"github.com/trenofacile/bot/plugins"
	"github.com/trenofacile/bot/witai"
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

	witAIClient, err := witai.NewClient(os.Getenv("WITAI_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	b.AddPlugin(plugins.NewWitAIPlugin(witAIClient))

	err = b.Start()
	if err != nil {
		log.Fatal(err)
	}
}
