package monzo

type TransactionCreatedEvent struct {
	AccountID   string   `json:"account_id"`
	Amount      int64    `json:"amount"`
	Created     string   `json:"created"`
	Currency    string   `json:"currency"`
	Description string   `json:"description"`
	ID          string   `json:"id"`
	Category    string   `json:"category"`
	IsLoad      bool     `json:"is_load"`
	Settled     bool     `json:"settled"`
	Merchant    Merchant `json:"merchant"`
}
