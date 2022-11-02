package client

import (
	"database/sql"
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

func TestAnswerQuestionError(t *testing.T) {
	_, _, err := AnswerQuestion(0, 0, "", "ip_address")

	if err == nil {
		t.Error("Empty text for answer: expected error, got nil")
	}

	var mock sqlmock.Sqlmock
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

	// check user database error
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnError(errors.New("test error"))
	_, _, err = AnswerQuestion(2, 0, "toto", "ip_address")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// question does not exists
	mock.ExpectQuery("SELECT COUNT").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT receiver_id").WithArgs(1).WillReturnError(sql.ErrNoRows)
	_, _, err = AnswerQuestion(3, 1, "toto", "ip_address")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Question does not exist: expected error, got nil")
	}

	// check question database error
	mock.ExpectQuery("SELECT COUNT").WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT receiver_id").WithArgs(2).WillReturnError(errors.New("test error"))
	_, _, err = AnswerQuestion(4, 2, "toto", "ip_address")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// question id and receiver id do not match
	mock.ExpectQuery("SELECT COUNT").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT receiver_id").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"receiver_id"}).AddRow(2))
	_, _, err = AnswerQuestion(5, 3, "toto", "ip_address")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Question id and receiver id do not match: expected error, got nil")
	}

	// user already answered the question
	mock.ExpectQuery("SELECT COUNT").WithArgs(6).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT receiver_id").WithArgs(7).WillReturnRows(sqlmock.NewRows([]string{"receiver_id"}).AddRow(6))
	mock.ExpectQuery("SELECT COUNT").WithArgs(7).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	_, _, err = AnswerQuestion(6, 7, "toto", "ip_address")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Question id and receiver id do not match: expected error, got nil")
	}

	// check if question already answered database error
	mock.ExpectQuery("SELECT COUNT").WithArgs(8).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT receiver_id").WithArgs(9).WillReturnRows(sqlmock.NewRows([]string{"receiver_id"}).AddRow(8))
	mock.ExpectQuery("SELECT COUNT").WithArgs(9).WillReturnError(errors.New("test error"))
	_, _, err = AnswerQuestion(8, 9, "toto", "ip_address")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Question id and receiver id do not match: expected error, got nil")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(10).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT receiver_id").WithArgs(11).WillReturnRows(sqlmock.NewRows([]string{"receiver_id"}).AddRow(10))
	mock.ExpectQuery("SELECT COUNT").WithArgs(11).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectExec("INSERT INTO answer").WithArgs(10, 11, "toto", "ip_address").WillReturnError(errors.New("test error"))
	_, _, err = AnswerQuestion(10, 11, "toto", "ip_address")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestAnswerQuestionSuccess(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(6).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT receiver_id").WithArgs(7).WillReturnRows(sqlmock.NewRows([]string{"receiver_id"}).AddRow(6))
	mock.ExpectQuery("SELECT COUNT").WithArgs(7).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectExec("INSERT INTO answer").WithArgs(6, 7, "toto", "ip_address").WillReturnResult(sqlmock.NewResult(1, 1))
	answerId, code, err := AnswerQuestion(6, 7, "toto", "ip_address")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while answering question: %s", err.Error())
	}
	if code != 201 {
		t.Errorf("Wrong status code: expected 200, got %d", code)
	}
	if answerId != 1 {
		t.Errorf("Wrong answer id: expected 1, got %d", answerId)
	}
}
