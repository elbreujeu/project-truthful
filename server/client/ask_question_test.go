package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestAskQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	database.DB = db
	if err != nil {
		t.Errorf("Error initializing mock database: %s", err)
	}
	// test for checkUserIdExists error
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnError(errors.New("error"))
	_, code, err := AskQuestion("question", 1, "ip_address", true, 1)
	if code != http.StatusInternalServerError {
		t.Errorf("Expected http.StatusInternalServerError, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	// test for checkUserIdExists not found
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	_, code, err = AskQuestion("question", 1, "ip_address", true, 1)
	if code != http.StatusNotFound {
		t.Errorf("Expected http.StatusNotFound, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	// test for checkQuestionInfos error
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	_, code, err = AskQuestion("", 1, "ip_address", true, 1)
	if code != http.StatusBadRequest {
		t.Errorf("Expected http.StatusBadRequest, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	// test for AddQuestion error
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectExec("INSERT INTO question").WithArgs("question", 1, "ip_address", true, 1).WillReturnError(errors.New("error"))
	_, code, err = AskQuestion("question", 1, "ip_address", true, 1)
	if code != http.StatusInternalServerError {
		t.Errorf("Expected http.StatusInternalServerError, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	// test for success
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectExec("INSERT INTO question").WithArgs("question", 1, "ip_address", true, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	id, code, err := AskQuestion("question", 1, "ip_address", true, 1)
	if code != http.StatusCreated {
		t.Errorf("Expected http.StatusCreated, got %d", code)
	}
	if err != nil {
		t.Errorf("Expected nil, got error")
	}
	if id != 1 {
		t.Errorf("Expected id 1, got %d", id)
	}
}
