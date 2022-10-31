package client

import (
	"errors"
	"project_truthful/client/database"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCheckAnswerInfos(t *testing.T) {
	err := checkAnswerInfos("")
	if err == nil {
		t.Error("Empty check answer info: expected error, got nil")
	}

	tooLongStr := strings.Repeat("a", 1500)
	err = checkAnswerInfos(tooLongStr)
	if err == nil {
		t.Error("Too long answer: expected error, got nil")
	}

	err = checkAnswerInfos("Toto")
	if err != nil {
		t.Error("Good answer text: expected nil, got error")
	}
}

func TestAnswerQuestionEmptyText(t *testing.T) {
	_, _, err := AnswerQuestion(0, 0, "", "ip_address")

	if err == nil {
		t.Error("Empty text for answer: expected error, got nil")
	}
}

func TestAnswerQuestionUserDoesNotExists(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer database.DB.Close()

	// user does not exists
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	_, _, err = AnswerQuestion(1, 0, "toto", "ip_address")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("User does not exist: expected error, got nil")
	}

	// database error
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnError(errors.New("test error"))
	_, _, err = AnswerQuestion(2, 0, "toto", "ip_address")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}
}
