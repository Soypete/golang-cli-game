package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// requires auth token to access the db
func (s State) getUsername(w http.ResponseWriter, r *http.Request) {
	username, err := getAndValidateUsername(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: return all values from db
	userData, err := s.db.GetUserData(username)
	if err != nil {
		handle500Err(w, " unable to get user")
		return
	}
	counter200Code.Add(1)
	w.Write([]byte(userData))
}

// does not require header
func (s State) updateUsername(w http.ResponseWriter, r *http.Request) {

	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username needs to be provided", http.StatusBadRequest)
		counter400Code.Add(1)
	}

	password := r.URL.Query().Get("password")
	if password == "" {
		log.Println("password parameter is empty")
		password = genPassword()
	}

	err := s.db.UpsertUsername(username, password)
	if err != nil {
		handle500Err(w, " unable to update user")
		return
	}
	counter200Code.Add(1)
	w.Write([]byte("user %s registered with password %s ")) //TODO: return updated username
}

func (s State) deleteUsername(w http.ResponseWriter, r *http.Request) {
	username, err := getAndValidateUsername(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		counter400Code.Add(1)
		return
	}
	err = s.db.DeleteUsername(username)
	if err != nil {
		handle500Err(w, " unable to delete user")
		return
	}
	counter200Code.Add(1)
	w.Write([]byte("user deleted")) // TODO: return deleted username
}

// TODO: dont allow them to start a game without a userID
// /game/start
// userID is attached to permissions which are in headers
func (s State) startGame(w http.ResponseWriter, r *http.Request) {
	username, err := usernameFromHeader(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	//  return gameID and an error
	GameID, err := s.db.CreateGame(username)
	if err != nil {
		handle500Err(w, " unable to start game")
		return
	}
	respTest := fmt.Sprintf("Game started with id %d.\n share this link so others can join %s\n", GameID, s.makeGamePath(GameID))
	counter200Code.Add(1)
	w.WriteHeader(http.StatusCreated) // Created
	w.Write([]byte(respTest))
}

// /game/gameID/join?
func (s State) joinGame(w http.ResponseWriter, r *http.Request) {
	username, err := usernameFromHeader(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	gameID, err := getAndValidateGameID(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.db.AddUserToGame(username, gameID)
	if err != nil {
		handle500Err(w, " unable to join game")
		return
	}
	counter200Code.Add(1)
	respText := fmt.Sprintf("User %s joined game %d", username, gameID)
	w.Write([]byte(respText))
}

func (s State) getGameState(w http.ResponseWriter, r *http.Request) {
	_, err := usernameFromHeader(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	gameID, err := getAndValidateGameID(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gameData, err := s.db.GetGameData(gameID)
	if err != nil {
		handle500Err(w, " unable to get game")
		return
	}
	gameJson, err := json.Marshal(gameData)
	if err != nil {
		handle500Err(w, " unable to marshal game data")
		return
	}
	counter200Code.Add(1)
	w.Write([]byte(gameJson))
}

func (s State) stopGame(w http.ResponseWriter, r *http.Request) {
	_, err := usernameFromHeader(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	gameID, err := getAndValidateGameID(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.db.StopGame(gameID)
	if err != nil {
		handle500Err(w, "unable to delete game")
		return
	}
	counter200Code.Add(1)
	w.Write([]byte("game deleted"))
}

// /game/{gameID}/play?answer=...
func (s State) playGame(w http.ResponseWriter, r *http.Request) {
	username, err := usernameFromHeader(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	gameID, err := getAndValidateGameID(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	answer := r.URL.Query().Get("answer")
	// send answer to db
	// TODO: add more db functions
	fmt.Println(username, answer, gameID)
}
