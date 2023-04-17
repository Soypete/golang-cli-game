// This cli game will be hosted on a server and will be played by multiple users
// we don't know what the game will do yet, but we will start with a simple /register
// and database to store usernames

// TODO: add to readme with DB example:
// - https://github.com/jackc/pgx
// - https://github.com/lib/pq
// - https://github.com/jmoiron/sqlx
package main

import (
	"net/http"

	_ "github.com/lib/pq"
	"github.com/soypete/golang-cli-game/server"
)

func main() {

	gameState := server.NewState()

	// setup chi server
	// curl http://localhost:3000
	http.ListenAndServe(":3000", gameState.Router)
}
