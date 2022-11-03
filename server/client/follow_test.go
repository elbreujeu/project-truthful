package client

import (
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
