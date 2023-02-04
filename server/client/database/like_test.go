package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCheckLikeExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 1).WillReturnError(errors.New("error for db test"))
	_, err = CheckLikeExists(1, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with one row returned
	mock.ExpectQuery("SELECT COUNT").WithArgs(2, 2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckLikeExists(2, 2, db)
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
	mock.ExpectQuery("SELECT COUNT").WithArgs(3, 3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckLikeExists(3, 3, db)
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

func TestAddLike(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error
	mock.ExpectExec("INSERT INTO answer_like").WithArgs(1, 1).WillReturnError(errors.New("error for db test"))
	err = AddLike(1, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	//Test with no error
	mock.ExpectExec("INSERT INTO answer_like").WithArgs(2, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = AddLike(2, 2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
}

func TestRemoveLike(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error
	mock.ExpectExec("DELETE FROM answer_like").WithArgs(1, 1).WillReturnError(errors.New("error for db test"))
	err = RemoveLike(1, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	//Test with no error
	mock.ExpectExec("DELETE FROM answer_like").WithArgs(2, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = RemoveLike(2, 2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
}
