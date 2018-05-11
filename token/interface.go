package token

import "github.com/tomwright/monzoroundup/monzo"

type Model interface {
	FetchByID(id int64) (*monzo.Token, error)
	FetchByUserID(userID string) (*monzo.Token, error)
	Insert(userID int64, token *monzo.Token) (*monzo.Token, error)
	FetchMostRecentForAllUsers() (map[int64]*monzo.Token, error)
}
