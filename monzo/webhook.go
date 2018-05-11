package monzo

import (
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strings"
	"net/url"
)

type WebHook struct {
	AccountID string `json:"account_id"`
	ID        string `json:"id"`
	URL       string `json:"url"`
}

func ListWebHooks(token *Token, accountID string) ([]*WebHook, error) {
	u := fmt.Sprintf("https://api.monzo.com/webhooks?account_id=%s", url.QueryEscape(accountID))

	req, err := http.NewRequest("GET", u, nil)
	req = wrapRequestWithToken(req, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		return nil, errors.New("could not list webhooks: " + string(body))
	}

	result := map[string][]*WebHook{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result["webhooks"], nil
}

func CreateWebHook(token *Token, hook WebHook) (*WebHook, error) {
	u := fmt.Sprintf("https://api.monzo.com/webhooks?account_id=%s", hook.AccountID)

	form := url.Values{}
	form.Add("account_id", hook.AccountID)
	form.Add("url", hook.URL)

	req, err := http.NewRequest("POST", u, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Println("form data: ", string(form.Encode()))
	req = wrapRequestWithToken(req, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, errors.New("could not create webhook: " + string(body))
	}

	result := WebHook{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
