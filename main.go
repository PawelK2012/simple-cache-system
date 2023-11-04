package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID       int
	UserName string
}

type Server struct {
	db         map[int]*User
	cache      map[int]*User
	dbHit      int
	listenAddr string
}

type ApiError struct {
	Error string `json:"error"`
}

func main() {
	server := NewServer()
	server.Run()
}

func NewServer() *Server {
	db := make(map[int]*User)
	for i := 0; i < 100; i++ {
		db[i+1] = &User{
			ID:       i + 1,
			UserName: fmt.Sprintf("user_%d", i+1),
		}
	}

	return &Server{
		db:         db,
		cache:      make(map[int]*User),
		listenAddr: ":3000",
	}
}

func (s *Server) tryCache(id int) (*User, bool) {
	user, ok := s.cache[id]
	return user, ok
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	// first hit cache
	user, ok := s.tryCache(id)
	if ok {
		fmt.Printf(" user %v found in cache \n", id)
		json.NewEncoder(w).Encode(user)
		return nil
	}

	// hit db if not found in cache
	user, ok = s.db[id]
	if !ok {
		return fmt.Errorf("user %d not found", id)
	}
	s.dbHit++

	// insert in cache
	s.cache[id] = user

	json.NewEncoder(w).Encode(user)
	return nil

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//handle error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (s *Server) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/user", makeHTTPHandleFunc(s.handleGetUser))

	log.Println("API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}
