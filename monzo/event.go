package monzo

import (
	"encoding/json"
	"errors"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (e Event) ToType() (interface{}, error) {
	dataJSON, err := json.Marshal(e.Data)
	if err != nil {
		return nil, err
	}
	switch e.Type {
	case "transaction.created":
		r := TransactionCreatedEvent{}
		json.Unmarshal(dataJSON, &r)
	}

	return nil, errors.New("unhandled event type: " + e.Type)
}
