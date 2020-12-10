package apiserver

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type server struct {
	router *mux.Router
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func NewServer() *server {
	s := &server{
		router: mux.NewRouter(),
	}
	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/login/{password}", s.Login()).Methods("GET")
	s.router.HandleFunc("/rows/{key}", s.GetValueByKey()).Methods("GET")
	s.router.HandleFunc("/rows", s.CreateRow()).Methods("POST")
	s.router.HandleFunc("/rows/{key}", s.DeleteRow()).Methods("DELETE")
}

func (s *server) Login() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "LOGIN")
	}
}

func (s *server) GetValueByKey() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "GETTED")
	}
}

func (s *server) CreateRow() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "CREATED")
	}
}

func (s *server) DeleteRow() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "DELETED")
	}
}
