package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCheckAdminStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	_, err = CheckAdminStatus(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking admin status: %s", err.Error())
	}
}

func TestCheckAdminStatusError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnError(err)
	_, err = CheckAdminStatus(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestCheckModeratorStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	_, err = CheckModeratorStatus(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking moderator status: %s", err.Error())
	}
}

func TestCheckModeratorStatusError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnError(err)
	_, err = CheckModeratorStatus(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestPromoteUserToAdmin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectExec("UPDATE user").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = PromoteUserToAdmin(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while promoting user to admin: %s", err.Error())
	}
}

func TestPromoteUserToAdminError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectExec("UPDATE user").WithArgs(1).WillReturnError(err)
	err = PromoteUserToAdmin(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestPromoteUserToModerator(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectExec("UPDATE user").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = PromoteUserToModerator(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while promoting user to moderator: %s", err.Error())
	}
}

func TestPromoteUserToModeratorError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectExec("UPDATE user").WithArgs(1).WillReturnError(err)
	err = PromoteUserToModerator(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}
