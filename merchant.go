package monzoroundup

type merchant struct {
	Address  merchantAddress `json:"address"`
	Created  string          `json:"created"`
	GroupID  string          `json:"group_id"`
	ID       string          `json:"id"`
	Logo     string          `json:"logo"`
	Emoji    string          `json:"emoji"`
	Name     string          `json:"name"`
	Category string          `json:"category"`
}

type merchantAddress struct {
	Address   string  `json:"address"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Postcode  string  `json:"postcode"`
	Region    string  `json:"region"`
}
