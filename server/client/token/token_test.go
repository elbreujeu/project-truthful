package token

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParseAccessToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test missing header
	_, status, err := ParseAccessToken(c)
	if err == nil {
		t.Errorf("Expected error when header is missing, got nil")
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}

	// Test header with incorrect format
	req.Header.Set("Authorization", "incorrect format")
	_, status, err = ParseAccessToken(c)
	if err == nil {
		t.Errorf("Expected error when header has incorrect format, got nil")
	}
	if status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}

	// Test valid header
	req.Header.Set("Authorization", "Bearer token")
	token, status, err := ParseAccessToken(c)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
	if token != "token" {
		t.Errorf("Expected token %s, got %s", "token", token)
	}
}
