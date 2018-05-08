package monzoroundup

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"
	"github.com/tomwright/monzoroundup/auth"
	"fmt"
)

var (
	bindAddress string
	httpClient  *http.Client
)

const (
	WebHookEndpoint = "/webhook"
)

func Init(bindAddr string, clientId string, clientSecret string) {
	bindAddress = bindAddr

	httpClient = http.DefaultClient

	auth.Init(bindAddress, clientId, clientSecret)
}

func Listen() {
	fin := make(chan struct{})

	go func() {
		http.HandleFunc(WebHookEndpoint, handleWebHook)
		http.HandleFunc(auth.LoginEndpoint, auth.LoginHandler)
		http.HandleFunc(auth.CallbackEndpoint, auth.CallbackHandler)

		fmt.Println("Login at:", "http://"+bindAddress+auth.LoginEndpoint)

		err := http.ListenAndServe(bindAddress, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		log.Println("Waiting for auth")
		auth.Wait()
		log.Println("Successfully authorized")

		err := prepareAccounts()
		if err != nil {
			log.Println(err)
		}

		fin <- struct{}{}
	}()

	<-fin
	// start the Event listener
	waitForEvents()
}

func prepareAccounts() error {
	accounts, err := listAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if account.Closed {
			continue
		}
		err = ensureWebHookExists(account.ID, "http://"+bindAddress+WebHookEndpoint)
		if err != nil {
			log.Fatal("Could not ensure webhook exists: ", err)
		}
	}

	return nil
}

func handleWebHook(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		sendResponse(w, map[string]interface{}{
			"error": "Missing request body",
		}, 400)
		return
	}

	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendResponse(w, map[string]interface{}{
			"error": "Could not read request body",
		}, 400)
		return
	}

	event := Event{}
	err = json.Unmarshal(bodyData, &event)
	if err != nil {
		sendResponse(w, map[string]interface{}{
			"error": "Could not parse event",
		}, 400)
	}

	eventChan <- event

	sendResponse(w, map[string]interface{}{
		"error": "Event accepted",
	}, 202)

}

func sendResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonData)
	w.WriteHeader(statusCode)
}
