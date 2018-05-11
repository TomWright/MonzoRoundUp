package monzo

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"errors"
	"encoding/json"
)

type Pot struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Style    string `json:"style"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
	Deleted  bool   `json:"deleted"`
}

func ListPots(token *Token) ([]Pot, error) {
	url := fmt.Sprintf("https://api.monzo.com/pots")

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

	result := map[string][]Pot{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result["pots"], nil
}
