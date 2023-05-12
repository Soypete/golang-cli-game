package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/soypete/golang-cli-game/database"
)

type passDB struct{}

func getAuthHeader() string {
	encoding := base64.StdEncoding
	ed := encoding.EncodeToString([]byte("captainnobody1:password"))
	return fmt.Sprintf("Basic %s", ed)
}

func (db *passDB) GetUserData(username string) (string, error) {
	return "captainnobody1", nil
}

func (db *passDB) UpsertUsername(username, password string) error {
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
func (db *passDB) GetGameData(gameID int64) (database.Game, error) {
	return database.Game{
		GameID: 321,
	}, nil
}
func (db *passDB) StopGame(gameID int64) error {
	return nil
}
func (db *passDB) CheckUserValid(string, string) (bool, error) {
	return true, nil
}

type failDB struct{}

func (db *failDB) GetUserData(username string) (string, error) {
	return "", fmt.Errorf("failed to get username %s from db", username)
}
func (db *failDB) UpsertUsername(username, password string) error {
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
func (db *failDB) GetGameData(gameID int64) (database.Game, error) {
	return database.Game{}, fmt.Errorf("failed to get game %d from db", gameID)
}
func (db *failDB) StopGame(gameID int64) error {
	return fmt.Errorf("failed to stop game %d from db", gameID)
}
func (db *failDB) CheckUserValid(string, string) (bool, error) {
	return false, nil
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
			r.Get("/status", s.getGameState) // GET /game/123/status
			// 	// only the host can get thummary
			// 	r.Get("/summary", s.getSummary) // GET /game/123/summary
			// 	// only the host can stop the game
			r.Get("/stop", s.stopGame) // GET /game/123/stop
		})
		// /abandoned returns all games that have been abandoned without being finished
		// r.Get("/abandoned", s.getAbandonedGames) // GET /game/abandoned
	})
	return r
}

func TestUserEndpoints(t *testing.T) {
	t.Run("get username: Pass", testPassGetUserName)
	t.Run("get username: No Username", testPassGetUsernameEmpty)
	// t.Run("get username: No Header", testFailGetUsernameNoHeader)
	t.Run("get username:Fail", testFailGetUsernameDB)
	t.Run("update username: Pass", testPassUpdateUser)
	t.Run("update username: No Username", testFailUpdateUsernameEmpty)
	t.Run("update username:Fail", testFailUpdateUsernameDB)
	t.Run("delete username: Pass", testPassDeleteUser)
	// t.Run("delete username:Fail No Header", testFailDeleteUsernameNoHeader)
	t.Run("delete username:Fail", testFailDeleteUsernameDB)
	t.Run("create game: Pass", testPassStartGame)
	t.Run("create game: Fail", testFailStartGame)
	t.Run("create game: Fail no header", testFailStartGameNoHeader)
	t.Run("join game: Pass", testPassJoinGame)
	t.Run("join game: No gameID", testFailJoinNoGameID)
	t.Run("join game: No Header", testFailJoinGameNoHeader)
	t.Run("join game:Fail", testFailJoinGameDB)
}
func testPassGetUserName(t *testing.T) {
	sPass := State{
		db: new(passDB),
	}
	sPass.Router = setupTestRouter(sPass, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/register/captainnobody1/get", nil)
	req.Header.Set("Authorization", getAuthHeader())
	sPass.Router.ServeHTTP(w, req)
	// check status code
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// test no username
func testPassGetUsernameEmpty(t *testing.T) {
	sPass := State{
		db: new(passDB),
	}
	sPass.Router = setupTestRouter(sPass, t)
	w := httptest.NewRecorder()
	reqNoUser := httptest.NewRequest("GET", "/register//get", nil)
	reqNoUser.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
	sPass.Router.ServeHTTP(w, reqNoUser)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// test not match header
func testFailGetUsernameNotMatch(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	reqNoUser := httptest.NewRequest("GET", "/register/523/get", nil)
	reqNoUser.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
	sFail.Router.ServeHTTP(w, reqNoUser)
	if status := w.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

// test no  header
func testFailGetUsernameNoHeader(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	reqNoUser := httptest.NewRequest("GET", "/register/captainnobody1/get", nil)
	sFail.Router.ServeHTTP(w, reqNoUser)
	if status := w.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
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
	req.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
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
	req.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
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
	req.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
	sPass.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// test delete user no header
func testFailDeleteUsernameNoHeader(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/register/captainnobody1/delete", nil)
	sFail.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
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
	req.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
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
	req.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
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
	req.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
	sFail.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func testFailStartGameNoHeader(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/game/start", nil)
	sFail.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func testPassJoinGame(t *testing.T) {
	sPass := State{
		db: new(passDB),
	}
	sPass.Router = setupTestRouter(sPass, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/game/1234/join", nil)
	req.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
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
	reqNoUser.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
	sFail.Router.ServeHTTP(w, reqNoUser)
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func testFailJoinGameNoHeader(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/game/1232/join", nil)
	sFail.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}
func testFailJoinGameDB(t *testing.T) {
	sFail := State{
		db: new(failDB),
	}
	sFail.Router = setupTestRouter(sFail, t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/game/1232/join", nil)
	req.Header.Set("Authorization", "Basic Y2FwdGFpbm5vYm9keTE6cGFzc3dvcmQK")
	sFail.Router.ServeHTTP(w, req)

	// check status code
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}
