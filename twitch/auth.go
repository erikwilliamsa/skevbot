package twitch

import (
	"context"
	"fmt"
	"math"
	"time"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

type tokengenerator struct {
	oauth2Config *clientcredentials.Config
	RetryCount   int
	retries      int
}

//ClientCredAuth used for creating client credentials grant type tokens
func ClientCredAuth(cid, cs string) TokenGenerator {

	return &tokengenerator{
		oauth2Config: &clientcredentials.Config{
			ClientID:     cid,
			ClientSecret: cs,
			TokenURL:     twitch.Endpoint.TokenURL,
		},
		RetryCount: 3,
	}
}

func (tg *tokengenerator) Token() (string, error) {
	token, err := tg.oauth2Config.Token(context.Background())
	if err != nil {

		if tg.retries == tg.RetryCount {
			return "", err
		}

		fmt.Println("Retrying bearer token retrieval")
		time.Sleep(backoff(time.Second, float64(tg.retries)))
		tg.retries++
		return tg.Token()
	}
	return token.AccessToken, err
}

//TokenGenerator generates the appropriate OAuth2 Token
type TokenGenerator interface {
	Token() (string, error)
}

func backoff(t time.Duration, count float64) time.Duration {
	n := float64(t)
	n = n * math.Pow(2, count-1)
	return time.Duration(n)
}
