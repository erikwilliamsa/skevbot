package main

import (
	"net/http"
	"os"

	"github.com/erikwilliamsa/skev-bot/chat"
	"github.com/erikwilliamsa/skev-bot/data/chatrepo"
	"github.com/erikwilliamsa/skev-bot/twitch"
)

func main() {

	cc := os.Getenv("CLIENT_ID")
	cs := os.Getenv("CLIENT_SECRET")
	c := os.Getenv("CHANNEL")
	if cc == "" || cs == "" || c == "" {
		panic("Must set environment variables: ")
	}
	oat, cu := os.Getenv("OAUTH_CHAT_TOKEN"), os.Getenv("TWITCH_USERNAME")

	ca := twitch.ClientCredAuth(cc, cs)
	f := twitch.NewUserFollowers(cc, ca, &http.Client{})
	user := c

	repo := chatrepo.NewMemUserRepositry()
	handlers := chat.UserJoinHandler(f, repo)
	tcc := chat.NewtwitchchatitchChatClient(cu, oat)
	engine := chat.EventEngine{Client: tcc, Handlers: handlers}

	bot := chat.NewTwitchChatBot(user, tcc)

	engine.Listen(user)
	bot.Start()

}
