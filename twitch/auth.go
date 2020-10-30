package twitch

import (
	"context"
	"log"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

type tokengenerator struct {
	oauth2Config *clientcredentials.Config
}

//ClientCredAuth used for creating client credentials grant type tokens
func ClientCredAuth(cid, cs string) TokenGenerator {

	return &tokengenerator{
		oauth2Config: &clientcredentials.Config{
			ClientID:     cid,
			ClientSecret: cs,
			TokenURL:     twitch.Endpoint.TokenURL,
		},
	}
}

func (tg *tokengenerator) Token() (string, error) {
	token, err := tg.oauth2Config.Token(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return token.AccessToken, err
}

//TokenGenerator generates the appropriate OAuth2 Token
type TokenGenerator interface {
	Token() (string, error)
}
