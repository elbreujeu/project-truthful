package client

import (
	"database/sql"
	"log"
	"os"
	"project_truthful/client/database"
	"project_truthful/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLogin(t *testing.T) {
	//inits the sqlmock

	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer database.DB.Close()

	// tests that the username does not exist
	mock.ExpectQuery("SELECT id FROM user").WithArgs("username").WillReturnError(sql.ErrNoRows)
	_, _, err = Login(models.LoginInfos{Username: "username", Password: "password"})
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// tests that the password is wrong
	mock.ExpectQuery("SELECT id FROM user").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT password FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("toto"))
	_, _, err = Login(models.LoginInfos{Username: "username", Password: "password"})
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// tests that the login is successful
	hashedPassword, err := encryptPassword("password")
	if err != nil {
		t.Errorf("Error while encrypting password: %s", err.Error())
	}
	mock.ExpectQuery("SELECT id FROM user").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(44))
	mock.ExpectQuery("SELECT password FROM user").WithArgs(44).WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(hashedPassword))

	os.Setenv("IS_TEST", "true")
	_, _, err = Login(models.LoginInfos{Username: "username", Password: "password"})
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		log.Printf("Error while logging in: %s", err.Error())
		t.Errorf("Error should be nil")
	}
	os.Setenv("IS_TEST", "false")
}
