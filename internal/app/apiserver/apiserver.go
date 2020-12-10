package apiserver

import (
	"github.com/MeguMan/geoTrain/internal/app/memcache"
	"net/http"
)

func Start() error {
	cache := memcache.NewLru(256)
	server := NewServer(cache)
	return http.ListenAndServe(":8080", server)
}
