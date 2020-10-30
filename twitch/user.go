package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	TwitchAPIURL = "https://api.twitch.tv/"
	UserPath     = "helix/users"
	FollowsPath  = UserPath + "/follows"
)

//UserFollowers is used to determine information about user followers
type UserFollowers interface {
	//IsFollowing is the User following user?
	IsFollowing(from string, to string) (bool, error)
	//HasFollowers does the user have followers?
	HasFollowers(string) (bool, error)
	//GetFollowers Get a list of Followers
	GetFollowers(string) (Followers, error)
}

type userfollowers struct {
	clientID       string
	tokenGenerator TokenGenerator
	Client         *http.Client
}

func (uf userfollowers) IsFollowing(from, to string) (bool, error) {
	ids, err := uf.getIDS(from, to)
	if err != nil {
		log.Println("IDS:", err)
		return false, err
	}
	url := fmt.Sprintf("%s%s?from_id=%s&to_id=%s", TwitchAPIURL, FollowsPath, ids[from], ids[to])
	b, err := uf.get(url)

	if err != nil {
		log.Println("UF Get:", err)

		return false, err
	}

	var f Followers
	err = json.Unmarshal(b, &f)

	if err != nil {
		log.Println(err)

		return false, err
	}
	return f.Total > 0, nil
}

func (uf userfollowers) HasFollowers(user string) (bool, error) {
	url := fmt.Sprintf("%s%s?to_id=%s", TwitchAPIURL, FollowsPath, user)
	body, err := uf.get(url)

	if err != nil {
		return false, err
	}
	var f Followers
	err = json.Unmarshal(body, &f)

	if err != nil {
		return false, err
	}

	return f.Total > 0, nil

}

func (uf userfollowers) GetFollowers(user string) (Followers, error) {
	return Followers{}, nil
}

func (uf userfollowers) get(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err

	}

	token, err := uf.tokenGenerator.Token()

	if err != nil {
		return []byte{}, err

	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("client-id", uf.clientID)

	res, err := uf.Client.Do(req)

	if res.StatusCode != http.StatusOK {
		return []byte{}, errorMessage(res)
	}
	if err != nil {
		return []byte{}, err

	}
	return ioutil.ReadAll(res.Body)

}

func (uf userfollowers) getIDS(users ...string) (map[string]string, error) {

	uids := map[string]string{}
	url := fmt.Sprintf("%s%s?", TwitchAPIURL, UserPath)

	for i, u := range users {

		pattern := "%s&login=%s"
		if i == 0 {
			pattern = "%slogin=%s"
		}
		url = fmt.Sprintf(pattern, url, u)
	}

	b, err := uf.get(url)

	if err != nil {
		log.Println("Error looking up IDs", err)
		return uids, err
	}

	usersList := UsersList{}
	json.Unmarshal(b, &usersList)

	for _, u := range usersList.Users {
		uids[strings.ToLower(u.DisplayName)] = u.ID
	}
	return uids, nil
}

func NewUserFollowers(clientID string, tg TokenGenerator, c *http.Client) UserFollowers {
	return &userfollowers{tokenGenerator: tg, clientID: clientID, Client: c}
}

func errorMessage(res *http.Response) error {
	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	tem := TwitchErrorMessage{}
	json.Unmarshal(b, &tem)

	return fmt.Errorf(string(b))

}
