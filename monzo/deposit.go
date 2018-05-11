package monzo

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"io/ioutil"
	"errors"
	"encoding/json"
	"net/url"
)

func DepositToPotFromAccount(token *Token, potID string, accountID string, amount int64, deDupeID string) (*Pot, error) {
	u := fmt.Sprintf("https://api.monzo.com/pots/%s/deposit", url.QueryEscape(potID))

	form := url.Values{}
	form.Add("source_account_id", accountID)
	form.Add("amount", fmt.Sprint(amount))
	form.Add("dedupe_id", deDupeID)

	req, err := http.NewRequest("POST", u, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req = wrapRequestWithToken(req, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		return nil, errors.New(fmt.Sprintf("could not deposit to pot `%s` from account `%s`: %s", potID, accountID, string(body)))
	}

	result := Pot{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	return &result, nil
}
