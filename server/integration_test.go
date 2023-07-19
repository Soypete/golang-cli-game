package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/soypete/golang-cli-game/database"
)

func setupIntegrationRouter(s State, t *testing.T) *chi.Mux {
	r := chi.NewRouter()

	// setup routes
	r.Route("/register", func(r chi.Router) {
		// subroutes for register
		r.Route("/{username}", func(r chi.Router) {
			r.Get("/get", s.getUsername)          // GET /register/123/get
			r.Get("/update", s.updateUsername)    // PUT /register/123/update // TODO: this is a get because I am not providing a body
			r.Delete("/delete", s.deleteUsername) // DELETE /register/123/delete
		})
	})
	r.With(s.middlewareHandler).Route("/game", func(r chi.Router) {
		r.Get("/start", s.startGame) // GET /game/start
		r.Route("/{gameID}", func(r chi.Router) {
			r.Get("/join", s.joinGame)       // GET /game/123/join
			r.Get("/status", s.getGameState) // GET /game/123/status
			r.Get("/stop", s.stopGame)       // GET /game/123/stop
		})
	})
	return r
}

func TestAuthFlow(t *testing.T) {
	if os.Getenv("CI") != "true" {
		t.Skip("Skipping integration tests")
	}
	db, err := database.Setup()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
	}
	defer t.Cleanup(func() {
		db.TestCleanup()
	})
	err = db.TestSetup()
	if err != nil {
		t.Errorf("Error adding data to database: %v", err)
	}
	game := State{
		db:      db,
		BaseURL: "http://localhost:3000", // TODO: this should be a config
		Port:    ":3000",
	}
	game.Router = setupIntegrationRouter(game, t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/game/start", nil)
	req.Header.Set("Authorization", getAuthHeader())
	game.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	row := db.TestSetupGame()
	if row.Err() != nil {
		t.Errorf("Error adding data to database: %v", err)
	}
	var count int
	err = row.Scan(&count)
	if err != nil {
		t.Errorf("Error scanning rows: %v", err)
	}
	// check that one game is created
	if count != 1 {
		t.Errorf("Error: expected at least one row in games table")
	}
}
