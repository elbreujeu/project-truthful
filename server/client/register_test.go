package client

import (
	"errors"
	"net/http"
	"os"
	"project_truthful/client/database"
	"project_truthful/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestEncryptPassword(t *testing.T) {
	hash, err := encryptPassword("password")
	if err != nil {
		t.Errorf("Error while encrypting password: %s", err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("password"))
	if err != nil {
		t.Errorf("Error while comparing hash and password: %s", err.Error())
	}
}

func TestIsUsernameValid(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	err = isUsernameValid("us")
	if err == nil {
		t.Errorf("Error while validating username: %s", err.Error())
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	err = isUsernameValid("username")
	if err != nil {
		t.Errorf("Error while checking username: %s", err.Error())
	}
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	err = isUsernameValid("toto")
	if err == nil {
		t.Errorf("Error while checking username: %s", err.Error())
	}
}

func TestIsPasswordValid(t *testing.T) {
	err := isPasswordValid("pass")
	if err == nil {
		t.Errorf("Error while checking password: %s", err.Error())
	}
	err = isPasswordValid("password")
	if err == nil {
		t.Errorf("Error while validating password: %s", err.Error())
	}
	err = isPasswordValid("Password1")
	if err != nil {
		t.Errorf("Error while validating password: %s", err.Error())
	}
}

func TestIsEmailValid(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	err = isEmailValid("email")
	if err == nil {
		t.Errorf("Error while validating email: %s", err.Error())
	}
	//generates a string of 400 characters
	var email string
	for i := 0; i < 400; i++ {
		email += "a"
	}
	err = isEmailValid(email + "@gmail.com")
	if err == nil {
		t.Errorf("Error while validating email: %s", err.Error())
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs("email@email.com").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	err = isEmailValid("email@email.com")
	if err != nil {
		t.Errorf("Error while checking email: %s", err.Error())
	}
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto@toto.fr").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	err = isEmailValid("toto@toto.fr")
	if err == nil {
		t.Errorf("Error while checking email: %s", err.Error())
	}
}

func TestIsBirthdateValid(t *testing.T) {
	err := isBirthdateValid("20")
	if err == nil {
		t.Errorf("No error while validating invalid format birthdate")
	}
	err = isBirthdateValid("2020-01-01")
	if err == nil {
		t.Errorf("No error while validating too young birthdate")
	}
	err = isBirthdateValid("2042-01-01")
	if err == nil {
		t.Errorf("No error while validating too birthdate in the future")
	}
	err = isBirthdateValid("2000-01-01")
	if err != nil {
		t.Errorf("Error while validating birthday: %s", err.Error())
	}
}

func TestRegisterInvalidUsername(t *testing.T) {
	userInfos := models.RegisterInfos{Username: "us", Password: "password", Email: "email@email.fr", Birthdate: "2000-01-01"}

	_, _, err := Register(userInfos)
	if err == nil {
		t.Errorf("No error while creating user with invalid username")
	}
}

func TestRegisterInvalidPassword(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))

	userInfos := models.RegisterInfos{Username: "username", Password: "pass", Email: "email@email.fr", Birthdate: "2000-01-01"}

	_, _, err = Register(userInfos)
	if err == nil {
		t.Errorf("No error while creating user with invalid password")
	}
}

func TestRegisterInvalidEmail(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))

	userInfos := models.RegisterInfos{Username: "username", Password: "Password123@", Email: "email", Birthdate: "2000-01-01"}

	_, _, err = Register(userInfos)
	if err == nil {
		t.Errorf("No error while creating user with invalid email")
	}
}

func TestRegisterInvalidBirthdate(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectQuery("SELECT COUNT").WithArgs("email@email.fr").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))

	userInfos := models.RegisterInfos{Username: "username", Password: "Password123@", Email: "email@email.fr", Birthdate: "2025-01-01"}

	_, _, err = Register(userInfos)
	if err == nil {
		t.Errorf("No error received while creating user with invalid birthdate")
	}
}

