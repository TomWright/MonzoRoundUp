package auth

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	"time"
	url2 "net/url"
	"strings"
)

func exchangeAuthCode(authCode string) {
	url := fmt.Sprintf("https://api.monzo.com/oauth2/token")

	form := url2.Values{}

	form.Add("grant_type", "authorization_code")
	form.Add("client_id", oauthClientId)
	form.Add("client_secret", oauthClientSecret)
	form.Add("redirect_uri", "http://"+bindAddress+CallbackEndpoint)
	form.Add("code", authCode)

	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		log.Println("could not exchange auth code", string(body))
		return
	}

	t := Token{}

	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Fatal(err)
	}

	t.ExpiresAt = time.Now().Add(time.Second * time.Duration(t.ExpiresIn))

	ActiveToken = t

	log.Println("Exchanged ActiveToken ("+ActiveToken.AccessToken+"). Expires at", ActiveToken.ExpiresAt.Format(time.RFC3339))
	AccessTokenGrantedChan <- struct{}{}
	AccessTokenGranted = true
}

func refreshToken() {
	url := fmt.Sprintf("https://api.monzo.com/oauth2/token")

	form := url2.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("client_id", oauthClientId)
	form.Add("client_secret", oauthClientSecret)
	form.Add("refresh_token", ActiveToken.RefreshToken)

	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		log.Println("could not refresh token", string(body))
		return
	}

	t := Token{}

	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Fatal(err)
	}

	t.ExpiresAt = time.Now().Add(time.Second * time.Duration(t.ExpiresIn))

	ActiveToken = t

	log.Println("Refreshed token ("+ActiveToken.AccessToken+"). Expires at", ActiveToken.ExpiresAt.Format(time.RFC3339))
	AccessTokenGrantedChan <- struct{}{}
	AccessTokenGranted = true
}

func WrapRequest(r *http.Request) *http.Request {
	if ActiveToken.ExpiresAt.Before(time.Now()) {
		AccessTokenGranted = false
		refreshToken()
	}

	r.Header.Add("Authorization", ActiveToken.TokenType+" "+ActiveToken.AccessToken)
	return r
}
