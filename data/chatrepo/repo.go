package chatrepo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/erikwilliamsa/skev-bot/chat"
)

const fd = "/.skevbot/v1/"

func NewMemUserRepositry() chat.UserRepository {
	setup()
	return &userrepo{mutex: &sync.Mutex{}}

}

var data repodata

func setup() {

	hd, err := os.UserHomeDir()

	if err != nil {
		panic("Unable to access data home directory: " + err.Error())
	}

	p := path.Join(hd, fd)
	if data.userdata == nil {
		fmt.Println("Initilizing data")
		data = repodata{
			path:     p,
			datafile: "users.json",
			userdata: make(map[string]users),
		}
	}
	_, err = os.Stat(p)

	if err != nil {
		if err := os.MkdirAll(p, 0755); err != nil {
			panic("Cannot create data directory: " + err.Error())
		}
	}

	fp := path.Join(data.path, data.datafile)

	_, err = os.Stat(fp)
	if err == nil {
		b, err := ioutil.ReadFile(fp)
		if err != nil {
			panic("Cannot read data file: " + err.Error())
		}

		err = json.Unmarshal(b, &data.userdata)

		if err != nil {
			log.Fatal("Cannot parse data: ", err.Error())
		}
	}

}

type userrepo struct {
	mutex *sync.Mutex
}

func (mup userrepo) GetUser(name, channel string) chat.ChatUser {
	_, ok := data.userdata[channel]

	if !ok {
		data.userdata[channel] = users{}
	}
	cu, ok := data.userdata[channel][name]

	if !ok {
		cu.Name = name
		mup.SaveUser(channel, cu)
		return cu
	}

	return cu
}

func (mup userrepo) SaveUser(channel string, cu chat.ChatUser) {
	mup.mutex.Lock()
	defer mup.mutex.Unlock()
	cu.Returning = true

	data.userdata[channel][cu.Name] = cu
	data.Save()

}

type users map[string]chat.ChatUser

type repodata struct {
	userdata map[string]users
	path     string
	datafile string
}

func (rd repodata) Save() {
	fp := path.Join(data.path, data.datafile)

	b, err := json.Marshal(data.userdata)

	if err != nil {
		log.Println("Unable to parse data", err)
	}

	err = ioutil.WriteFile(fp, b, 0755)
	if err != nil {
		log.Println("unable to save data")
	}
}
