package user

func (m *SqlModel) Create(userName string, oauthID string, oauthSecret string) (*User, error) {
	stmt, err := m.db.Prepare("" +
		"INSERT INTO users (username, oauthId, oauthSecret)" +
		"VALUES (?, ?, ?);")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(userName, oauthID, oauthSecret)
	if err != nil {
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return m.FetchByID(lastId)
}
