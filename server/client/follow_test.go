package client

import (
	"database/sql"
	"errors"
	"net/http"
	"project_truthful/client/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFollowUser(t *testing.T) {
	//tests for followeeId != followerId
	code, err := FollowUser(1, 1)
	if code != http.StatusBadRequest {
		t.Errorf("Expected http.StatusBadRequest, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	//inits the sqlmock
	var mock sqlmock.Sqlmock
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error initializing mock database: %s", err)
	}

	//tests for followeeId error
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(2).WillReturnError(errors.New("error"))
	code, err = FollowUser(1, 2)
	if code != http.StatusInternalServerError {
		t.Errorf("Expected http.StatusInternalServerError, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	//tests for followeeId not found
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	code, err = FollowUser(1, 2)
	if code != http.StatusNotFound {
		t.Errorf("Expected http.StatusNotFound, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	//tests for followerId error
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnError(errors.New("error"))
	code, err = FollowUser(1, 2)
	if code != http.StatusInternalServerError {
		t.Errorf("Expected http.StatusInternalServerError, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	//tests for followerId not found
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	code, err = FollowUser(1, 2)
	if code != http.StatusNotFound {
		t.Errorf("Expected http.StatusNotFound, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	//tests for check follow already exists error
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnError(errors.New("error"))
	code, err = FollowUser(1, 2)
	if code != http.StatusInternalServerError {
		t.Errorf("Expected http.StatusInternalServerError, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	//tests for follow already exists
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	code, err = FollowUser(1, 2)
	if code != http.StatusBadRequest {
		t.Errorf("Expected http.StatusBadRequest, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	//tests for follow insert error
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO follow").WillReturnError(errors.New("error"))
	code, err = FollowUser(1, 2)
	if code != http.StatusInternalServerError {
		t.Errorf("Expected http.StatusInternalServerError, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	//tests for follow insert success
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO follow").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	code, err = FollowUser(1, 2)
	if code != http.StatusOK {
		t.Errorf("Expected http.StatusOK, got %d", code)
	}
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
}

func TestUnfollowUser(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error initializing mock database: %s", err)
	}

	// test for error when checking if follow exists
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnError(errors.New("error"))
	code, err := UnfollowUser(1, 2)
	if code != http.StatusInternalServerError {
		t.Errorf("Expected http.StatusInternalServerError, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	// test for follow does not exists
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	code, err = UnfollowUser(1, 2)
	if code != http.StatusBadRequest {
		t.Errorf("Expected http.StatusBadRequest, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	// test for error when deleting follow
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectExec("DELETE FROM follow").WithArgs(1, 2).WillReturnError(errors.New("error"))
	code, err = UnfollowUser(1, 2)
	if code != http.StatusInternalServerError {
		t.Errorf("Expected http.StatusInternalServerError, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	// test for success when deleting follow
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectExec("DELETE FROM follow").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	code, err = UnfollowUser(1, 2)
	if code != http.StatusOK {
		t.Errorf("Expected http.StatusOK, got %d", code)
	}
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
}

func TestGetFollowers(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error initializing mock database: %s", err)
	}

	// test for user not found
	mock.ExpectQuery("SELECT (.+) FROM user").WithArgs("username").WillReturnError(sql.ErrNoRows)
	followers, code, err := GetFollowers("username", 10, 0)
	if code != http.StatusNotFound {
		t.Errorf("Expected http.StatusNotFound, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if len(followers) != 0 {
		t.Errorf("Expected empty followers list, got %d followers", len(followers))
	}

	// test for error when getting user ID
	mock.ExpectQuery("SELECT (.+) FROM user").WithArgs("username").WillReturnError(errors.New("error"))
	followers, code, err = GetFollowers("username", 10, 0)
	if code != http.StatusInternalServerError {
		t.Errorf("Expected http.StatusInternalServerError, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if len(followers) != 0 {
		t.Errorf("Expected empty followers list, got %d followers", len(followers))
	}

	// test for error when getting followers
	mock.ExpectQuery("SELECT (.+) FROM user").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT (.+) FROM follow").WithArgs(1, 10, 0).WillReturnError(errors.New("error"))
	followers, code, err = GetFollowers("username", 10, 0)
	if code != http.StatusInternalServerError {
		t.Errorf("Expected http.StatusInternalServerError, got %d", code)
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if len(followers) != 0 {
		t.Errorf("Expected empty followers list, got %d followers", len(followers))
	}

	// test for success
	mock.ExpectQuery("SELECT (.+) FROM user").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT (.+) FROM follow").WithArgs(1, 10, 0).WillReturnRows(sqlmock.NewRows([]string{"follower"}).AddRow(2).AddRow(3))
	mock.ExpectQuery("SELECT (.+) FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("user1", "User 1"))
	mock.ExpectQuery("SELECT (.+) FROM user").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("user2", "User 2"))
	followers, code, err = GetFollowers("username", 10, 0)
	if code != http.StatusOK {
		t.Errorf("Expected http.StatusOK, got %d", code)
	}
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if len(followers) != 2 {
		t.Errorf("Expected 2 followers, got %d followers", len(followers))
	}
	if followers[0].Id != 2 || followers[0].Username != "user1" || followers[0].DisplayName != "User 1" {
		t.Errorf("Expected follower 1 to have ID 2, Username 'user1', and DisplayName 'User 1'")
	}
	if followers[1].Id != 3 || followers[1].Username != "user2" || followers[1].DisplayName != "User 2" {
		t.Errorf("Expected follower 2 to have ID 3, Username 'user2', and DisplayName 'User 2'")
	}
}
