package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func (s State) getUsername(w http.ResponseWriter, r *http.Request) {
	username, err := getAndValidateUsername(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: return all values from db
	userData, err := s.db.GetUserData(username)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to get user", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(userData))
}

func (s State) updateUsername(w http.ResponseWriter, r *http.Request) {
	username, err := getAndValidateUsername(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.db.UpsertUsername(username)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to update user", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("user registered")) //TODO: return updated username
}

func (s State) deleteUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	err := s.db.DeleteUsername(username)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to delete user", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("user deleted")) // TODO: return deleted username
}

func (s State) makeGamePath(gameID int64) string {
	return fmt.Sprintf("%s/game/%d", s.BaseURL, gameID)
}

// /game/start?username=...
// userID is attached to permissions which are in headers
func (s State) startGame(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("username parameter is empty")
		// TODO: retrieve userID from header and use that to get username
	}
	//  return gameID and an error
	GameID, err := s.db.CreateGame(username)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to start game", http.StatusInternalServerError)
		return
	}
	respTest := fmt.Sprintf("Game started with id %d.\n share this link so others can join %s\n", GameID, s.makeGamePath(GameID))
	w.WriteHeader(http.StatusCreated) // Created
	w.Write([]byte(respTest))
}

var randomNames = []string{"ShadowDragon", "CrimsonReaper", "NightAssassin", "SavageHunter", "Thunderbolt", "StormBringer", "IceQueen", "FireDemon", "TheOneTrueHero", "SilverWolf", "GoldenKnight", "MasterMind", "CyborgAssassin", "ElectricEagle", "GalacticGamer", "NeonNinja", "DarkPhoenix", "DiamondDragon", "ChaosKing", "MysticMage", "IronGiant", "CelestialSiren", "ShadowHunter", "DeathWish", "SnowLeopard", "CosmicCrusader", "EternalKnight", "PhoenixBlaze", "ThunderStorm"}

// /game/gameID/join?username=...
func (s State) joinGame(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("username parameter is empty")
		// TODO: retrieve userID from header and use that to get username

		// add new user if not exists
		username = randomNames[rand.Intn(len(randomNames))]
		s.db.UpsertUsername(username)
	}
	gameID, err := getAndValidateGameID(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.db.AddUserToGame(username, gameID)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to join game", http.StatusInternalServerError)
		return
	}
	respText := fmt.Sprintf("User %s joined game %s", username, gameID)
	w.Write([]byte(respText))
}

func (s State) getGameState(w http.ResponseWriter, r *http.Request) {
	gameID, err := getAndValidateGameID(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gameData, err := s.db.GetGameData(gameID)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to get game", http.StatusInternalServerError)
		return
	}
	gameJson, err := json.Marshal(gameData)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to marshal game data", 500)
		return
	}
	w.Write([]byte(gameJson))
}

func (s State) stopGame(w http.ResponseWriter, r *http.Request) {
	gameID, err := getAndValidateGameID(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.db.StopGame(gameID)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to delete game", 500)
		return
	}
	w.Write([]byte("game deleted"))
}

func (s State) playGame(w http.ResponseWriter, r *http.Request) {
	gameID, err := getAndValidateGameID(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := r.URL.Query().Get("username")
	// check if the user is the game host
	answer := r.URL.Query().Get("answer")
	// send answer to db
	// TODO: add more db functions
	fmt.Println(username, answer, gameID)
}
