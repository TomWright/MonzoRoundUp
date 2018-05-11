package user

type Model interface {
	FetchByUserName(userName string) (*User, error)
	FetchByID(id int64) (*User, error)
	FetchByAccountID(accountID string) (*User, error)
	Create(userName string, oauthID string, oauthSecret string) (*User, error)
}
