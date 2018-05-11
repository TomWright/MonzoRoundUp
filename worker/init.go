package worker

import (
	"github.com/tomwright/monzoroundup/user"
	"github.com/tomwright/monzoroundup/monzo"
	"github.com/tomwright/monzoroundup/token"
)

func init() {
	HandleEvent = &handleEvent{
		C: make(chan interface{}, 5),
	}
	HandleTransactionCreatedEvent = &handleTransactionCreatedEvent{
		C: make(chan monzo.TransactionCreatedEvent, 5),
	}

	InitUser = &initUser{
		C: make(chan user.User, 5),
	}
}

func Listen(tokenModel token.Model, userModel user.Model) {
	go InitUser.Work(tokenModel)
	go HandleEvent.Work()
	go HandleTransactionCreatedEvent.Work(userModel, tokenModel)
}
