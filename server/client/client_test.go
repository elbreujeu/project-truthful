package client

import (
	"os"
	"testing"
)

func TestGetServerVersion(t *testing.T) {
	os.Setenv("SERVER_VERSION", "1")
	version := GetServerVersion()
	if version != "1" {
		t.Errorf("Expected version 1, got %s", version)
	}
}
