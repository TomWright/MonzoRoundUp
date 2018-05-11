package monzo

import "time"

type Token struct {
	ID           int64
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	UserID       int64  `json:"user_id"`
	ExpiresAt    time.Time
}

func (t Token) Expired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t Token) ExpiresWithin(length time.Duration) bool {
	return time.Now().Add(length).After(t.ExpiresAt)
}
