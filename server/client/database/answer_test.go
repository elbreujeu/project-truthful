package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCheckAnswerIdExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnError(errors.New("error for db test"))
	_, err = CheckAnswerIdExists(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with one row returned
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckAnswerIdExists(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
	if !exists {
		t.Errorf("Database error: expected true, got false")
	}

	// Test with no row returned
	mock.ExpectQuery("SELECT COUNT").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckAnswerIdExists(3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
	if exists {
		t.Errorf("Database error: expected false, got true")
	}
}

func TestAddAnswer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with query fail
	mock.ExpectExec("INSERT INTO answer").WithArgs(1, 1, "content", "ip_address").WillReturnError(errors.New("error for db test"))
	_, err = AddAnswer(1, 1, "content", "ip_address", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with no error
	mock.ExpectExec("INSERT INTO answer").WithArgs(2, 2, "content", "ip_address").WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = AddAnswer(2, 2, "content", "ip_address", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
}

func TestHasQuestionBeenAnswered(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := HasQuestionBeenAnswered(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if question has been answered: %s", err.Error())
	}
	if !exists {
		t.Errorf("Question should have been answered")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = HasQuestionBeenAnswered(3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if question has been answered: %s", err.Error())
	}
	if exists {
		t.Errorf("Question should not have been answered")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(4).WillReturnError(errors.New("error"))
	_, err = HasQuestionBeenAnswered(4, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGetAnswerAuthorId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectQuery("SELECT user_id FROM answer").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"author_id"}).AddRow(1))
	authorId, err := GetAnswerAuthorId(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while getting answer author id: %s", err.Error())
	}
	if authorId != 1 {
		t.Errorf("Author id should be 1")
	}

	mock.ExpectQuery("SELECT user_id FROM answer").WithArgs(2).WillReturnError(errors.New("error"))
	_, err = GetAnswerAuthorId(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestMarkAnswerAsDeleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectExec("UPDATE answer").WithArgs(1).WillReturnError(errors.New("error"))
	err = MarkAnswerAsDeleted(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	mock.ExpectExec("UPDATE answer").WithArgs(2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = MarkAnswerAsDeleted(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error should be nil")
	}
}
