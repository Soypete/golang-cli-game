package server

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

// this only needs to happen once per program execution
func init() {
	rand.Seed(time.Now().UnixNano())
}

func (s State) makeGamePath(gameID int64) string {
	return fmt.Sprintf("%s/game/%d", s.BaseURL, gameID)
}

func usernameFromHeader(w http.ResponseWriter, r *http.Request) (string, error) {
	username, _, ok := r.BasicAuth()
	if !ok {
		return "", errors.New("Authorization header must be in the form username:password")
	}
	return username, nil
}

// we want the header to include basic auth - username:password
// we wnt this method to return the username
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication
func (s *State) authMiddlelware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//https://pkg.go.dev/net/http#Request.BasicAuth
		username, password, ok := r.BasicAuth()
		if !ok {
			counter400Code.Add(1)
			http.Error(w, http.StatusText(http.StatusBadRequest)+", Authorization header must be in the form username:password", http.StatusBadRequest)
			return
		}
		if isValid, err := s.db.CheckUserValid(username, password); err != nil || !isValid {
			counter400Code.Add(1)
			http.Error(w, http.StatusText(http.StatusUnauthorized)+", Username or password do not exist", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func getAndValidateUsername(w http.ResponseWriter, r *http.Request) (string, error) {
	headerUsername, err := usernameFromHeader(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		counter400Code.Add(1)
		err := errors.New(http.StatusText(http.StatusBadRequest) + err.Error())
		return "", err
	}
	username := chi.URLParam(r, "username")
	if username == "" {
		username = headerUsername
		return username, nil
	}
	// stop user from editing other users' data
	if headerUsername != username {
		counter400Code.Add(1)
		return "", errors.New("username does not match")
	}

	// TODO: should be able to remove if auth is on all endpoints
	// check if username is empty
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		counter400Code.Add(1)
		err := errors.New(http.StatusText(http.StatusBadRequest) + ", username parameter cannot be empty")
		return "", err
	}
	fmt.Println(username)
	return username, nil
}

func getAndValidateGameID(w http.ResponseWriter, r *http.Request) (int64, error) {
	gameID := chi.URLParam(r, "gameID")
	if gameID == "" {
		w.WriteHeader(http.StatusBadRequest)
		counter400Code.Add(1)
		err := errors.New(http.StatusText(http.StatusBadRequest) + ", gameID parameter cannot be empty")
		return 0, err
	}
	gameID64, err := strconv.ParseInt(gameID, 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		counter400Code.Add(1)
		err := errors.New(http.StatusText(http.StatusBadRequest) + ", gameID parameter must be an integer")
		return 0, err
	}
	return gameID64, nil
}

var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789!@#$%^&*()_+")
var randomNames = []string{"ShadowDragon", "CrimsonReaper", "NightAssassin", "SavageHunter", "Thunderbolt", "StormBringer", "IceQueen", "FireDemon", "TheOneTrueHero", "SilverWolf", "GoldenKnight", "MasterMind", "CyborgAssassin", "ElectricEagle", "GalacticGamer", "NeonNinja", "DarkPhoenix", "DiamondDragon", "ChaosKing", "MysticMage", "IronGiant", "CelestialSiren", "ShadowHunter", "DeathWish", "SnowLeopard", "CosmicCrusader", "EternalKnight", "PhoenixBlaze", "ThunderStorm"}

func genPassword() string {
	var password string
	for i := 0; i < 10; i++ {
		password += string(symbols[rand.Intn(len(symbols))])
	}
	return password
}

func (s State) genUsername() (string, error) {
	// add new user if not exists
	username := randomNames[rand.Intn(len(randomNames))]
	password := genPassword()
	err := s.db.UpsertUsername(username, password)
	if err != nil {
		return "", err
	}
	return username, nil
}

func handle500Err(w http.ResponseWriter, err string) {
	log.Println(errors.New(err))
	http.Error(w, err, http.StatusInternalServerError)
	counter500Code.Add(1)
}
