package auth

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"github.com/julienschmidt/httprouter"
)

type registerRequest struct {
	UserName    string `json:"username"`
	OAuthID     string `json:"oauthClientId"`
	OAuthSecret string `json:"oauthClientSecret"`
}

// Fetch the user with the provided username and redirect to the monzo oauth register
func registerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("missing request body"))
		w.WriteHeader(400)
		return
	}

	req := registerRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("could not decode request: %v\n", err)
		w.Write([]byte("invalid request provided"))
		w.WriteHeader(400)
		return
	}

	u, err := userModel.FetchByUserName(req.UserName)
	if err != nil {
		log.Printf("could not fetch user by username: %v\n", err)
		w.Write([]byte("Internal server error"))
		w.WriteHeader(500)
		return
	}

	if u != nil {
		w.Write([]byte("Username already taken"))
		w.WriteHeader(401)
		return
	}

	u, err = userModel.Create(req.UserName, req.OAuthID, req.OAuthSecret)
	if err != nil {
		log.Printf("could not create user: %v\n", err)
		w.Write([]byte("Internal server error"))
		w.WriteHeader(500)
		return
	}

	http.Redirect(w, r, LoginEndpoint, 301)
}
