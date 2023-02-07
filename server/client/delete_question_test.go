package client

import (
	"database/sql"
	"errors"
	"project_truthful/client/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMarkQuestionAsDeleted(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}

	// check with question id not found
	mock.ExpectQuery("SELECT receiver_id FROM question WHERE id").WithArgs(1).WillReturnError(sql.ErrNoRows)
	_, err = MarkQuestionAsDeleted(1, 1)
	if err == nil || err.Error() != "question not found" {
		t.Errorf("Expected error, got nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// check with error while checking answer id
	mock.ExpectQuery("SELECT receiver_id FROM question WHERE id").WithArgs(2).WillReturnError(errors.New("test error"))
	_, err = MarkQuestionAsDeleted(2, 2)
	if err == nil || err.Error() != "test error" {
		t.Errorf("Expected error, got nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// check with answer id found but user id not matching
	mock.ExpectQuery("SELECT receiver_id FROM question WHERE id").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(3))
	_, err = MarkQuestionAsDeleted(2, 3)
	if err == nil || err.Error() != "user is not the receiver of the question" {
		t.Errorf("Expected error, got nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// check with error while checking if question has been answered
	mock.ExpectQuery("SELECT receiver_id FROM question WHERE id").WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(4))
	mock.ExpectQuery("SELECT id FROM answer").WithArgs(4).WillReturnError(errors.New("test error"))
	_, err = MarkQuestionAsDeleted(4, 4)
	if err == nil || err.Error() != "test error" {
		t.Errorf("Expected error, got nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// check with question has been answered but error while deleting answer
	mock.ExpectQuery("SELECT receiver_id FROM question WHERE id").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(5))
	mock.ExpectQuery("SELECT id FROM answer").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))
	mock.ExpectExec("UPDATE answer SET has_been_deleted = 1").WithArgs(5).WillReturnError(errors.New("test error"))
	_, err = MarkQuestionAsDeleted(5, 5)
	if err == nil || err.Error() != "test error" {
		t.Errorf("Expected error, got nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// check with question has no answer but error while deleting question
	mock.ExpectQuery("SELECT receiver_id FROM question WHERE id").WithArgs(6).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(6))
	mock.ExpectQuery("SELECT id FROM answer").WithArgs(6).WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("UPDATE question SET has_been_deleted = 1").WithArgs(6).WillReturnError(errors.New("test error"))
	_, err = MarkQuestionAsDeleted(6, 6)
	if err == nil || err.Error() != "test error" {
		t.Errorf("Expected error, got nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// check with question has no answer and question deleted successfully
	mock.ExpectQuery("SELECT receiver_id FROM question WHERE id").WithArgs(7).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(7))
	mock.ExpectQuery("SELECT id FROM answer").WithArgs(7).WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("UPDATE question SET has_been_deleted = 1").WithArgs(7).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = MarkQuestionAsDeleted(7, 7)
	if err != nil {
		t.Errorf("Expected nil, got error: %s", err.Error())
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
}
