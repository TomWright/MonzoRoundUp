package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up00001, Down00001)
}

func Up00001(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE users (
		id   INT AUTO_INCREMENT,
		username VARCHAR(100),
		oauthId VARCHAR(255),
		oauthSecret VARCHAR(255),
		PRIMARY KEY (id),
		UNIQUE KEY (username)
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func Down00001(tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE users;
	`)
	if err != nil {
		return err
	}
	return nil
}
