package monzo

import (
	"net/http"
)

func wrapRequestWithToken(r *http.Request, token *Token) *http.Request {
	r.Header.Add("Authorization", token.TokenType+" "+token.AccessToken)
	return r
}
