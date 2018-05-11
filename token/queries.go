package token

import (
	"database/sql"
	"github.com/tomwright/monzoroundup/monzo"
	"errors"
	"fmt"
)

func NewSqlModel(db *sql.DB) *SqlModel {
	m := new(SqlModel)
	m.db = db
	return m
}

type SqlModel struct {
	db *sql.DB
}

func (m *SqlModel) FetchMostRecentForAllUsers() (map[int64]*monzo.Token, error) {
	stmt, err := m.db.Prepare("SELECT " +
		"id, accessToken, clientId, refreshToken, tokenType, userId, expiresAt " +
		"FROM userTokens GROUP BY userId ORDER BY expiresAt DESC;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	res := make(map[int64]*monzo.Token)

	for rows.Next() {
		u := monzo.Token{}
		err = rows.Scan(&u.ID, &u.AccessToken, &u.ClientID, &u.RefreshToken, &u.TokenType, &u.UserID, &u.ExpiresAt)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("could not scan token row: %s", err))
		}
		res[u.UserID] = &u
	}

	return res, nil
}

func (m *SqlModel) FetchByID(id int64) (*monzo.Token, error) {
	stmt, err := m.db.Prepare("SELECT " +
		"id, accessToken, clientId, refreshToken, tokenType, userId, expiresAt " +
		"FROM userTokens WHERE id = ? LIMIT 1;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	u := monzo.Token{}

	err = row.Scan(&u.ID, &u.AccessToken, &u.ClientID, &u.RefreshToken, &u.TokenType, &u.UserID, &u.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m *SqlModel) FetchByUserID(userID string) (*monzo.Token, error) {
	stmt, err := m.db.Prepare("SELECT " +
		"id, accessToken, clientId, refreshToken, tokenType, userId, expiresAt " +
		"FROM userTokens WHERE userId = ? ORDER BY expiresAt DESC LIMIT 1;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID)
	u := monzo.Token{}

	err = row.Scan(&u.ID, &u.AccessToken, &u.ClientID, &u.RefreshToken, &u.TokenType, &u.UserID, &u.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m *SqlModel) Insert(userID int64, token *monzo.Token) (*monzo.Token, error) {
	stmt, err := m.db.Prepare("" +
		"INSERT INTO userTokens (accessToken, clientId, refreshToken, tokenType, userId, expiresAt)" +
		"VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(token.AccessToken, token.ClientID, token.RefreshToken, token.TokenType, userID, token.ExpiresAt)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return m.FetchByID(lastId)
}