func TestRegisterValid(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	os.Setenv("IS_TEST", "true")

	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectQuery("SELECT COUNT").WithArgs("email@email.fr").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectExec("INSERT INTO user").WithArgs("username", "username", "Password123@", "email@email.fr", "2000-01-01").WillReturnResult(sqlmock.NewResult(4, 1))

	userInfos := models.RegisterInfos{Username: "username", Password: "Password123@", Email: "email@email.fr", Birthdate: "2000-01-01"}

	id, returnStatus, err := Register(userInfos)
	if err != nil {
		t.Errorf("Error while creating user: %s", err.Error())
	}
	if returnStatus != http.StatusCreated {
		t.Errorf("Error: return status is %d instead of %d", returnStatus, http.StatusCreated)
	}
	if id != 4 {
		t.Errorf("Error: id is %d instead of %d", id, 4)
	}
	os.Setenv("IS_TEST", "false")
}

func TestRegisterOauthEmailInvalid(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}

	username := "username"
	email := "InvalidEmail"
	birthdate := "2000-01-01"

	mock.ExpectQuery("SELECT COUNT").WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectQuery("SELECT COUNT").WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))

	_, _, err = RegisterOauth(username, email, birthdate)
	if err == nil {
		t.Errorf("No error while creating user with invalid email")
	}
}

func TestRegisterOauthBirthdateInvalid(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}

	username := "username"
	email := "email@mail.com"
	birthdate := "2025-01-01"

	mock.ExpectQuery("SELECT COUNT").WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)) // In order to trigger username auto-generation
	mock.ExpectQuery("SELECT COUNT").WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))

	_, _, err = RegisterOauth(username, email, birthdate)
	if err == nil {
		t.Errorf("No error while creating user with invalid birthdate")
	}
}

func TestRegisterOauthErrorInserting(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}

	username := "username"
	email := "email@email.com"
	birthdate := "2000-01-01"

	mock.ExpectQuery("SELECT COUNT").WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectQuery("SELECT COUNT").WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectExec("INSERT INTO user").WithArgs(username, username, "", email, birthdate).WillReturnError(errors.New("Error while inserting user"))

	id, returnStatus, err := RegisterOauth(username, email, birthdate)
	if err == nil {
		t.Errorf("No error while creating user with invalid email, were expecting an error while inserting user")
	}
	if returnStatus != http.StatusInternalServerError {
		t.Errorf("Error: return status is %d instead of %d", returnStatus, http.StatusInternalServerError)
	}
	if id != 0 {
		t.Errorf("Error: id is %d instead of %d", id, 0)
	}
}

func TestRegisterOauthValid(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}

	username := "username"
	email := "email@email.com"
	birthdate := "2000-01-01"

	mock.ExpectQuery("SELECT COUNT").WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectQuery("SELECT COUNT").WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectExec("INSERT INTO user").WithArgs(username, username, "", email, birthdate).WillReturnResult(sqlmock.NewResult(4, 1))

	id, returnStatus, err := RegisterOauth(username, email, birthdate)
	if err != nil {
		t.Errorf("Error while creating user: %s", err.Error())
	}
	if returnStatus != http.StatusCreated {
		t.Errorf("Error: return status is %d instead of %d", returnStatus, http.StatusCreated)
	}
	if id != 4 {
		t.Errorf("Error: id is %d instead of %d", id, 4)
	}
}

func TestRegisterOauthValidLongUsername(t *testing.T) {
	os.Setenv("IS_TEST", "true")

	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}

	username := "User User User Toooo Looooooooong 166514 23132"
	email := "email@email.com"
	birthdate := "2000-01-01"

	mock.ExpectQuery("SELECT COUNT").WithArgs("user-1000").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectQuery("SELECT COUNT").WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectExec("INSERT INTO user").WithArgs("user-1000", "User User User Toooo", "", email, birthdate).WillReturnResult(sqlmock.NewResult(4, 1))

	id, returnStatus, err := RegisterOauth(username, email, birthdate)
	os.Setenv("IS_TEST", "false")
	if err != nil {
		t.Errorf("Error while creating user: %s", err.Error())
	}
	if returnStatus != http.StatusCreated {
		t.Errorf("Error: return status is %d instead of %d", returnStatus, http.StatusCreated)
	}
	if id != 4 {
		t.Errorf("Error: id is %d instead of %d", id, 4)
	}
}
