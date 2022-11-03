package client

import (
	"errors"
	"project_truthful/client/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLike(t *testing.T) {
	//inits the sqlmock
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer database.DB.Close()

	// Test that the like function returns an error when the user does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	_, err = LikeAnswer(1, 1)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Expected error when user does not exist")
	}

	// Test that the like function returns an error when getting an error looking for user
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnError(errors.New("error"))
	_, err = LikeAnswer(1, 1)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// Test that the like function returns an error when getting an error when checking if the user id exists
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnError(errors.New("error"))
	_, err = LikeAnswer(1, 1)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// Test that the like function returns an error when the post id returns an error
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnError(errors.New("error"))
	_, err = LikeAnswer(1, 1)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// Test that the like function returns an error when the post id is not found
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
	_, err = LikeAnswer(1, 2)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// Test that the like function returns an error when the user already likes the post
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err = LikeAnswer(1, 2)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// Test with an error when checking if like exists
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 2).WillReturnError(errors.New("error"))
	_, err = LikeAnswer(1, 2)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// Test with an error when inserting like
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
	mock.ExpectExec("INSERT INTO").WithArgs(1, 2).WillReturnError(errors.New("error"))
	_, err = LikeAnswer(1, 2)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// Test that the like function returns no error when the user likes the post
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
	mock.ExpectExec("INSERT INTO answer_like").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = LikeAnswer(1, 2)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error should be nil")
	}
}

func TestUnlikeAnswer(t *testing.T) {
	//inits the sqlmock
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer database.DB.Close()

	// Test that the remove like function returns an error when getting an error when checking if the like exists
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 2).WillReturnError(errors.New("error"))
	_, err = UnlikeAnswer(1, 2)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// Test that the remove like function returns an error when the like does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
	_, err = UnlikeAnswer(1, 2)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// Test that the remove like function returns an error when getting an error when deleting the like
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("DELETE FROM").WithArgs(1, 2).WillReturnError(errors.New("error"))
	_, err = UnlikeAnswer(1, 2)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// Test that the remove like function returns no error when the like is removed
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("DELETE FROM").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = UnlikeAnswer(1, 2)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error should be nil")
	}
}
