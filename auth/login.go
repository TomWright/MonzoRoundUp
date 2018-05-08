package auth

import (
	"net/http"
	"fmt"
)

const LoginEndpoint = "/oauth/login"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf(
		"https://auth.monzo.com/?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		oauthClientId,
		"http://"+bindAddress+CallbackEndpoint,
		state,
	), 301)
}
