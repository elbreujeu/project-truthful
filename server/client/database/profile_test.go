package database

import (
	"errors"
	"project_truthful/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserProfileInfos(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error on query to get name and display name
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(1).WillReturnError(errors.New("error for db test"))
	_, err = GetUserProfileInfos(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with an error when getting follower count
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username", "display_name"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(2).WillReturnError(errors.New("error for db test"))
	_, err = GetUserProfileInfos(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with an error when getting following count
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username", "display_name"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(3).WillReturnError(errors.New("error for db test"))
	_, err = GetUserProfileInfos(3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with an error when getting answer count
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username", "display_name"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(4).WillReturnError(errors.New("error for db test"))
	_, err = GetUserProfileInfos(4, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with an error when getting answers
	// todo, function is not coded yet

	// Test with no error
	var expected models.UserProfileInfos
	expected.Username = "username"
	expected.DisplayName = "display_name"
	expected.FollowerCount = 1
	expected.FollowingCount = 1
	expected.AnswerCount = 1
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username", "display_name"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	profile, err := GetUserProfileInfos(5, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
	// compares each field of the struct
	if profile.Username != expected.Username {
		t.Errorf("Database error: expected %s, got %s", expected.Username, profile.Username)
	}
	if profile.DisplayName != expected.DisplayName {
		t.Errorf("Database error: expected %s, got %s", expected.DisplayName, profile.DisplayName)
	}
	if profile.FollowerCount != expected.FollowerCount {
		t.Errorf("Database error: expected %d, got %d", expected.FollowerCount, profile.FollowerCount)
	}
	if profile.FollowingCount != expected.FollowingCount {
		t.Errorf("Database error: expected %d, got %d", expected.FollowingCount, profile.FollowingCount)
	}
	if profile.AnswerCount != expected.AnswerCount {
		t.Errorf("Database error: expected %d, got %d", expected.AnswerCount, profile.AnswerCount)
	}
}
