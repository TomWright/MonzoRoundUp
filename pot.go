package monzoroundup

import (
	"fmt"
	"net/http"
	"github.com/tomwright/monzoroundup/auth"
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

func listPots() ([]Pot, error) {
	url := fmt.Sprintf("https://api.monzo.com/pots")

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
		return nil, errors.New("could not list accounts: " + string(body))
	}

	result := map[string][]Pot{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result["pots"], nil
}

func getRoundUpPot() (Pot, error) {
	pots, err := listPots()
	if err != nil {
		return Pot{}, err
	}

	for _, pot := range pots {
		if pot.Name == "RoundUp" {
			return pot, nil
		}
	}

	return Pot{}, errors.New("no roundup pot exists")
}
