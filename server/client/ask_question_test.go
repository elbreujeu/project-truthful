package client

import (
	"strings"
	"testing"
)

func TestCheckQuestionInfos(t *testing.T) {
	err := checkQuestionInfos("")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	str := strings.Repeat("a", 501)
	err = checkQuestionInfos(str)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	err = checkQuestionInfos("Hey there, how are you?")
	if err != nil {
		t.Error("Expected nil, got error")
	}
}
