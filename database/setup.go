package database

import (
	"database/sql"
	"net/url"

	"github.com/jmoiron/sqlx"
)

// Connection is an interface that defines the methods that
// the database client must implement. This allows us to
// mock the database client in our tests and swap it out
// with the real one in our main function.
type Connection interface {
	GetUserData(string) (string, error)
	UpsertUsername(string, string) error
	DeleteUsername(string) error
	CreateGame(string) (int64, error)
	AddUserToGame(string, int64) error
	GetGameData(int64) (Game, error)
	StopGame(int64) error
	CheckUserValid(string, string) (bool, error)
}

// Client is the real database client that satisfies the
// Connection interface. It embeds a sqlx.DB struct, which
// contains all the methods we need to interact with the
// database.
type Client struct {
	db *sqlx.DB
}

func (db Client) GetSqlDB() *sql.DB {
	return db.db.DB
}

// Setup is a function that returns a new database client.
// It also creates the users table if it doesn't exist.
// The current implementation uses a PostgreSQL database, that is
// running in a Docker container.
func Setup() *Client {
	// connect to db
	params := url.Values{}
	params.Set("sslmode", "disable")

	connectionString := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost:5431",
		Path:     "postgres",
		RawQuery: params.Encode(),
	}
	db, err := sqlx.Connect("postgres", connectionString.String())
	if err != nil {
		panic(err)
	}

	tableQuery := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS games (
		id SERIAL PRIMARY KEY,
		host VARCHAR(255) NOT NULL,
		players VARCHAR(255)[5],
		answer VARCHAR(255) NOT NULL,
		questions VARCHAR(255)[],
		guesses VARCHAR(255)[],
		start_time TIMESTAMP NOT NULL DEFAULT NOW(),
		end_time TIMESTAMP,
		ended BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS questions (
		id SERIAL PRIMARY KEY,
		question VARCHAR(255) NOT NULL,
		user_id INTEGER NOT NULL references users(id),
		game_id INTEGER NOT NULL references games(id)
	);

	CREATE TABLE IF NOT EXISTS guesses (
		id SERIAL PRIMARY KEY,
		guess VARCHAR(255) NOT NULL,
		user_id INTEGER NOT NULL references users(id),
		game_id INTEGER NOT NULL references games(id),
		correct BOOLEAN DEFAULT FALSE
	);

`
	db.MustExec(tableQuery)

	return &Client{
		db: db,
	}
}
