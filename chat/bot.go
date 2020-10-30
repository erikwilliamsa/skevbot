package chat

import (
	"fmt"

	twitchchat "github.com/gempir/go-twitch-irc/v2"
)

var p = fmt.Println

type Bot interface {
	Start() error
	Stop() error
}
type twitchchatitchBot struct {
	PrimaryChannel string
	client         *twitchchat.Client
}

func NewTwitchChatBot(primaryChannel string, c *twitchchat.Client) Bot {

	return &twitchchatitchBot{
		PrimaryChannel: primaryChannel,
		client:         c,
	}
}

func (tb *twitchchatitchBot) Start() error {

	tb.client.Join(tb.PrimaryChannel)

	tb.client.OnConnect(func() {
		p("Connected to twitch chat!", "Joined channel", tb.PrimaryChannel)
	})
	return tb.client.Connect()
}

func (tb *twitchchatitchBot) Stop() error {
	return tb.client.Disconnect()
}

func NewtwitchchatitchChatClient(user, oauth string) *twitchchat.Client {
	if user == "" || oauth == "" {
		p("User or Oauth token missing, connecting as Anonymous")
		return twitchchat.NewAnonymousClient()
	}

	return twitchchat.NewClient(user, oauth)

}
