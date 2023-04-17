package server

import (
	"fmt"
	"log"
	"net/http"

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
	fmt.Println(username)
	err := s.db.DeleteUsername(username)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500)+" unable to delete user", 500)
		return
	}
	w.Write([]byte("user deleted")) // TODO: return deleted username
}
