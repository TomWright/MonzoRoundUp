package auth

import (
	"time"
)

var (
	bindAddress            string
	state                  string
	oauthClientId          string
	oauthClientSecret      string
	ActiveToken            Token
	AccessTokenGranted     bool
	AccessTokenGrantedChan chan struct{}
)

func init() {
	AccessTokenGranted = false
	AccessTokenGrantedChan = make(chan struct{}, 100)
}

func Init(authBindAddr string, clientId string, clientSecret string) {
	bindAddress = authBindAddr
	state = "stateAbc123"
	oauthClientId = clientId
	oauthClientSecret = clientSecret
}

func Wait() {
	<-AccessTokenGrantedChan
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`
	ExpiresAt    time.Time
}
