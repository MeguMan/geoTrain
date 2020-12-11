package apiserver

import (
	"github.com/MeguMan/geoTrain/internal/app/memcache"
	"github.com/gorilla/sessions"
	"net/http"
)

func Start(config *memcache.Config) error {
	cache := memcache.NewLru(config)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	server := NewServer(cache, sessionStore)
	return http.ListenAndServe(":8080", server)
}
