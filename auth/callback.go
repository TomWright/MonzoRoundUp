package auth

import (
	"net/http"
	"log"
	"fmt"
)

const CallbackEndpoint = "/oauth/callback"

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	authCode := r.Form.Get("code")
	requestState := r.Form.Get("state")

	if requestState != requestState {
		log.Println("invalid state. expected `", state, "`, got `", requestState, "`")
	}

	fmt.Println("got auth code", authCode)

	exchangeAuthCode(authCode)
}
