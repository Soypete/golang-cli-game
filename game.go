// This cli game will be hosted on a server and will be played by multiple users
// we don't know what the game will do yet, but we will start with a simple /register
// and database to store usernames

package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/soypete/golang-cli-game/server"
)

func main() {

	gameState, err := server.NewState()
	if err != nil {
		log.Fatal(err)
	}

	// setup chi server
	// curl http://localhost:3000
	http.ListenAndServe(gameState.Port, gameState.Router)
}
