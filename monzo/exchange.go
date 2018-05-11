package monzo

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/url"
	"errors"
)

func ExchangeAuthCode(oauthClientID string, oauthClientSecret string, auth string) (*Token, error) {
	u := fmt.Sprintf("https://api.monzo.com/oauth2/token")

	form := url.Values{}

	form.Add("grant_type", "authorization_code")
	form.Add("client_id", oauthClientID)
	form.Add("client_secret", oauthClientSecret)
	form.Add("redirect_uri", oauthCallbackURL)
	form.Add("code", auth)

	req, err := http.NewRequest("POST", u, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		return nil, errors.New("could not exchange auth code:" + string(body))
	}

	t := Token{}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}

	t.ExpiresAt = time.Now().Add(time.Second * time.Duration(t.ExpiresIn))

	return &t, nil
}
