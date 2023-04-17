// server is a module that contains the server and all of its routes.
package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/soypete/golang-cli-game/database"
)

// State is the global state of the server.
type State struct {
	db     database.Connection
	Router *chi.Mux
}

// NewState creates a new server state.
func NewState() *State {
	// create table if not exists
	db := database.Setup()
	// setup chi server
	//
	// curl http://localhost:3000
	r := chi.NewRouter()
	// add logger for all requests
	r.Use(middleware.Logger)

	// TODO(soypete): investigate middleware to see if we should add any for intro
	// topics or wait until later
	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	// r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to game server"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// TODO(soypete): investigate middleware to see if we should add any for intro
	// topics or wait until later
	//r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	//r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to game server"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	s := &State{
		db:     db,
		Router: r,
	}

	// setup routes
	r.Route("/register", func(r chi.Router) {
		// subroutes for register
		r.Route("/{username}", func(r chi.Router) {
			r.Get("/get", s.getUsername)          // GET /register/123/get
			r.Get("/update", s.updateUsername)    // PUT /register/123/update // TODO: this is a get because I am not providing a body
			r.Delete("/delete", s.deleteUsername) // DELETE /register/123/delete //TODO: should this have /delete in the path?
		})
	})

	// TODO: add more routes
	// - start game - generate a game id and join code
	// - join game - use the join code to join a game
	// - leave game - leave a game
	// - get game history - user's win loss history
	// - get game status - get current game status. if thre is a game in progress, return the game id
	// - get abandoned games - get a list of games that have been abandoned
	//  - stop an abandoned game - stop a game that has been abandoned

	return s
}
