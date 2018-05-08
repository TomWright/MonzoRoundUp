package monzoroundup

import (
	"errors"
	"log"
	"github.com/tomwright/monzoroundup/auth"
	"encoding/json"
	"fmt"
)

var eventChan chan Event

func waitForEvents() {
	eventChan = make(chan Event, 10)

	log.Println("Waiting for events")

eventsLoop:
	for {
		event, open := <-eventChan
		if ! open {
			break eventsLoop
		}
		if ! auth.AccessTokenGranted {
			log.Println("Cannot accept events until access token is granted")
		}
		err := handleEvent(event)
		if err != nil {
			log.Println(err)
		}
	}
}

func handleEvent(event Event) error {
	switch event.Type {
	case "transaction.created":
		byteArr, _ := json.Marshal(event)
		fmt.Println("handle transaction notification: ", string(byteArr))
		
		pennies := event.Data.Amount % 100
		roundUpAmount := 100 - pennies
		log.Println("Transfer", roundUpAmount)
		return nil
	default:
		return errors.New("Event `" + event.Type + "` not handled")
	}
}
