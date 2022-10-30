package basicfuncs

import "testing"

func TestGenerateRandomString(t *testing.T) {
	s := GenerateRandomString(10)
	if len(s) != 10 {
		t.Errorf("Expected string length 10, got %d", len(s))
	}
}
