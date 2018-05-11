package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
	"fmt"
	"os"
)

func init() {
	goose.AddMigration(Up00002, Down00002)
}

func Up00002(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`
		CREATE TABLE %s.users (
		id   INT AUTO_INCREMENT,
		username VARCHAR(100),
		oauthId VARCHAR(255),
		oauthSecret VARCHAR(255),
		PRIMARY KEY (id),
		UNIQUE KEY (username)
		);
	`, os.Getenv("MYSQL_DB")))
	if err != nil {
		return err
	}
	return nil
}

func Down00002(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`
		DROP TABLE %s.users;
	`, os.Getenv("MYSQL_DB")))
	if err != nil {
		return err
	}
	return nil
}
