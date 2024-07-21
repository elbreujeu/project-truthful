package client

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
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

func TestGoogleLoginVerifyTokenFail(t *testing.T) {
	userToken := "toto123"
	_, _, err := GoogleLogin("google", userToken)
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGoogleLoginGetOAuthProviderFail(t *testing.T) {
	os.Setenv("IS_TEST", "true")
	userToken := "toto123"

	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT id FROM oauth_provider").WithArgs("google").WillReturnError(sql.ErrNoRows)
	_, _, err = GoogleLogin("google", userToken)
	os.Setenv("IS_TEST", "false")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGoogleLoginGetUserIdBySubjectFail(t *testing.T) {
	os.Setenv("IS_TEST", "true")
	userToken := "toto123"

	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT id FROM oauth_provider").WithArgs("google").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT user_id FROM oauth_login").WithArgs(1, "123456").WillReturnError(sql.ErrNoRows)
	_, _, err = GoogleLogin("google", userToken)
	os.Setenv("IS_TEST", "false")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGoogleLoginUserExists(t *testing.T) {
	os.Setenv("IS_TEST", "true")
	userToken := "toto123"

	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT id FROM oauth_provider").WithArgs("google").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT user_id FROM oauth_login").WithArgs(1, "123456").WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

	loginToken, code, err := GoogleLogin("google", userToken)
	os.Setenv("IS_TEST", "false")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error should be nil")
	}
	if code != http.StatusOK {
		t.Errorf("Code should be http.StatusOK")
	}
	if loginToken != "test" {
		t.Errorf("Token should be \"test\"")
	}
}

func TestGoogleLoginRegisterOauthFail(t *testing.T) {
	os.Setenv("IS_TEST", "true")
	userToken := "toto123"

	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT id FROM oauth_provider").WithArgs("google").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT user_id FROM oauth_login").WithArgs(1, "123456").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto123").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto123@gmail.com").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectExec("INSERT INTO user").WithArgs("toto123", "toto123", "", "toto123@gmail.com", "2000-01-01").WillReturnError(errors.New("error for test register oauth"))

	_, code, err := GoogleLogin("google", userToken)
	os.Setenv("IS_TEST", "false")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
	if code != http.StatusInternalServerError {
		t.Errorf("Code should be http.StatusInternalServerError")
	}
}

func TestGoogleLoginRegisterOauthSuccess(t *testing.T) {
	os.Setenv("IS_TEST", "true")
	userToken := "toto123"

	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT id FROM oauth_provider").WithArgs("google").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT user_id FROM oauth_login").WithArgs(1, "123456").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto123").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto123@gmail.com").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectExec("INSERT INTO user").WithArgs("toto123", "toto123", "", "toto123@gmail.com", "2000-01-01").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO oauth_login").WithArgs(1, "123456", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	token, code, err := GoogleLogin("google", userToken)
	os.Setenv("IS_TEST", "false")
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error should be nil")
	}
	if code != http.StatusOK {
		t.Errorf("Code should be http.StatusCreated")
	}
	if token != "test" {
		t.Errorf("Token should be \"test\"")
	}
}
