package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up00003, Down00003)
}

func Up00003(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE userTokens (
		id INT AUTO_INCREMENT,
		accessToken   VARCHAR(255),
		clientId   VARCHAR(255),
		refreshToken   VARCHAR(255),
		tokenType   ENUM('BEARER'),
		userId   INT,
		expiresAt   DATETIME,
		PRIMARY KEY (id),
		INDEX(userId)
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func Down00003(tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE userTokens;
	`)
	if err != nil {
		return err
	}
	return nil
}
