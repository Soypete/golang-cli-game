package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

type passDB struct{}

func (db *passDB) GetUserData(username string) (string, error) {
	return "captainnobody1", nil
}

func (db *passDB) UpsertUsername(username string) error {
	return nil
}

func (db *passDB) DeleteUsername(username string) error {
	return nil
}
func (db *passDB) CreateGame(username string) (int64, error) {
	return 1234, nil
}
func (db *passDB) AddUserToGame(username string, gameID int64) error {
	return nil
}

type failDB struct{}

func (db *failDB) GetUserData(username string) (string, error) {
	return "", fmt.Errorf("failed to get username %s from db", username)
}

func (db *failDB) UpsertUsername(username string) error {
	return fmt.Errorf("failed to update username %s from db", username)
}

func (db *failDB) DeleteUsername(username string) error {
	return fmt.Errorf("failed to delete username %s from db", username)
}
func (db *failDB) CreateGame(username string) (int64, error) {
	return 0, fmt.Errorf("failed to start game for username %s from db", username)
}
func (db *failDB) AddUserToGame(username string, gameID int64) error {
	return fmt.Errorf("failed to add user %s to game %d from db", username, gameID)
}

func setupTestRouter(s State, t *testing.T) *chi.Mux {
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
	r.Route("/game", func(r chi.Router) {
		// /start add you to the host role
		r.Get("/start", s.startGame) // GET /game/start
		// // subroutes for game
		r.Route("/{gameID}", func(r chi.Router) {
			r.Get("/join", s.joinGame) // GET /game/123/join
			// 	r.Get("/leave", s.leaveGame) // GET /game/123/leave
			// 	// starting = no answer submitted, in progess = asking questions, finished = guest guessed or game stopped
			// 	r.Get("/status", s.getGameStatus) // GET /game/123/status
			// 	// only the host can get the summary
			// 	r.Get("/summary", s.getSummary) // GET /game/123/summary
			// 	// only the host can stop the game
			// 	r.Get("/stop", s.stopGame) // GET /game/123/stop
		})
		// /abandoned returns all games that have been abandoned without being finished
		// r.Get("/abandoned", s.getAbandonedGames) // GET /game/abandoned
	})
	return r
}

func TestUserEndpoints(t *testing.T) {
	t.Run("get username: Pass", testPassGetUserName)
	t.Run("get username: No Username", testFailGetUsernameEmpty)
	t.Run("get username:Fail", testFailGetUsernameDB)
	t.Run("update username: Pass", testPassUpdateUser)
	t.Run("update username: No Username", testFailUpdateUsernameEmpty)
	t.Run("update username:Fail", testFailUpdateUsernameDB)
	t.Run("delete username: Pass", testPassDeleteUser)
	t.Run("delete username:Fail", testFailDeleteUsernameDB)
	t.Run("create game: Pass", testPassStartGame)
	t.Run("create game: Fail", testFailStartGame)
	t.Run("join game: Pass", testPassJoinGame)
	t.Run("join game: No Username", testFailJoinNoGameID)
	t.Run("join game:Fail", testFailJoinGameDB)
}
func testPassGetUserName(t *testing.T) {
	sPass := State{
		db: new(passDB),
	}
	sPass.Router = setupTestRouter(sPass, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/register/captainnobody1/get", nil)
	sPass.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// test no username
func testFailGetUsernameEmpty(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	reqNoUser := httptest.NewRequest("GET", "/register//get", nil)
	sFail.Router.ServeHTTP(w, reqNoUser)
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// test db error
func testFailGetUsernameDB(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/register/captainnobody1/get", nil)
	sFail.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func testPassUpdateUser(t *testing.T) {
	sPass := State{
		db: new(passDB),
	}
	sPass.Router = setupTestRouter(sPass, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/register/captainnobody1/update", nil)
	sPass.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// test no username
func testFailUpdateUsernameEmpty(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	reqNoUser := httptest.NewRequest("GET", "/register//update", nil)
	sFail.Router.ServeHTTP(w, reqNoUser)
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// test db error
func testFailUpdateUsernameDB(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/register/captainnobody1/update", nil)
	sFail.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func testPassDeleteUser(t *testing.T) {
	sPass := State{
		db: new(passDB),
	}
	sPass.Router = setupTestRouter(sPass, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/register/captainnobody1/delete", nil)
	sPass.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// test db error
func testFailDeleteUsernameDB(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/register/captainnobody1/delete", nil)
	sFail.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func testPassStartGame(t *testing.T) {
	sPass := State{
		db: new(passDB),
	}
	sPass.Router = setupTestRouter(sPass, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/game/start", nil)
	sPass.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func testFailStartGame(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/game/start", nil)
	sFail.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}
func testPassJoinGame(t *testing.T) {
	sPass := State{
		db: new(passDB),
	}
	sPass.Router = setupTestRouter(sPass, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/game/1234/join", nil)
	sPass.Router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// test no GameID
func testFailJoinNoGameID(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	reqNoUser := httptest.NewRequest("GET", "/game//join", nil)
	sFail.Router.ServeHTTP(w, reqNoUser)
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func testFailJoinGameDB(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/game/1232/join", nil)
	sFail.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}
