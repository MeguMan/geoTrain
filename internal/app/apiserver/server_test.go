package apiserver_test

import (
	"github.com/MeguMan/geoTrain/internal/app/apiserver"
	"github.com/MeguMan/geoTrain/internal/app/memcache"
	"net/http"
	"net/http/httptest"
	"testing"
)

var config *memcache.Config

func TestServer_SessionsCreate(t *testing.T) {

	s := apiserver.NewServer(memcache.TestLru(t), emailsender.Sender{
		Email:    "",
		Password: "",
		TLSPort:  "",
	})

	req, err := http.NewRequest("GET", "/login?password=supersecretpassword", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.)
}