package apiserver

import (
	"net/http"
)

func Start() error {
	server := NewServer()
	return http.ListenAndServe(":8080", server)
}
