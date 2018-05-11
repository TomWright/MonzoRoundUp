package monzo

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
)

type Account struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Created     string `json:"created"`
	Closed      bool   `json:"closed"`
	Type        string `json:"type"`
}

func ListAccounts(token *Token) ([]*Account, error) {
	url := fmt.Sprintf("https://api.monzo.com/accounts")

	req, err := http.NewRequest("GET", url, nil)
	req = wrapRequestWithToken(req, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, errors.New("could not list accounts: " + string(body))
	}

	result := map[string][]*Account{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result["accounts"], nil
}
