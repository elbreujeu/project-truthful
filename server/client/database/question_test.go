package database

import (
	"errors"
	"project_truthful/helpunittesting"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAddQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO question").WithArgs("question", "ip address", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	id, err := AddQuestion("question", 0, "ip address", true, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while adding question: %s", err.Error())
	}
	if id != 1 {
		t.Errorf("Id should be 1")
	}

	mock.ExpectExec("INSERT INTO question").WithArgs("question", 2, "ip address", true, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	id, err = AddQuestion("question", 2, "ip address", true, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while adding question: %s", err.Error())
	}
	if id != 1 {
		t.Errorf("Id should be 1")
	}

	mock.ExpectExec("INSERT INTO question").WithArgs("question", 3, "ip address", false, 1).WillReturnError(errors.New("error"))
	_, err = AddQuestion("question", 3, "ip address", false, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGetQuestionReceiverId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error
	mock.ExpectQuery("SELECT receiver_id FROM question").WithArgs(1).WillReturnError(errors.New("error for db test"))
	_, err = GetQuestionReceiverId(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with one row returned
	mock.ExpectQuery("SELECT receiver_id FROM question").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"author_id"}).AddRow(1))
	receiverId, err := GetQuestionReceiverId(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
	if receiverId != 1 {
		t.Errorf("Database error: expected 1, got %d", receiverId)
	}
}

func TestGetQuestions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock: %s", err.Error())
	}

	// test for an SQL error
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 30).WillReturnError(errors.New("error for db test"))
	_, err = GetQuestions(1, 0, 30, db)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	expectErr := mock.ExpectationsWereMet()
	if expectErr != nil {
		t.Error("Error while checking expectations")
	}

	// test for no rows returned
	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "creation_date"})
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 30).WillReturnRows(rows)
	questions, err := GetQuestions(1, 0, 30, db)
	if err != nil {
		t.Errorf("Error while getting questions: %s", err.Error())
	}
	if len(questions) != 0 {
		t.Errorf("Expected 0 questions, got %d", len(questions))
	}
	expectErr = mock.ExpectationsWereMet()
	if expectErr != nil {
		t.Error("Error while checking expectations")
	}
}

func TestGetQuestionsMultipleRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock: %s", err.Error())
	}
	// generate questions
	curTime := time.Now()
	questions := helpunittesting.GenerateTestQuestions(30, 1, curTime)

	// generate rows
	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "creation_date"})
	for _, question := range questions {
		rows.AddRow(question.Id, question.Text, question.AuthorId, question.IsAuthorAnonymous, question.ReceiverId, question.CreatedAt)
	}
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 30).WillReturnRows(rows)
	// expects all the queries for getting answers
	for _, question := range questions {
		mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(question.Id).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	}

	returnedQuestions, err := GetQuestions(1, 0, 30, db)
	if err != nil {
		t.Error("Expected nil, got", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
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

func TestGetQuestionsMultipleRowsAndAnswers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock: %s", err.Error())
	}
	// generate questions
	curTime := time.Now()
	questions := helpunittesting.GenerateTestQuestions(30, 1, curTime)

	// generate rows
	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "creation_date"})
	for _, question := range questions {
		rows.AddRow(question.Id, question.Text, question.AuthorId, question.IsAuthorAnonymous, question.ReceiverId, question.CreatedAt)
	}
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 30).WillReturnRows(rows)
	// expects all the queries for getting answers. Every query will return no answers except the 25th one
	for i, question := range questions {
		if i == 25 {
			mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(question.Id).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
		} else {
			mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(question.Id).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
		}
	}

	returnedQuestions, err := GetQuestions(1, 0, 30, db)
	if err != nil {
		t.Error("Expected nil, got", err)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
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
