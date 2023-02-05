package database

import (
	"database/sql"
	"errors"
	"fmt"
	"project_truthful/helpunittesting"
	"project_truthful/models"
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
	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "created_at"})
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
	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "created_at"})
	for _, question := range questions {
		rows.AddRow(question.Id, question.Text, question.Author.Id, question.IsAuthorAnonymous, question.ReceiverId, question.CreatedAt)
	}
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 30).WillReturnRows(rows)
	// expects all the queries for getting answers
	for i, question := range questions {
		mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(question.Id).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
		mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(question.Author.Id).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username"+fmt.Sprintf("%d", i), "display_name"+fmt.Sprintf("%d", i)))
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
		if q.Author.Id != questions[i].Author.Id {
			t.Errorf("Expected AuthorId to be %d, but got %d", questions[i].Author.Id, q.Author.Id)
		}
		if q.Author.Username != "username"+fmt.Sprintf("%d", i) {
			t.Errorf("Expected AuthorUsername to be %s, but got %s", "username"+fmt.Sprintf("%d", i), q.Author.Username)
		}
		if q.Author.DisplayName != "display_name"+fmt.Sprintf("%d", i) {
			t.Errorf("Expected AuthorDisplayName to be %s, but got %s", "display_name"+fmt.Sprintf("%d", i), q.Author.DisplayName)
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
	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "created_at"})
	for _, question := range questions {
		rows.AddRow(question.Id, question.Text, question.Author.Id, question.IsAuthorAnonymous, question.ReceiverId, question.CreatedAt)
	}
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 30).WillReturnRows(rows)
	// expects all the queries for getting answers. Every query will return no answers except the 25th one
	for i, question := range questions {
		if i == 25 {
			mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(question.Id).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
		} else {
			mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(question.Id).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
			mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(question.Author.Id).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username"+fmt.Sprintf("%d", i), "display_name"+fmt.Sprintf("%d", i)))
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
		if q.Author.Id != questions[i].Author.Id {
			t.Errorf("Expected AuthorId to be %d, but got %d", questions[i].Author.Id, q.Author.Id)
		}
		if q.Author.Username != "username"+fmt.Sprintf("%d", i) {
			t.Errorf("Expected AuthorUsername to be %s, but got %s", "username"+fmt.Sprintf("%d", i), q.Author.Username)
		}
		if q.Author.DisplayName != "display_name"+fmt.Sprintf("%d", i) {
			t.Errorf("Expected AuthorDisplayName to be %s, but got %s", "display_name"+fmt.Sprintf("%d", i), q.Author.DisplayName)
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

func TestGetQuestionById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	question := models.Question{
		Id:                1,
		Text:              "What is the meaning of life?",
		Author:            models.UserPreview{Id: 42, Username: "username", DisplayName: "display_name"},
		IsAuthorAnonymous: false,
		ReceiverId:        12,
		CreatedAt:         time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "created_at"}).
		AddRow(question.Id, question.Text, question.Author.Id, question.IsAuthorAnonymous, question.ReceiverId, question.CreatedAt)

	mock.ExpectQuery("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(question.Author.Id).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow(question.Author.Username, question.Author.DisplayName))

	gotQuestion, err := GetQuestionById(1, db)
	if err != nil {
		t.Fatalf("error was not expected while retrieving question: %s", err)
	}
	if gotQuestion.Id != question.Id ||
		gotQuestion.Text != question.Text ||
		gotQuestion.Author.Id != question.Author.Id ||
		gotQuestion.Author.Username != question.Author.Username ||
		gotQuestion.Author.DisplayName != question.Author.DisplayName ||
		gotQuestion.IsAuthorAnonymous != question.IsAuthorAnonymous ||
		gotQuestion.ReceiverId != question.ReceiverId ||
		gotQuestion.CreatedAt != question.CreatedAt {
		t.Errorf("question was not as expected.\nExpected: %v\nGot: %v\n", question, gotQuestion)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetQuestionByIdNullAuthor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	question := models.Question{
		Id:                1,
		Text:              "What is the meaning of life?",
		Author:            models.UserPreview{Id: 0, Username: "", DisplayName: ""},
		IsAuthorAnonymous: true,
		ReceiverId:        13,
		CreatedAt:         time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "created_at"}).
		AddRow(question.Id, question.Text, nil, question.IsAuthorAnonymous, question.ReceiverId, question.CreatedAt)

	mock.ExpectQuery("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	gotQuestion, err := GetQuestionById(1, db)
	if err != nil {
		t.Fatalf("error was not expected while retrieving question: %s", err)
	}
	if gotQuestion.Id != question.Id ||
		gotQuestion.Text != question.Text ||
		gotQuestion.Author.Id != question.Author.Id ||
		gotQuestion.Author.Username != question.Author.Username ||
		gotQuestion.Author.DisplayName != question.Author.DisplayName ||
		gotQuestion.IsAuthorAnonymous != question.IsAuthorAnonymous ||
		gotQuestion.ReceiverId != question.ReceiverId ||
		gotQuestion.CreatedAt != question.CreatedAt {
		t.Errorf("question was not as expected.\nExpected: %v\nGot: %v\n", question, gotQuestion)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetQuestionByIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question WHERE id = \\?").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err = GetQuestionById(1, db)
	if err == nil {
		t.Errorf("error was expected but got nil")
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
