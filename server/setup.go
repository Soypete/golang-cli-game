// server is a module that contains the server and all of its routes.
package server

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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

var (
	// expvar counters for the number of status codes returned
	counter200Code expvar.Int
	counter400Code expvar.Int
	counter500Code expvar.Int
)

// NewState creates a new server state.
func NewState() (*State, error) {
	// create table if not exists
	db, err := database.Setup()
	if err != nil {
		return nil, err
	}
	// setup chi server
	//
	// curl http://localhost:3000
	r := chi.NewRouter()

	// add prebuild middleware for all requests
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// add expvars at /debug/vars
	expvar.Publish("counter200Code", &counter200Code)
	expvar.Publish("counter400Code", &counter400Code)
	expvar.Publish("counter500Code", &counter500Code)
	// add pprof at /debug/
	r.Mount("/debug", middleware.Profiler())

	// add expvar and db metrics to prometheus
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewBuildInfoCollector())
	reg.MustRegister(collectors.NewDBStatsCollector(db.GetSqlDB(), "postgres"))
	// TODO: figure out how to add types
	reg.MustRegister(collectors.NewExpvarCollector(
		map[string]*prometheus.Desc{
			"counter200Code": prometheus.NewDesc("expvar_200Status", "number of status 200 api calls", nil, nil),
			"counter400Code": prometheus.NewDesc("expvar_400status", "number of status 400 api calls", nil, nil),
			"counter500Code": prometheus.NewDesc("expvar_500status", "number of status 500 api calls", nil, nil),
		},
	))
	// add prometheus endpoint at /metrics. The above collectors will be show
	// in the reverse order the are registered.
	r.Mount("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to game server"))
	})

	s := &State{
		db:      db,
		Router:  r,
		BaseURL: "http://localhost:3000", // TODO: this should be a config
		Port:    ":3000",
	}

	// TODO: Change register routes.
	// /register?username=...&password... (no auth)
	// user/{username}/get (with auth)
	// user/{username}/update?password=... (with auth)
	// user/{username}/delete (with auth)

	// setup routes
	r.Route("/register", func(r chi.Router) {
		// subroutes for register
		r.Route("/{username}", func(r chi.Router) {
			r.Get("/get", s.getUsername)          // GET /register/123/get
			r.Get("/update", s.updateUsername)    // PUT /register/123/update?password=..
			r.Delete("/delete", s.deleteUsername) // DELETE /register/123/delete //TODO: should this have /delete in the path?
		})
	})

	// add middleware to /game routes
	r.With(s.middlewareHandler).Route("/game", func(r chi.Router) {
		// /start add you to the host role
		r.Get("/start", s.startGame) // GET /game/start
		// // subroutes for game
		r.Route("/{gameID}", func(r chi.Router) {
			r.Get("/join", s.joinGame)       // GET /game/123/join?
			r.Get("/status", s.getGameState) // GET /game/123/status
			r.Get("/play", s.playGame)       // GET /game/123/play?&answer=...
			// r.Get("/turn", s.TakeTurn) // GET /game/123/turn?action=question/answer&question=...&answer=...
			// 	// only the host can get the summary
			// 	r.Get("/summary", s.getSummary) // GET /game/123/summary
			// 	// only the host can stop the game
			r.Get("/stop", s.stopGame) // GET /game/123/stop
		})
		// /abandoned returns all games that have been abandoned without being finished
		// r.Get("/abandoned", s.getAbandonedGames) // GET /game/abandoned

	})

	return s, nil
}

func (s *State) middlewareHandler(h http.Handler) func(next http.Handler) http.Handler {
	s.authMiddleware(h)
}
