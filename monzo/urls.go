package monzo

import "fmt"

var (
	oauthCallbackURL string
)

func GenerateOauthUrl(clientID string, state string) string {
	return fmt.Sprintf(
		"https://auth.monzo.com/?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		clientID,
		oauthCallbackURL,
		state,
	)
}
