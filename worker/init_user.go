package worker

import (
	"github.com/tomwright/monzoroundup/user"
	"log"
	"github.com/tomwright/monzoroundup/monzo"
	"github.com/tomwright/monzoroundup/token"
	"fmt"
	"os"
	"github.com/tomwright/monzoroundup/transport/http/webhook"
)

var InitUser *initUser

type initUser struct {
	C chan user.User
}

func (i *initUser) Work(tokenModel token.Model) {
	users := i.C
	log.Println("init user worker started")

	expectedURL := fmt.Sprintf("%s://%s%s", os.Getenv("HTTP_PROTOCOL"), os.Getenv("PUBLIC_DOMAIN"), webhook.ReceiveEventEndpoint)

	for {
		u, open := <-users
		if ! open {
			log.Println("init u worker stopping")
			return
		}

		t, err := tokenModel.FetchByUserID(u.ID)
		if err != nil {
			log.Printf("could not fetch user token: %s\n", err)
			continue
		}

		accounts, err := monzo.ListAccounts(t)
		if err != nil {
			log.Printf("could not get user accounts: %s\n", err)
			continue
		}

		for _, account := range accounts {
			hooks, err := monzo.ListWebHooks(t, account.ID)
			if err != nil {
				log.Printf("could not fetch webhooks: %s\n", err)
				continue
			}

			if account.Closed {
				log.Println("skipping closed account")
				continue
			}

			foundHook := false
		hookLoop:
			for _, hook := range hooks {
				if hook.URL == expectedURL {
					log.Printf("Hook `%s` was found as expected", expectedURL)
					foundHook = true
					break hookLoop
				}
			}

			if ! foundHook {
				log.Printf("Hook `%s` was not found. it will be created now", expectedURL)

				_, err := monzo.CreateWebHook(t, monzo.WebHook{
					AccountID: account.ID,
					URL:       expectedURL,
				})
				if err != nil {
					log.Printf("could not create webhook: %s\n", err)
					continue
				}
			}
		}
	}
}
