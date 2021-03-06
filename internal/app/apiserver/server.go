package apiserver

import (
	"encoding/json"
	"errors"
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
	s.router.HandleFunc("/login", s.SessionsCreate()).Methods("GET")
	s.router.HandleFunc("/save", s.SaveCache()).Methods("GET")
	s.router.HandleFunc("/keys", s.GetKeys()).Methods("GET")

	private := s.router.PathPrefix("/rows").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/{key}", s.GetValueByKey()).Methods("GET")
	private.HandleFunc("/hash/{hashName}/{field}", s.GetHashValue()).Methods("GET")
	private.HandleFunc("/hash", s.CreateHashRow()).Methods("POST")
	private.HandleFunc("", s.CreateRow()).Methods("POST")
	private.HandleFunc("/{key}", s.DeleteRow()).Methods("DELETE")
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, ok := session.Values["status"]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Not authenticated")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *server) SessionsCreate() func(http.ResponseWriter, *http.Request) {
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
		session.Values["status"] = "authorized"
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
		value, err := s.cache.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, value)
	}
}

func (s *server) GetHashValue() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		hash := vars["hashName"]
		field := vars["field"]
		value, err := s.cache.HGet(hash, field)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, value)
	}
}

func (s *server) GetKeys() func(http.ResponseWriter, *http.Request) {
	type Response struct{
		Keys []string `json:"keys"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		pattern := r.URL.Query().Get("pattern")
		if pattern == "" {
			err :=  errors.New("parameter pattern is empty")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp := Response{s.cache.GetKeys(pattern)}
		jsResp, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsResp)
	}
}


func (s *server) CreateRow() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		key := r.URL.Query().Get("key")
		if key == "" {
			err :=  errors.New("parameter key is empty")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		value := r.URL.Query().Get("value")
		if value == "" {
			err :=  errors.New("parameter value is empty")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ttl, _ := strconv.ParseInt(r.URL.Query().Get("ttl"), 10, 64)
		if ttl < 0 {
			err := errors.New("invalid ttl in set")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.cache.Set(key, value, ttl)
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *server) CreateHashRow() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		hash := r.URL.Query().Get("hash")
		if hash == "" {
			err :=  errors.New("parameter hash is empty")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		field := r.URL.Query().Get("field")
		if field == "" {
			err :=  errors.New("parameter field is empty")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		value := r.URL.Query().Get("value")
		if value == "" {
			err :=  errors.New("parameter value is empty")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := s.cache.HSet(hash, field, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *server) DeleteRow() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		key := vars["key"]
		err := s.cache.Delete(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) SaveCache() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := s.cache.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}