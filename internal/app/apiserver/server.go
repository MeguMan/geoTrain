package apiserver

import (
	"fmt"
	"github.com/MeguMan/geoTrain/internal/app/memcache"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

const sessionName = "geo"

type server struct {
	router       *mux.Router
	cache        *memcache.LRU
	sessionStore sessions.Store
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func NewServer(cache *memcache.LRU, sessionStore sessions.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		cache: cache,
		sessionStore: sessionStore,
	}
	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/login", s.HandleSessionsCreate()).Methods("GET")
	s.router.HandleFunc("/rows/{key}", s.GetValueByKey()).Methods("GET")
	s.router.HandleFunc("/rows", s.CreateRow()).Methods("POST")
	s.router.HandleFunc("/rows/{key}", s.DeleteRow()).Methods("DELETE")
	s.router.HandleFunc("/save", s.SaveCache()).Methods("GET")
}

func (s *server) HandleSessionsCreate() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		password := r.URL.Query().Get("password")

		authorized := s.cache.CheckPassword(password)

		if !authorized {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Wrong password")
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session.Values["shit"] = "suka"
		if err := s.sessionStore.Save(r, w, session); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
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

func (s *server) SaveCache() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		s.cache.Save()
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Saved")
	}
}