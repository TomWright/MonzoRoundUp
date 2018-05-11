package auth

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"github.com/tomwright/monzoroundup/monzo"
	"github.com/julienschmidt/httprouter"
)

type loginRequest struct {
	UserName string `json:"username"`
}

// Fetch the user with the provided username and redirect to the monzo oauth login
func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("missing request body"))
		w.WriteHeader(400)
		return
	}

	req := loginRequest{}
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

	if u == nil {
		w.Write([]byte("Invalid username"))
		w.WriteHeader(401)
		return
	}

	http.Redirect(w, r, monzo.GenerateOauthUrl(u.OAuthID, u.ID), 301)
}
