package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
	"project_truthful/helpunittesting"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetQuestionsFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	database.DB = db
	if err != nil {
		t.Fatalf("Error while creating mock: %s", err.Error())
	}

	// test with not existing user id
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	_, status, err := GetQuestions(1, 0, 30)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if status != http.StatusNotFound {
		t.Error("Expected status 404, got", status)
	}

	// test with error while checking user id + too low count and start
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnError(errors.New("error while checking user id"))
	_, status, err = GetQuestions(1, -1, -1)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if status != http.StatusInternalServerError {
		t.Error("Expected status 500, got", status)
	}

	// test with existing user id but error while getting questions + too high count
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 30).WillReturnError(errors.New("error while getting questions"))
	_, status, err = GetQuestions(1, 0, 50)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if status != http.StatusInternalServerError {
		t.Error("Expected status 500, got", status)
	}
}

func TestGetQuestionsNoQuestions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock: %s", err.Error())
	}
	database.DB = db
	// test with existing user id and nil questions
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 30).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "text", "created_at", "updated_at"}))
	questions, status, err := GetQuestions(1, 0, 30)
	if err != nil {
		t.Error("Expected nil, got", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
	if len(questions) != 0 {
		t.Errorf("Expected number of questions to be %d, but got %d", 0, len(questions))
	}
}

func TestGetQuestionsNoAnswers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock: %s", err.Error())
	}
	database.DB = db

	// test with existing user id and questions, all not answered
	// generate questions
	curTime := time.Now()
	questions := helpunittesting.GenerateTestQuestions(30, 1, curTime)

	// generate rows
	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "creation_date"})
	for _, question := range questions {
		rows.AddRow(question.Id, question.Text, question.AuthorId, question.IsAuthorAnonymous, question.ReceiverId, question.CreatedAt)
	}
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 30).WillReturnRows(rows)
	// expects all the queries for getting answers
	for _, question := range questions {
		mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(question.Id).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	}

	returnedQuestions, status, err := GetQuestions(1, 0, 30)
	if err != nil {
		t.Error("Expected nil, got", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// checks if the questions are correctly returned
	for i, q := range returnedQuestions {
		if q.Id != questions[i].Id {
			t.Errorf("Expected Id to be %d, but got %d", questions[i].Id, q.Id)
		}
		if q.Text != questions[i].Text {
			t.Errorf("Expected Text to be %s, but got %s", questions[i].Text, q.Text)
		}
		if q.AuthorId != questions[i].AuthorId {
			t.Errorf("Expected AuthorId to be %d, but got %d", questions[i].AuthorId, q.AuthorId)
		}
		if q.IsAuthorAnonymous != questions[i].IsAuthorAnonymous {
			t.Errorf("Expected IsAuthorAnonymous to be %t, but got %t", questions[i].IsAuthorAnonymous, q.IsAuthorAnonymous)
		}
		if q.ReceiverId != questions[i].ReceiverId {
			t.Errorf("Expected ReceiverId to be %d, but got %d", questions[i].ReceiverId, q.ReceiverId)
		}
		if q.CreatedAt != questions[i].CreatedAt {
			t.Errorf("Expected CreatedAt to be %v, but got %v", questions[i].CreatedAt, q.CreatedAt)
		}
	}
}

func TestGetQuestionsWithAnswers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock: %s", err.Error())
	}
	database.DB = db

	// test with existing user id and questions, all not answered
	// generate questions
	curTime := time.Now()
	questions := helpunittesting.GenerateTestQuestions(30, 1, curTime)

	// generate rows
	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "creation_date"})
	for _, question := range questions {
		rows.AddRow(question.Id, question.Text, question.AuthorId, question.IsAuthorAnonymous, question.ReceiverId, question.CreatedAt)
	}
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 30).WillReturnRows(rows)
	// expects all the queries for getting answers. Every query will return no answers except the 25th one
	for i, question := range questions {
		if i == 25 {
			mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(question.Id).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
		} else {
			mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(question.Id).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
		}
	}

	returnedQuestions, status, err := GetQuestions(1, 0, 30)
	if err != nil {
		t.Error("Expected nil, got", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// checks if the questions are correctly returned until the 25th one
	for i, q := range returnedQuestions {
		if i >= 25 {
			i += 1
		}
		if q.Id != questions[i].Id {
			t.Errorf("Expected Id to be %d, but got %d", questions[i].Id, q.Id)
		}
		if q.Text != questions[i].Text {
			t.Errorf("Expected Text to be %s, but got %s", questions[i].Text, q.Text)
		}
		if q.AuthorId != questions[i].AuthorId {
			t.Errorf("Expected AuthorId to be %d, but got %d", questions[i].AuthorId, q.AuthorId)
		}
		if q.IsAuthorAnonymous != questions[i].IsAuthorAnonymous {
			t.Errorf("Expected IsAuthorAnonymous to be %t, but got %t", questions[i].IsAuthorAnonymous, q.IsAuthorAnonymous)
		}
		if q.ReceiverId != questions[i].ReceiverId {
			t.Errorf("Expected ReceiverId to be %d, but got %d", questions[i].ReceiverId, q.ReceiverId)
		}
		if q.CreatedAt != questions[i].CreatedAt {
			t.Errorf("Expected CreatedAt to be %v, but got %v", questions[i].CreatedAt, q.CreatedAt)
		}
	}
}
