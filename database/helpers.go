// These are helper functions for testing
package database

import "database/sql"

func (db *Client) TestCleanup() {
	db.db.Exec("DELETE TABLE IF EXISTS guesses")
	db.db.Exec("DELETE TABLE IF EXISTS questions")
	db.db.Exec("DELETE TABLE IF EXISTS games")
	db.db.Exec("DELETE TABLE IF EXISTS players")
	db.db.Close()
}

func (db *Client) TestSetup() error {
	_, err := db.db.Exec(`
INSERT INTO players (username, password) VALUES ('captainnobody1', 'password');
`)
	return err
}

func (db *Client) TestSetupGame() *sql.Row {
	rows := db.db.QueryRow(`select count(*) from games`)
	return rows
}
