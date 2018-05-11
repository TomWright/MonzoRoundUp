package auth

import (
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"github.com/tomwright/monzoroundup/monzo"
	"strconv"
	"github.com/tomwright/monzoroundup/worker"
)

// Fetch the user with the provided username and redirect to the monzo oauth callback
func callbackHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("could not parse form: %v\n", err)
		w.WriteHeader(400)
		w.Write([]byte("invalid request provided"))
		return
	}

	authCode := r.Form.Get("code")
	requestState := r.Form.Get("state")

	userId, err := strconv.ParseInt(requestState, 10, 64)
	if err != nil {
		log.Printf("Invalid state. Not int. Got `%v` of type `%t`\n", requestState, requestState)
		w.WriteHeader(400)
		w.Write([]byte("invalid state"))
		return
	}

	u, err := userModel.FetchByID(userId)
	if err != nil {
		log.Printf("Could not fetch user: %s\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}

	if u == nil {
		log.Printf("Invalid state. Invalid user ID `%d`\n", userId)
		w.WriteHeader(400)
		w.Write([]byte("invalid state"))
		return
	}

	fmt.Println("got auth code", authCode)

	token, err := monzo.ExchangeAuthCode(u.OAuthID, u.OAuthSecret, authCode)
	if err != nil {
		log.Printf("could not exchange auth code `%s`: %s", authCode, err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
	}

	_, err = tokenModel.Insert(userId, token)
	if err != nil {
		log.Printf("could not insert token: %s\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}

	worker.InitUser.C <- *u

	w.WriteHeader(200)
	w.Write([]byte("Success"))
}
