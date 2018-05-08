package monzoroundup

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
