package main

import (
	"fmt"
	"sync"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v2"
)

func main() {
	client := twitch.NewAnonymousClient()
	client.Join("justskev")
	users := Users{}

	l := sync.Mutex{}
	client.OnUserJoinMessage(func(m twitch.UserJoinMessage) {

		l.Lock()

		users[m.User] = User{
			Name:   m.User,
			Joined: time.Now(),
		}

		fmt.Printf("%s has joined your chat at %v.\n", m.User, users[m.User].Joined)
		l.Unlock()

	})
	client.OnUserPartMessage(func(m twitch.UserPartMessage) {
		l.Lock()

		if u, ok := users[m.User]; ok {
			u.Parted = time.Now()
			users[m.User] = u
			chattime := u.Parted.Sub(u.Joined)
			fmt.Printf("%s has left your chat at %v for a total of %v.\n", m.User, users[m.User].Parted, chattime)
		}

		l.Unlock()

	})
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

type User struct {
	Name   string
	Joined time.Time
	Parted time.Time
}

type Users map[string]User
