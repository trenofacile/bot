package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	messenger "github.com/maciekmm/messenger-platform-go-sdk"
)

// PluginITF defines the interface of a bot plugin
type PluginITF interface {
	Reply(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage, mess *messenger.Messenger)
	ReplyToPostback(event messenger.Event, opts messenger.MessageOpts, pb messenger.Postback, mess *messenger.Messenger)
}

// Bot responds to Messenger messages from Facebook users.
type Bot struct {
	port    int
	mess    *messenger.Messenger
	plugins []PluginITF
}

// NewBot returns a newly initialized bot with the given parameters.
func NewBot(httpPort string, messengerAccessToken, messengerVerifyToken string) (*Bot, error) {
	b := new(Bot)

	port, err := strconv.Atoi(httpPort)
	if err != nil || port < 80 || port > 65535 {
		return nil, errors.New("Invalid http port given")
	}
	b.port = port

	if messengerAccessToken == "" || messengerVerifyToken == "" {
		return nil, errors.New("Invalid parameters")
	}

	mess := &messenger.Messenger{
		AccessToken:     messengerAccessToken,
		VerifyToken:     messengerVerifyToken,
		MessageReceived: b.messageReceived,
		Postback:        b.postbackHandler,
	}

	b.mess = mess

	return b, nil
}

// Start activates the bot.
func (b *Bot) Start() error {
	http.HandleFunc("/", b.mess.Handler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", b.port), nil)
	if err != nil {
		return err
	}

	return nil
}

// AddPlugin adds new features to the bot with the given plugin.
func (b *Bot) AddPlugin(plugin PluginITF) {
	b.plugins = append(b.plugins, plugin)
}

func (b *Bot) messageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	for _, plugin := range b.plugins {
		plugin.Reply(
			event,
			opts,
			msg,
			b.mess,
		)
	}
}

func (b *Bot) postbackHandler(event messenger.Event, opts messenger.MessageOpts, pb messenger.Postback) {
	for _, plugin := range b.plugins {
		plugin.ReplyToPostback(
			event,
			opts,
			pb,
			b.mess,
		)
	}
}
