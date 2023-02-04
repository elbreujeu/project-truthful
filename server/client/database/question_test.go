package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAddQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO question").WithArgs("question", "ip address", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	id, err := AddQuestion("question", 0, "ip address", true, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while adding question: %s", err.Error())
	}
	if id != 1 {
		t.Errorf("Id should be 1")
	}

	mock.ExpectExec("INSERT INTO question").WithArgs("question", 2, "ip address", true, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	id, err = AddQuestion("question", 2, "ip address", true, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while adding question: %s", err.Error())
	}
	if id != 1 {
		t.Errorf("Id should be 1")
	}

	mock.ExpectExec("INSERT INTO question").WithArgs("question", 3, "ip address", false, 1).WillReturnError(errors.New("error"))
	_, err = AddQuestion("question", 3, "ip address", false, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGetQuestionReceiverId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error
	mock.ExpectQuery("SELECT receiver_id FROM question").WithArgs(1).WillReturnError(errors.New("error for db test"))
	_, err = GetQuestionReceiverId(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with one row returned
	mock.ExpectQuery("SELECT receiver_id FROM question").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"author_id"}).AddRow(1))
	receiverId, err := GetQuestionReceiverId(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
	if receiverId != 1 {
		t.Errorf("Database error: expected 1, got %d", receiverId)
	}
}
