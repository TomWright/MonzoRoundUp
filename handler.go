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

		pot, err := getRoundUpPot()
		if err != nil {
			return errors.New("could not fetch round up pot: " + err.Error())
		}

		deposit := Deposit{
			PotID:     pot.ID,
			AccountID: event.Data.AccountID,
			Amount:    roundUpAmount,
			DeDupeID:  fmt.Sprintf("%s_%s_%s_%d", pot.ID, event.Data.AccountID, event.Data.ID, roundUpAmount),
		}

		newPot, err := depositToPotFromAccount(deposit)
		if err != nil {
			return errors.New("could not round up pot: " + err.Error())
		}

		log.Printf("Pot `%s` rounded up from account `%s`\n"+
			"\tTransaction amount: %d\n"+
			"\tRound up amount: %d\n"+
			"\tOld pot value: %d\n"+
			"\tNew pot value: %d\n", pot.ID, event.Data.AccountID, event.Data.Amount, roundUpAmount, pot.Balance, newPot.Balance)
		return nil
	default:
		return errors.New("Event `" + event.Type + "` not handled")
	}
}
