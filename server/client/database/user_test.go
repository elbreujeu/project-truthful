package database

import (
	"database/sql"
	"errors"
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

	mock.ExpectExec("INSERT INTO user").WithArgs("username_error", "username_error", "password", "email", "birthdate").WillReturnError(errors.New("error"))
	_, err = InsertUser("username_error", "password", "email", "birthdate", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestCheckUsernameExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckUsernameExists("username", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if username exists: %s", err.Error())
	}
	if !exists {
		t.Errorf("Username should exist")
	}
	// tests that the username does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckUsernameExists("toto", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if username exists: %s", err.Error())
	}
	if exists {
		t.Errorf("Username should not exist")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs("username_error").WillReturnError(errors.New("error"))
	_, err = CheckUsernameExists("username_error", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGetUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT id").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	id, err := GetUserId("username", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while getting user id: %s", err.Error())
	}
	if id != 1 {
		t.Errorf("id should be 1, but is %d", id)
	}

	mock.ExpectQuery("SELECT id").WithArgs("username_not_existing").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
	id, err = GetUserId("username_not_existing", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while getting user id: %s", err.Error())
	}
	if id != 0 {
		t.Errorf("id should be 0, but is %d", id)
	}

	mock.ExpectQuery("SELECT id").WithArgs("username_error").WillReturnError(errors.New("error"))
	_, err = GetUserId("username_error", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestCheckEmailExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs("email@email.fr").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckEmailExists("email@email.fr", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if email exists: %s", err.Error())
	}
	if !exists {
		t.Errorf("Email should exist")
	}
	// tests that the email does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto@toto.fr").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckEmailExists("toto@toto.fr", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if email exists: %s", err.Error())
	}
	if exists {
		t.Errorf("Email should not exist")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs("unexistant@toto.fr").WillReturnError(errors.New("error"))
	_, err = CheckEmailExists("unexistant@toto.fr", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGetHashedPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT password").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("password"))
	password, err := GetHashedPassword(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while getting hashed password: %s", err.Error())
	}
	if password != "password" {
		t.Errorf("password should be 'password', but is %s", password)
	}

	mock.ExpectQuery("SELECT password").WithArgs(3).WillReturnError(errors.New("error"))
	_, err = GetHashedPassword(3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestCheckUserIdExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckUserIdExists(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if user id exists: %s", err.Error())
	}
	if !exists {
		t.Errorf("User id should exist")
	}
	// tests that the user id does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckUserIdExists(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if user id exists: %s", err.Error())
	}
	if exists {
		t.Errorf("User id should not exist")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(3).WillReturnError(errors.New("error"))
	_, err = CheckUserIdExists(3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGetUsernameAndDisplayName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"username", "display_name"}).AddRow("john", "John Doe")
	mock.ExpectQuery("SELECT username, display_name FROM user WHERE id = \\?").WithArgs(1).WillReturnRows(rows)

	username, displayName, err := GetUsernameAndDisplayName(1, db)
	if err != nil {
		t.Errorf("error was not expected while getting username and display name: %s", err)
	}
	if username != "john" {
		t.Errorf("unexpected username: %s", username)
	}
	if displayName != "John Doe" {
		t.Errorf("unexpected display name: %s", displayName)
	}
}

func TestGetUsernameAndDisplayNameError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT username, display_name FROM user WHERE id = \\?").WithArgs(1).WillReturnError(sql.ErrNoRows)

	_, _, err = GetUsernameAndDisplayName(1, db)
	if err == nil {
		t.Errorf("error was expected while getting username and display name, but got nil")
	}
}

func TestUpdateUserInformations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := 1
	displayName := "John Doe"
	email := "john.doe@example.com"

	mock.ExpectExec("UPDATE user SET display_name = \\?, email = \\? WHERE id = \\?").
		WithArgs(displayName, email, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = UpdateUserInformations(id, displayName, email, db)
	if err != nil {
		t.Errorf("error was not expected while updating user informations: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUserInformationsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := 1
	displayName := "John Doe"
	email := "john.doe@example.com"

	mock.ExpectExec("UPDATE user SET display_name = \\?, email = \\? WHERE id = \\?").
		WithArgs(displayName, email, id).
		WillReturnError(errors.New("database error"))

	err = UpdateUserInformations(id, displayName, email, db)
	if err == nil {
		t.Errorf("error was expected while updating user informations, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetOAuthProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT id FROM oauth_provider WHERE name = \\?").WithArgs("google").WillReturnRows(rows)

	id, err := GetOAuthProvider("google", db)
	if err != nil {
		t.Errorf("error was not expected while getting oauth provider: %s", err)
	}
	if id != 1 {
		t.Errorf("unexpected provider id: %d", id)
	}
}

func TestGetOAuthProviderNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM oauth_provider WHERE name = \\?").WithArgs("google").WillReturnError(sql.ErrNoRows)

	_, err = GetOAuthProvider("google", db)
	if err == nil {
		t.Errorf("error was expected while getting oauth provider, but got nil")
	}
}

func TestGetOAuthProvierError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM oauth_provider WHERE name = \\?").WithArgs("google").WillReturnError(errors.New("database error"))

	_, err = GetOAuthProvider("google", db)
	if err == nil {
		t.Errorf("error was expected while getting oauth provider, but got nil")
	}
}
func TestGetUserIdBySubject(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	providerID := 1
	subject := "subject123"

	mock.ExpectQuery("SELECT user_id FROM oauth_login WHERE oauth_provider_id = \\? AND subject_id = \\?").
		WithArgs(providerID, subject).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

	userID, err := GetUserIdBySubject(providerID, subject, db)
	if err != nil {
		t.Errorf("error was not expected while getting user id for subject %s: %s", subject, err)
	}
	if userID != 1 {
		t.Errorf("unexpected user id: %d", userID)
	}
}

func TestGetUserIdBySubjectError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	providerID := 1
	subject := "subject123"

	mock.ExpectQuery("SELECT user_id FROM oauth_login WHERE oauth_provider_id = \\? AND subject_id = \\?").
		WithArgs(providerID, subject).
		WillReturnError(sql.ErrNoRows)

	_, err = GetUserIdBySubject(providerID, subject, db)
	if err == nil {
		t.Errorf("error was expected while getting user id for subject %s, but got nil", subject)
	}
}
