package monzoroundup

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/tomwright/monzoroundup/auth"
	url2 "net/url"
	"strings"
	"errors"
)

type WebHook struct {
	AccountID string `json:"account_id"`
	ID        string `json:"id"`
	URL       string `json:"url"`
}

func ensureWebHookExists(accountID string, url string) error {
	hooks, err := listWebHooks(accountID)
	if err != nil {
		return err
	}
	for _, h := range hooks {
		if h.URL == url {
			log.Println("hook already exists: ", url)
			return nil
		}
		log.Println("hook doesn't match: ", h.URL, "... ", url)
	}

	newHook := WebHook{
		AccountID: accountID,
		URL:       url,
	}
	_, err = createWebHook(newHook)
	return err
}

func listWebHooks(accountID string) ([]WebHook, error) {
	url := fmt.Sprintf("https://api.monzo.com/webhooks?account_id=%s", url2.QueryEscape(accountID))

	req, err := http.NewRequest("GET", url, nil)
	auth.WrapRequest(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		return nil, errors.New("could not list webhooks: " + string(body))
	}

	result := map[string][]WebHook{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result["webhooks"], nil
}

func createWebHook(hook WebHook) (WebHook, error) {
	url := fmt.Sprintf("https://api.monzo.com/webhooks?account_id=%s", hook.AccountID)

	form := url2.Values{}
	form.Add("account_id", hook.AccountID)
	form.Add("url", hook.URL)

	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Println("form data: ", string(form.Encode()))

	auth.WrapRequest(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		return WebHook{}, errors.New("could not create webhook: " + string(body))
	}

	result := WebHook{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}
