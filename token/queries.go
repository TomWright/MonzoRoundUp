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
		"t1.id, t1.accessToken, t1.clientId, t1.refreshToken, t1.tokenType, t1.monzoUserId, t1.userId, t1.expiresAt " +
		"FROM userTokens t1 " +
		"LEFT JOIN userTokens t2 ON t1.userId = t2.userId " +
		"WHERE t2.userId IS NULL ORDER BY t1.expiresAt DESC;")
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
		var userId int64
		err = rows.Scan(&u.ID, &u.AccessToken, &u.ClientID, &u.RefreshToken, &u.TokenType, &u.UserID, &userId, &u.ExpiresAt)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("could not scan token row: %s", err))
		}
		res[userId] = &u
	}

	return res, nil
}

func (m *SqlModel) FetchByID(id int64) (*monzo.Token, int64, error) {
	stmt, err := m.db.Prepare("SELECT " +
		"id, accessToken, clientId, refreshToken, tokenType, monzoUserId, userId, expiresAt " +
		"FROM userTokens WHERE id = ? LIMIT 1;")
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	u := monzo.Token{}

	var userId int64

	err = row.Scan(&u.ID, &u.AccessToken, &u.ClientID, &u.RefreshToken, &u.TokenType, &u.UserID, &userId, &u.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, userId, nil
	}
	if err != nil {
		return nil, userId, err
	}

	return &u, userId, nil
}

func (m *SqlModel) FetchByUserID(userID string) (*monzo.Token, error) {
	stmt, err := m.db.Prepare("SELECT " +
		"id, accessToken, clientId, refreshToken, tokenType, monzoUserId, expiresAt " +
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

	newToken, _, err := m.FetchByID(lastId)
	return newToken, err
}
