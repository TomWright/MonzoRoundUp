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
		w.WriteHeader(400)
		w.Write([]byte("missing request body"))
		return
	}

	req := loginRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("could not decode request: %v\n", err)
		w.WriteHeader(400)
		w.Write([]byte("invalid request provided"))
		return
	}

	u, err := userModel.FetchByUserName(req.UserName)
	if err != nil {
		log.Printf("could not fetch user by username: %v\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}

	if u == nil {
		w.WriteHeader(401)
		w.Write([]byte("Invalid username"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(monzo.GenerateOauthUrl(u.OAuthID, u.ID)))
}
