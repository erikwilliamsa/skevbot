package chat

import (
	"fmt"

	"github.com/erikwilliamsa/skev-bot/twitch"
	twitchchat "github.com/gempir/go-twitch-irc/v2"
	"github.com/logrusorgru/aurora/v3"
)

type EventEngine struct {
	Client   *twitchchat.Client
	Handlers Handlers
}

func (ee EventEngine) Listen(c string) {
	ee.Client.OnUserJoinMessage(ee.Handlers.OnJoin(c))
}

type Handlers struct {
	repo          UserRepository
	userfollowers twitch.UserFollowers
}

func (h *Handlers) OnJoin(channel string) func(twitchchat.UserJoinMessage) {
	return func(m twitchchat.UserJoinMessage) {
		fmt.Printf("||")
		cu := h.repo.GetUser(channel, m.User)

		isFollowing, err := h.userfollowers.IsFollowing(cu.Name, channel)
		cu.IsFollower = isFollowing
		if err != nil {
			p("Unable to look up followers for ", cu.Name, err)
		}

		output := "%s Joined Chat: | Follower: %s | Returning: %s | Suspect: %s | \n"

		name, isfollower, returning, suspect := aurora.Cyan(cu.Name).String(), formatGood(isFollowing), formatGood(cu.Returning), "coming soon"

		fmt.Printf(output, name, isfollower, returning, suspect)

		if !cu.Returning {
			h.repo.SaveUser(channel, cu)
		}

	}
}

func formatGood(b bool) string {

	if b {
		return aurora.BgBlack(aurora.Bold(aurora.Green(b))).String()
	}

	return aurora.BgBlack(aurora.Bold(aurora.Red(b))).String()

}

func UserJoinHandler(uf twitch.UserFollowers, repo UserRepository) Handlers {
	return Handlers{
		repo:          repo,
		userfollowers: uf,
	}
}

type ChatUser struct {
	Name       string
	ID         string
	IsFollower bool
	Suspect    bool
	Returning  bool
}
