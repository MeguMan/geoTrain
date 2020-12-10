package apiserver

import (
	"fmt"
	"github.com/MeguMan/geoTrain/internal/app/memcache"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type server struct {
	router *mux.Router
	cache  *memcache.LRU
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func NewServer(cache *memcache.LRU) *server {
	s := &server{
		router: mux.NewRouter(),
		cache: cache,
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
		vars := mux.Vars(r)
		key := vars["key"]

		value := s.cache.Get(key)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "GOT ", value)
	}
}

func (s *server) CreateRow() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		expiration, _ := strconv.ParseInt(r.URL.Query().Get("ttl"), 10, 64)

		s.cache.Set(key, value, expiration)

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
