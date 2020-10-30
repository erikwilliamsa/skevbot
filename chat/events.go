package chat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

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

		cu := h.repo.GetUser(m.User)

		isFollowing, err := h.userfollowers.IsFollowing(cu.Name, channel)

		if err != nil {
			p("Unable to look up followers for ", cu.Name, err)
		}

		output := "%s Joined Chat: | Follower: %s | Returning: %s | Suspect: %s | \n"

		name, isfollower, returning, suspect := aurora.Cyan(cu.Name).String(), formatGood(isFollowing), formatGood(cu.Returning), "coming soon"

		fmt.Printf(output, name, isfollower, returning, suspect)

		if !cu.Returning {
			h.repo.SaveUser(cu)
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

type UserRepository interface {
	GetUser(string) ChatUser
	SaveUser(ChatUser)
}

type users map[string]ChatUser

const fd = "/.skevbot/"

func NewMemUserRepositry() UserRepository {
	hd, _ := os.UserHomeDir()
	_, err := os.Stat(hd + fd)
	mutex := &sync.Mutex{}
	if err != nil {
		os.Mkdir(hd+fd, 0755)
		return userrepo{users{}, hd + fd + "users.json", mutex}
	}

	u, err := ioutil.ReadFile(hd + fd + "users.json")

	if err != nil {
		return userrepo{users{}, hd + fd + "users.json", mutex}
	}

	usrs := users{}
	json.Unmarshal(u, usrs)

	return userrepo{users: usrs, savepath: hd + fd + "users.json", mutex: mutex}

}

type userrepo struct {
	users    users
	savepath string
	mutex    *sync.Mutex
}

func (mup userrepo) GetUser(name string) ChatUser {

	cu, ok := mup.users[name]

	if !ok {
		cu.Name = name
		return cu
	}

	return cu
}

func (mup userrepo) SaveUser(cu ChatUser) {
	mup.mutex.Lock()
	defer mup.mutex.Unlock()

	mup.users[cu.Name] = cu
	ob, _ := json.Marshal(mup.users)
	ioutil.WriteFile(mup.savepath, ob, 0755)
}
