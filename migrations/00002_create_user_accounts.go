package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up00002, Down00002)
}

func Up00002(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE userAccounts (
		id   VARCHAR(100),
		userId INT,
		PRIMARY KEY (id),
		INDEX (userId)
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func Down00002(tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE userAccounts;
	`)
	if err != nil {
		return err
	}
	return nil
}
