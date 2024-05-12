package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPardonUserDbError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO pardon").WithArgs(1, 1).WillReturnError(errors.New("error"))
	_, err = PardonUser(1, 1, db)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestPardonUserLastInsertIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO pardon").WithArgs(1, 1).WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))
	_, err = PardonUser(1, 1, db)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestPardonUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO pardon").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = PardonUser(1, 1, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestCheckPardonExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// pardon exists
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	_, err = CheckPardonExists(1, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}

	// pardon does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	_, err = CheckPardonExists(2, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// db error
	mock.ExpectQuery("SELECT COUNT").WithArgs(3).WillReturnError(errors.New("error"))
	_, err = CheckPardonExists(3, db)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
