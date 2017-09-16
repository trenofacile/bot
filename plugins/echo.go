package plugins

import (
	"encoding/json"
	"fmt"
	"log"

	messenger "github.com/maciekmm/messenger-platform-go-sdk"
	"github.com/maciekmm/messenger-platform-go-sdk/template"
)

// EchoPlugin echoes whatever transits for the bot
type EchoPlugin struct{}

// Reply echoes the message
func (p *EchoPlugin) Reply(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage, mess *messenger.Messenger) {
	b, _ := json.Marshal(msg)
	log.Println(string(b))

	mq := messenger.MessageQuery{}
	mq.RecipientID(opts.Sender.ID)

	mq.Template(template.GenericTemplate{
		Title:    fmt.Sprintf("You wrote:"),
		Subtitle: msg.Text,
	})

	_, err := mess.SendMessage(mq)
	if err != nil {
		fmt.Println(err)
	}
}

// ReplyToPostback echoes the postback object
func (p *EchoPlugin) ReplyToPostback(event messenger.Event, opts messenger.MessageOpts, pb messenger.Postback, mess *messenger.Messenger) {
	b, _ := json.Marshal(pb)
	log.Println(string(b))
}
