package main

import (
	"net/http"
	"os"

	"github.com/erikwilliamsa/skev-bot/chat"
	"github.com/erikwilliamsa/skev-bot/twitch"
)

func main() {

	cc := os.Getenv("CLIENT_ID")
	cs := os.Getenv("CLIENT_SECRET")
	c := os.Getenv("CHANNEL")
	if cc == "" || cs == "" || c == "" {
		panic("Must set environment variables: ")
	}
	ca := twitch.ClientCredAuth(cc, cs)
	f := twitch.NewUserFollowers(cc, ca, &http.Client{})
	user := c

	repo := chat.NewMemUserRepositry()
	handlers := chat.UserJoinHandler(f, repo)
	tcc := chat.NewtwitchchatitchChatClient("", "")
	engine := chat.EventEngine{Client: tcc, Handlers: handlers}

	bot := chat.NewTwitchChatBot(user, tcc)

	engine.Listen(user)
	bot.Start()

}
