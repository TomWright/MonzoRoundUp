package worker

import (
	"time"
	"github.com/tomwright/monzoroundup/token"
	"log"
	"github.com/tomwright/monzoroundup/monzo"
	"github.com/tomwright/monzoroundup/user"
)

var RefreshToken *refreshToken

type refreshToken struct {
	ticker *time.Ticker
	T      time.Duration
}

func (i *refreshToken) work(tokenModel token.Model, userModel user.Model) {
	tokens, err := tokenModel.FetchMostRecentForAllUsers()
	if err != nil {
		log.Printf("could not fetch tokens: %s\n", err)
		return
	}

	for userId, t := range tokens {
		if t.ExpiresWithin(time.Minute * 5) {
			u, err := userModel.FetchByID(userId)
			if err != nil {
				log.Printf("token `%d` not refreshed. could not fetch user: %s\n", t.ID, err)
				return
			}

			newToken, err := monzo.RefreshToken(u.OAuthID, u.OAuthSecret, t)
			if err != nil {
				log.Printf("token `%d` not refreshed. could not refresh token: %s\n", t.ID, err)
				return
			}

			_, err = tokenModel.Insert(userId, newToken)
			if err != nil {
				log.Printf("token `%d` not refreshed. could not insert new token: %s\n", t.ID, err)
				return
			}
		}
	}
}

func (i *refreshToken) Work(tokenModel token.Model, userModel user.Model) {
	i.ticker = time.NewTicker(i.T)
	i.work(tokenModel, userModel)
	for {
		_, open := <-i.ticker.C
		if ! open {
			return
		}

		i.work(tokenModel, userModel)
	}
}
