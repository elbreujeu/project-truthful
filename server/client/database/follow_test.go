package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCheckFollowExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckFollowExists(1, 2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if follow exists: %s", err.Error())
	}
	if !exists {
		t.Errorf("Follow should exist")
	}
	// tests that the follow does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckFollowExists(1, 3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if follow exists: %s", err.Error())
	}
	if exists {
		t.Errorf("Follow should not exist")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 4).WillReturnError(errors.New("error"))
	_, err = CheckFollowExists(1, 4, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestAddFollow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectExec("INSERT INTO follow").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = AddFollow(1, 2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while adding follow: %s", err.Error())
	}

	mock.ExpectExec("INSERT INTO follow").WithArgs(1, 3).WillReturnError(errors.New("error"))
	err = AddFollow(1, 3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestRemoveFollow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectExec("DELETE FROM follow").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = RemoveFollow(1, 2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while removing follow: %s", err.Error())
	}

	mock.ExpectExec("DELETE FROM follow").WithArgs(1, 3).WillReturnError(errors.New("error"))
	err = RemoveFollow(1, 3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}
