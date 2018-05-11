package monzo

import (
	"net/http"
	"log"
)

func wrapRequestWithToken(r *http.Request, token *Token) *http.Request {
	authHeader := token.TokenType + " "
	if token.TokenType == "BEARER" {
		authHeader = "Bearer "
	}
	authHeader += token.AccessToken
	log.Println("Adding auth header: " + authHeader)
	r.Header.Add("Authorization", authHeader)
	return r
}
