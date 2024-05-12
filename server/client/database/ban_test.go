package database

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBanUserDbError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// no duration
	mock.ExpectExec("INSERT INTO ban").WithArgs(1, 1, "ban reason").WillReturnError(errors.New("error for db test"))
	_, err = BanUser(1, 1, 0, "ban reason", db)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}

	// with duration
	mock.ExpectExec("INSERT INTO ban").WithArgs(1, 1, "ban reason", time.Now().Add(time.Duration(1)*time.Hour)).WillReturnError(errors.New("error for db test"))
	_, err = BanUser(1, 1, 1, "ban reason", db)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestBanUserLastInsertIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// no duration
	mock.ExpectExec("INSERT INTO ban").WithArgs(1, 1, "ban reason").WillReturnResult(sqlmock.NewErrorResult(errors.New("error for last insert id")))
	_, err = BanUser(1, 1, 0, "ban reason", db)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}

	// with duration
	mock.ExpectExec("INSERT INTO ban").WithArgs(1, 1, "ban reason", time.Now().Add(time.Duration(1)*time.Hour)).WillReturnResult(sqlmock.NewErrorResult(errors.New("error for last insert id")))
	_, err = BanUser(1, 1, 1, "ban reason", db)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestBanUserSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// no duration
	mock.ExpectExec("INSERT INTO ban").WithArgs(1, 1, "ban reason").WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = BanUser(1, 1, 0, "ban reason", db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}

	// with duration
	mock.ExpectExec("INSERT INTO ban").WithArgs(1, 1, "ban reason", time.Now().Add(time.Duration(1)*time.Hour)).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = BanUser(1, 1, 1, "ban reason", db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestCheckUserBanStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// user is banned
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	_, err = CheckUserBanStatus(1, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}

	// user is not banned
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	_, err = CheckUserBanStatus(2, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}

	// db error
	mock.ExpectQuery("SELECT COUNT").WithArgs(3).WillReturnError(errors.New("error"))
	_, err = CheckUserBanStatus(3, db)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestCheckBanExistsByBanId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// ban exists
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	_, err = CheckBanExistsByBanId(1, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}

	// ban does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	_, err = CheckBanExistsByBanId(2, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}

	// db error
	mock.ExpectQuery("SELECT COUNT").WithArgs(3).WillReturnError(errors.New("error"))
	_, err = CheckBanExistsByBanId(3, db)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}
