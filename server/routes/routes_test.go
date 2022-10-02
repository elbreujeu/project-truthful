package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

//this function tests the SetupRoutes function
func TestSetupRoutes(t *testing.T) {
	router := mux.NewRouter()
	SetupRoutes(router)
	r, _ := http.NewRequest("GET", "/hello_world", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, `{"message": "Hello world !"}`, w.Body.String())
}

func TestMiddleware(t *testing.T) {
	router := mux.NewRouter()
	SetupRoutes(router)
	SetMiddleware(router)
	r, _ := http.NewRequest("GET", "/hello_world", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
