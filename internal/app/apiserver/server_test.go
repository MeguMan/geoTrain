package apiserver_test

import (
	"fmt"
	"github.com/MeguMan/geoTrain/internal/app/apiserver"
	"github.com/MeguMan/geoTrain/internal/app/memcache"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var cookie string = "geo=MTYwNzg3OTE3NXxEdi1CQkFFQ180SUFBUkFCRUFBQUt2L" +
"UNBQUVHYzNSeWFXNW5EQWdBQm5OMFlYUjFjd1p6ZEhKcGJtY01EQUFLWVhWMGFH" +
"OXlhWHBsWkE9PXyQ24-AdFGX7gf2Ucz1LIN6gjttNcJj9iFa5HrtL9TmeA=="

func TestServer_SessionsCreate(t *testing.T) {
	s := apiserver.NewServer(memcache.TestLru(t), sessions.NewCookieStore([]byte("4f02bd02c9cef5c05311c2225e05caa5")))

	testCases := []struct {
		name         string
		password     string
		expectedCode int
	}{
		{
			name: "valid",
			password: "supersecretpassword",
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid password",
			password: "superwrongpassword",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/login?password="+tc.password, nil)
			if err != nil {
				t.Fatal(err)
			}
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_CreateRow(t *testing.T) {
	s := apiserver.NewServer(memcache.TestLru(t), sessions.NewCookieStore([]byte("4f02bd02c9cef5c05311c2225e05caa5")))

	testCases := []struct {
		name         string
		key          string
		value        string
		ttl          string
		expectedCode int
	}{
		{
			name: "valid",
			key: "key1",
			value: "value1",
			ttl: "0",
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid",
			key: "key1",
			value: "",
			ttl: "0",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			key: "",
			value: "value1",
			ttl: "-5",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("POST", fmt.Sprintf("/rows?key=%s&value=%s", tc.key, tc.value), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Cookie", cookie)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_CreateHashRow(t *testing.T) {
	s := apiserver.NewServer(memcache.TestLru(t), sessions.NewCookieStore([]byte("4f02bd02c9cef5c05311c2225e05caa5")))

	testCases := []struct {
		name         string
		hash         string
		field        string
		value        string
		expectedCode int
	}{
		{
			name: "valid",
			hash: "myhash",
			field: "field1",
			value: "value1",
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid",
			hash: "myhash",
			field: "field1",
			value: "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			hash: "myhash",
			field: "",
			value: "value1",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			hash: "",
			field: "field1",
			value: "value1",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("POST", fmt.Sprintf("/rows/hash?hash=%s&field=%s&value=%s", tc.hash, tc.field, tc.value), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Cookie", cookie)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_GetHashValue(t *testing.T) {
	s := apiserver.NewServer(memcache.TestLru(t), sessions.NewCookieStore([]byte("4f02bd02c9cef5c05311c2225e05caa5")))
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/rows/hash?hash=myhash&field=field1&value=value1"), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", cookie)
	s.ServeHTTP(rec, req)

	testCases := []struct {
		name         string
		hash         string
		field        string
		expectedCode int
	}{
		{
			name: "valid",
			hash: "myhash",
			field: "field1",
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid",
			hash: "myhash123",
			field: "field1",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "invalid",
			hash: "myhash6124",
			field: "field1",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "invalid",
			hash: "myhash",
			field: "",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", fmt.Sprintf("/rows/hash/%s/%s", tc.hash, tc.field), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Cookie", cookie)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_GetValueByKey(t *testing.T) {
	s := apiserver.NewServer(memcache.TestLru(t), sessions.NewCookieStore([]byte("4f02bd02c9cef5c05311c2225e05caa5")))
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/rows?key=mykey&value=value1"), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", cookie)
	s.ServeHTTP(rec, req)

	testCases := []struct {
		name         string
		key         string
		expectedCode int
	}{
		{
			name: "valid",
			key: "mykey",
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid",
			key: "nokey",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", fmt.Sprintf("/rows/%s", tc.key), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Cookie", cookie)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_GetKeys(t *testing.T) {
	s := apiserver.NewServer(memcache.TestLru(t), sessions.NewCookieStore([]byte("4f02bd02c9cef5c05311c2225e05caa5")))
	testCases := []struct {
		name         string
		pattern         string
		expectedCode int
	}{
		{
			name: "valid",
			pattern: "*",
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid",
			pattern: "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", fmt.Sprintf("/keys?pattern=%s", tc.pattern), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Cookie", cookie)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_DeleteRow(t *testing.T) {
	s := apiserver.NewServer(memcache.TestLru(t), sessions.NewCookieStore([]byte("4f02bd02c9cef5c05311c2225e05caa5")))
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/rows?key=mykey&value=value1"), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", cookie)
	s.ServeHTTP(rec, req)

	testCases := []struct {
		name         string
		key         string
		expectedCode int
	}{
		{
			name: "valid",
			key: "mykey",
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid",
			key: "mykey",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/rows/%s", tc.key), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Cookie", cookie)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_SaveCache(t *testing.T) {
	s := apiserver.NewServer(memcache.TestLru(t), sessions.NewCookieStore([]byte("4f02bd02c9cef5c05311c2225e05caa5")))
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/save", nil)
	if err != nil {
		t.Fatal(err)
	}
	s.ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code)
}

