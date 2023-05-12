package server

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func getAndValidateUsername(w http.ResponseWriter, r *http.Request) (string, error) {
	username := chi.URLParam(r, "username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.New(http.StatusText(400) + "username parameter cannot be empty")
		return "", err
	}
	return username, nil
}

func getAndValidateGameID(w http.ResponseWriter, r *http.Request) (int64, error) {
	gameID := chi.URLParam(r, "gameID")
	if gameID == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.New(http.StatusText(400) + ", gameID parameter cannot be empty")
		return 0, err
	}
	gameID64, err := strconv.ParseInt(gameID, 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		err := errors.New(http.StatusText(400) + ", gameID parameter must be an integer")
		return 0, err
	}
	return gameID64, nil
}
