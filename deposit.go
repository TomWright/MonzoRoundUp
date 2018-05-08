package monzoroundup

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"github.com/tomwright/monzoroundup/auth"
	"io/ioutil"
	"errors"
	"encoding/json"
	url2 "net/url"
)

type Deposit struct {
	PotID     string
	AccountID string
	Amount    int64
	DeDupeID  string
}

func depositToPotFromAccount(deposit Deposit) (Pot, error) {
	url := fmt.Sprintf("https://api.monzo.com/pots/%s/deposit", url2.QueryEscape(deposit.PotID))

	form := url2.Values{}
	form.Add("source_account_id", deposit.AccountID)
	form.Add("amount", fmt.Sprint(deposit.Amount))
	form.Add("dedupe_id", deposit.DeDupeID)

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
		return Pot{}, errors.New(fmt.Sprintf("could not deposit to pot `%s` from account `%s`: %s", deposit.PotID, deposit.AccountID, string(body)))
	}

	result := Pot{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}
