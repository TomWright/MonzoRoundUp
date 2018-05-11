package user

import (
	"database/sql"
)

func NewSqlModel(db *sql.DB) *SqlModel {
	m := new(SqlModel)
	m.db = db
	return m
}

type SqlModel struct {
	db *sql.DB
}

func (m *SqlModel) FetchByID(id int64) (*User, error) {
	stmt, err := m.db.Prepare("SELECT id, username, oauthId, oauthSecret FROM users WHERE id = ? LIMIT 1;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	u := User{}

	err = row.Scan(&u.ID, &u.UserName, &u.OAuthID, &u.OAuthSecret)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m *SqlModel) FetchByUserName(userName string) (*User, error) {
	stmt, err := m.db.Prepare("SELECT id, username, oauthId, oauthSecret FROM users WHERE username = ? LIMIT 1;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(userName)
	u := User{}

	err = row.Scan(&u.ID, &u.UserName, &u.OAuthID, &u.OAuthSecret)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m *SqlModel) FetchByAccountID(accountID string) (*User, error) {
	stmt, err := m.db.Prepare("" +
		"SELECT users.id, users.username, users.oauthId, users.oauthSecret" +
		"FROM users" +
		"JOIN userAccounts ON userAccounts.userId = users.id" +
		"WHERE userAccounts.id = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(accountID)
	u := User{}

	err = row.Scan(&u.ID, &u.UserName, &u.OAuthID, &u.OAuthSecret)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}
