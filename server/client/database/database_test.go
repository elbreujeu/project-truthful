package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectExec("INSERT INTO user").WithArgs("username", "username", "password", "email", "birthdate").WillReturnResult(sqlmock.NewResult(4, 1))
	id, err := InsertUser("username", "password", "email", "birthdate", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if id != 4 {
		t.Errorf("id should be 4, but is %d", id)
	}
	if err != nil {
		t.Errorf("Error while inserting user: %s", err.Error())
	}
}

func TestCheckUsernameExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists := CheckUsernameExists("username", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if !exists {
		t.Errorf("Username should exist")
	}
	// tests that the username does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists = CheckUsernameExists("toto", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if exists {
		t.Errorf("Username should not exist")
	}
}

func TestCheckEmailExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs("email@email.fr").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists := CheckEmailExists("email@email.fr", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if !exists {
		t.Errorf("Email should exist")
	}
	// tests that the email does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto@toto.fr").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists = CheckEmailExists("toto@toto.fr", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if exists {
		t.Errorf("Email should not exist")
	}
}
