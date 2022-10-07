package token

import (
	"net/http"
	"testing"
)

func TestParseAccessToken(t *testing.T) {
	// creates a new request with an empty header
	r, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = ParseAccessToken(r)
	if err == nil {
		t.Errorf("no error while parsing nil request")
	}

	r.Header.Set("Authorization", "test")
	_, _, err = ParseAccessToken(r)
	if err == nil {
		t.Errorf("no error while parsing invalid token")
	}

	r.Header.Set("Authorization", "Bearer test")
	token, _, err := ParseAccessToken(r)
	if err != nil {
		t.Errorf("error while parsing valid token")
	}
	if token != "test" {
		t.Errorf("invalid token parsed")
	}
}
