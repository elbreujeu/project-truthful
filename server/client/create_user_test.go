package client

import (
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

func TestCreateUserInvalidUsername(t *testing.T) {
	userInfos := models.CreateUserInfos{Username: "us", Password: "password", Email: "email@email.fr", Birthdate: "2000-01-01"}

	_, _, err := CreateUser(userInfos)
	if err == nil {
		t.Errorf("No error while creating user with invalid username")
	}
}

func TestCreateUserInvalidPassword(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))

	userInfos := models.CreateUserInfos{Username: "username", Password: "pass", Email: "email@email.fr", Birthdate: "2000-01-01"}

	_, _, err = CreateUser(userInfos)
	if err == nil {
		t.Errorf("No error while creating user with invalid password")
	}
}

func TestCreateUserInvalidEmail(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))

	userInfos := models.CreateUserInfos{Username: "username", Password: "Password123@", Email: "email", Birthdate: "2000-01-01"}

	_, _, err = CreateUser(userInfos)
	if err == nil {
		t.Errorf("No error while creating user with invalid email")
	}
}

func TestCreateUserInvalidBirthdate(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	mock.ExpectQuery("SELECT COUNT").WithArgs("email@email.fr").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))

	userInfos := models.CreateUserInfos{Username: "username", Password: "Password123@", Email: "email@email.fr", Birthdate: "2025-01-01"}

	_, _, err = CreateUser(userInfos)
	if err == nil {
		t.Errorf("No error received while creating user with invalid birthdate")
	}
}

func TestCreateUserValid(t *testing.T) {
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

	userInfos := models.CreateUserInfos{Username: "username", Password: "Password123@", Email: "email@email.fr", Birthdate: "2000-01-01"}

	id, returnStatus, err := CreateUser(userInfos)
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
