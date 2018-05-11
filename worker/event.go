package worker

import (
	"log"
	"github.com/tomwright/monzoroundup/user"
	"time"
	"fmt"
	"github.com/tomwright/monzoroundup/monzo"
	"database/sql"
	"github.com/tomwright/monzoroundup/token"
)

var (
	HandleEvent                   *handleEvent
	HandleTransactionCreatedEvent *handleTransactionCreatedEvent
)

type handleEvent struct {
	C chan interface{}
}

func (i *handleEvent) Work() {
	events := i.C
	for {
		e, open := <-events
		if ! open {
			return
		}

		switch e.(type) {
		case monzo.TransactionCreatedEvent:
			HandleTransactionCreatedEvent.C <- e.(monzo.TransactionCreatedEvent)
		default:
			log.Printf("unhandled event type: %t\n", e)
		}
	}
}

type handleTransactionCreatedEvent struct {
	C chan monzo.TransactionCreatedEvent
}

func (i *handleTransactionCreatedEvent) Work(userModel user.Model, tokenModel token.Model) {
	events := i.C
	for {
		e, open := <-events
		if ! open {
			return
		}

		log.Printf("handle transaction created event: %v\n", e)

		u, err := userModel.FetchByAccountID(e.AccountID)
		if err == sql.ErrNoRows {
			log.Printf("account id `%s` is not stored against any user\n", e.AccountID)
		}
		if err != nil {
			log.Println("could not fetch user by account id: ", err)
			continue
		}
		t, err := tokenModel.FetchByUserID(u.ID)
		if err != nil {
			log.Printf("could not fetch token by user id: %s\n", err)
			continue
		}

		if t == nil {
			log.Printf("user `%s` has no token\n", u.UserName)
			continue
		}

		pot := i.getPot(t)

		amount := 100 - (e.Amount % 100)

		newPot, err := monzo.DepositToPotFromAccount(t, pot.ID, e.AccountID, amount, fmt.Sprintf("%d", time.Now().Unix()))
		if err != nil {
			log.Println("could not run deposit: ", err)
			continue
		}

		log.Printf("Pot value before/after: %d/%d", pot.Balance, newPot.Balance)
	}
}

func (i *handleTransactionCreatedEvent) getPot(token *monzo.Token) *monzo.Pot {
	pots, err := monzo.ListPots(token)
	if err != nil {
		log.Printf("could not list monzo pots: %s", err)
		return nil
	}

	for _, pot := range pots {
		if pot.Name == "Round Up" {
			return &pot
		}
	}

	return nil
}
