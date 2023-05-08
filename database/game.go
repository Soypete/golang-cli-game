package database

import (
	"fmt"
	"time"
)

// Game represents a game in the database.
type Game struct {
	GameID        int64    `db:"id"`
	Host          string   `db:"host"`
	Players       []string `db:"players"` // only 5 players allowed per game
	Answer        string   `db:"answer"`  // TODO: answer validation should account for capitalization and spelling errors.. maybe use a Levenshtein distance algorithm
	QuestionCount int64
	Questions     []Question `db:"questions"`
	Guesses       []Guess    `db:"guesses"`
	StartTime     time.Time  `db:"start_time"`
	EndTime       time.Time  `db:"end_time, omitempty""`
	Ended         bool       `db:"ended"`
}

// Question represents a question in the database.
type Question struct {
	QuestionID   string
	QuestionText string
	UserID       string // the user who asked the question
	GameID       string // the game the question is associated with
}

// Guess represents a guess in the database. This is a user's guess of the answer.
type Guess struct {
	GuessID   string
	GuessText string
	UserID    string // the user who made the guess
	GameID    string // the game the guess is associated with
	Correct   bool   // whether the guess is correct or not
}

// CreateGame starts a new game for the user with the given username
// as the host. A new game is created and the user is added to the game.
// The game id is returned, or an error if one occurs.
func (c *Client) CreateGame(username string) (int64, error) {
	query := `INSERT INTO games (host, players)
					VALUES ($1, $1)`
	results, err := c.db.Exec(query, username)
	if err != nil {
		return 0, fmt.Errorf("unable to create game instance: %w", err)
	}
	gameID, err := results.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("unable to get game id: %w", err)
	}
	return gameID, nil
}

// AddUserToGame adds the user with the given username to the game with the
// given game id. An error is returned if one occurs.
func (c *Client) AddUserToGame(username string, gameID int64) error {
	// TODO: check if game is full or started. If so, return an error. We shouldn't add a user if answer guessing has already started.
	query := `UPDATE games
					SET players = array_append(players, $1)
					WHERE game_id = $2`
	_, err := c.db.Exec(query, username, gameID)
	if err != nil {
		return fmt.Errorf("unable to add user to game: %w", err)
	}
	return nil
}

// GetGameData returns the game info for the game with the given game id.
func (c *Client) GetGameData(gameID int64) (Game, error) {
	query := `SELECT * FROM games WHERE game_id = $1`
	var game Game
	err := c.db.QueryRowx(query, gameID).StructScan(&game)
	if err != nil {
		return Game{}, fmt.Errorf("unable to get game info: %w", err)
	}
	return game, nil
}

func (c *Client) StopGame(gameID int64) error {
	query := `UPDATE games
					SET ended = true, end_time = NOW()
					WHERE game_id = $1`
	_, err := c.db.Exec(query, gameID)
	if err != nil {
		return fmt.Errorf("unable to stop game: %w", err)
	}
	return nil
}
