package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

func (s State) getUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, http.StatusText(400)+", username parameter cannot be empty", 400)
		return
	}
	// TODO: return all values from db
	userData, err := s.db.GetUserData(username)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to get user", 500)
		return
	}
	w.Write([]byte(userData))
}

func (s State) updateUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, http.StatusText(400)+", username parameter cannot be empty", 400)
		return
	}
	err := s.db.UpsertUsername(username)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to update user", 500)
		return
	}
	w.Write([]byte("user registered")) //TODO: return updated username
}

func (s State) deleteUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	err := s.db.DeleteUsername(username)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to delete user", 500)
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
		http.Error(w, http.StatusText(500)+" unable to start game", 500)
		return
	}
	respTest := fmt.Sprintf("Game started with id %d.\n share this link so others can join %s\n", GameID, s.makeGamePath(GameID))
	w.WriteHeader(201) // Created
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
	gameID := chi.URLParam(r, "gameID")
	if gameID == "" {
		http.Error(w, http.StatusText(400)+", gameID parameter cannot be empty", 400)
		return
	}
	gameID64, err := strconv.ParseInt(gameID, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400)+", gameID parameter must be an integer", 400)
		return
	}
	err = s.db.AddUserToGame(username, gameID64)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to join game", 500)
		return
	}
	respText := fmt.Sprintf("User %s joined game %s", username, gameID)
	w.Write([]byte(respText))
}
