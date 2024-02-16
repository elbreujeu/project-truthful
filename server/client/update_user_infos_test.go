package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCheckUserInfos(t *testing.T) {
	// Test case 1: Empty display name
	err := checkUserInfos("", "test@example.com")
	if err == nil || err.Error() != "display name is empty" {
		t.Errorf("Expected error: display name is empty, got: %v", err)
	}

	// Test case 2: Display name is too long
	err = checkUserInfos("ThisIsAVeryLongDisplayNameThatExceedsTheMaximumLength", "test@example.com")
	if err == nil || err.Error() != "display name is too long" {
		t.Errorf("Expected error: display name is too long, got: %v", err)
	}

	// Test case 3: Empty email address
	err = checkUserInfos("John Doe", "")
	if err == nil || err.Error() != "email address is empty" {
		t.Errorf("Expected error: email address is empty, got: %v", err)
	}

	// Test case 4: Email address is too long
	err = checkUserInfos("John Doe", "test@example.com"+generateLongString(350))
	if err == nil || err.Error() != "email address is too long" {
		t.Errorf("Expected error: email address is too long, got: %v", err)
	}

	// Test case 5: Invalid email address
	err = checkUserInfos("John Doe", "invalid_email")
	if err == nil || err.Error() != "mail: missing '@' or angle-addr" {
		t.Errorf("Expected error: mail: missing '@' or angle-addr, got: %v", err)
	}

	// Test case 6: Valid input
	err = checkUserInfos("John Doe", "test@example.com")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Helper function to generate a long string
func generateLongString(length int) string {
	str := ""
	for i := 0; i < length; i++ {
		str += "a"
	}
	return str
}

// func UpdateUserInformations(requesterId int, displayName string, email string) (int, error) {
// 	exists, err := database.CheckUserIdExists(requesterId, database.DB)
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}
// 	if !exists {
// 		return http.StatusNotFound, errors.New("user not found")
// 	}

// 	err = checkUserInfos(displayName, email)
// 	if err != nil {
// 		return http.StatusBadRequest, err
// 	}

// 	err = database.UpdateUserInformations(requesterId, displayName, email, database.DB)
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}
// 	return 0, nil
// }

func TestUpdateUserInternalServerError(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnError(errors.New("sql error"))
	code, err := UpdateUserInformations(1, "John Doe", "ffdqsjfsd@gmail.com")
	if code != http.StatusInternalServerError || err.Error() != "sql error" {
		t.Errorf("Expected error: sql error, got: %v", err)
	}
}

func TestUpdateUserUserNotFound(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	code, err := UpdateUserInformations(1, "John Doe", "toto@gmail.com")
	if code != http.StatusNotFound || err.Error() != "user not found" {
		t.Errorf("Expected error: user not found, got: %v", err)
	}
}

func TestUpdateUserInvalidDisplayName(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	code, err := UpdateUserInformations(1, "", "toto@toto.fr")
	if code != http.StatusBadRequest || err.Error() != "display name is empty" {
		t.Errorf("Expected error: display name is empty, got: %v", err)
	}
}

func TestUpdateUserDbError(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectExec("UPDATE user SET display_name = \\?, email = \\? WHERE id = \\?").
		WithArgs("John Doe", "toto@toto.fr", 1).
		WillReturnError(errors.New("sql error"))
	code, err := UpdateUserInformations(1, "John Doe", "toto@toto.fr")
	if code != http.StatusInternalServerError || err.Error() != "sql error" {
		t.Errorf("Expected error: sql error, got: %v", err)
	}
}

func TestUpdateUserValid(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock: %s", err.Error())
	}
	defer database.DB.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectExec("UPDATE user SET display_name = \\?, email = \\? WHERE id = \\?").
		WithArgs("John Doe", "toto@toto.fr", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	code, err := UpdateUserInformations(1, "John Doe", "toto@toto.fr")
	if code != 0 || err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
