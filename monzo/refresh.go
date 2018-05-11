package monzo

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"time"
	"fmt"
	"strings"
	"net/url"
	"errors"
)

func RefreshToken(oauthClientID string, oauthClientSecret string, token *Token) (*Token, error) {
	u := fmt.Sprintf("https://api.monzo.com/oauth2/token")

	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("client_id", oauthClientID)
	form.Add("client_secret", oauthClientSecret)
	form.Add("refresh_token", token.RefreshToken)

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
		return nil, errors.New("could not refresh token: " + string(body))
	}

	t := Token{}

	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Fatal(err)
	}

	t.ExpiresAt = time.Now().Add(time.Second * time.Duration(t.ExpiresIn))

	return token, nil
}
