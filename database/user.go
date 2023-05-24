package database

import "fmt"

// GetUserData returns the username from the database.
// TODO: This is a placeholder function for now it will return
// all the user data from the database.
func (c *Client) GetUserData(username string) (string, error) {
	var userData string
	err := c.db.QueryRow("SELECT username FROM users WHERE username = $1;", username).Scan(&userData)
	if err != nil {
		return "", fmt.Errorf("failed to get username %s from db: %w", username, err)
	}
	return userData, nil
}

// UpsertUsername inserts a new username into the database.
func (c *Client) UpsertUsername(username, password string) error {
	registerUser := `INSERT INTO users (username) VALUES ($1)
										ON CONFLICT (username) DO NOTHING;`
	_, err := c.db.Exec(registerUser, username)
	if err != nil {
		return fmt.Errorf("failed to register user %s: %w", username, err)
	}
	return nil
}

// DeleteUsername deletes a username from the database.
func (c *Client) DeleteUsername(username string) error {
	deleteUser := `DELETE FROM users WHERE username = $1;`
	_, err := c.db.Exec(deleteUser, username)
	if err != nil {
		return fmt.Errorf("failed to delete user %s: %w", username, err)
	}
	return nil
}

// CheckUserValid checks if the user is valid my checking that it
// exists in the database and that the password matches.
func (c *Client) CheckUserValid(username, password string) (bool, error) {
	query := `SELECT username, password FROM users WHERE username = $1 ;`
	var user, pass string
	err := c.db.QueryRow(query, username).Scan(&user, &pass)
	if err != nil {
		return false, fmt.Errorf("failed to get user %s: %w", username, err)
	}

	// having a password check at this level makes using sqlc pointless
	if password != pass {
		return false, nil
	}
	return true, nil
}
