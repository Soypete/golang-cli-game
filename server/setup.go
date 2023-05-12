// server is a module that contains the server and all of its routes.
package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soypete/golang-cli-game/database"
)

// State is the global state of the server.
type State struct {
	db      database.Connection
	Router  *chi.Mux
	BaseURL string
	Port    string
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
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// add pprof
	r.Mount("/debug", middleware.Profiler())
	// add prometheus metrics
	r.Mount("/metrics", promhttp.Handler())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to game server"))
	})

	s := &State{
		db:      db,
		Router:  r,
		BaseURL: "http://localhost:3000", // TODO: this should be a config
		Port:    ":3000",
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

	r.Route("/game", func(r chi.Router) {
		// /start add you to the host role
		r.Get("/start", s.startGame) // GET /game/start
		// // subroutes for game
		r.Route("/{gameID}", func(r chi.Router) {
			r.Get("/join", s.joinGame)       // GET /game/123/join?username=
			r.Get("/status", s.getGameState) // GET /game/123/status
			r.Get("/play", s.playGame)       // GET /game/123/play?username=&answer=...
			// r.Get("/turn", s.TakeTurn) // GET /game/123/turn?username=&action=question/answer&question=...&answer=...
			// 	// only the host can get the summary
			// 	r.Get("/summary", s.getSummary) // GET /game/123/summary
			// 	// only the host can stop the game
			r.Get("/stop", s.stopGame) // GET /game/123/stop
		})
		// /abandoned returns all games that have been abandoned without being finished
		// r.Get("/abandoned", s.getAbandonedGames) // GET /game/abandoned
	})

	return s
}
