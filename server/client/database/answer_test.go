package database

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCheckAnswerIdExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnError(errors.New("error for db test"))
	_, err = CheckAnswerIdExists(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with one row returned
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckAnswerIdExists(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
	if !exists {
		t.Errorf("Database error: expected true, got false")
	}

	// Test with no row returned
	mock.ExpectQuery("SELECT COUNT").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckAnswerIdExists(3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
	if exists {
		t.Errorf("Database error: expected false, got true")
	}
}

func TestAddAnswer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with query fail
	mock.ExpectExec("INSERT INTO answer").WithArgs(1, 1, "content", "ip_address").WillReturnError(errors.New("error for db test"))
	_, err = AddAnswer(1, 1, "content", "ip_address", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with no error
	mock.ExpectExec("INSERT INTO answer").WithArgs(2, 2, "content", "ip_address").WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = AddAnswer(2, 2, "content", "ip_address", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
}

func TestHasQuestionBeenAnswered(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := HasQuestionBeenAnswered(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if question has been answered: %s", err.Error())
	}
	if !exists {
		t.Errorf("Question should have been answered")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = HasQuestionBeenAnswered(3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if question has been answered: %s", err.Error())
	}
	if exists {
		t.Errorf("Question should not have been answered")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(4).WillReturnError(errors.New("error"))
	_, err = HasQuestionBeenAnswered(4, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGetAnswerAuthorId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectQuery("SELECT user_id FROM answer").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"author_id"}).AddRow(1))
	authorId, err := GetAnswerAuthorId(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while getting answer author id: %s", err.Error())
	}
	if authorId != 1 {
		t.Errorf("Author id should be 1")
	}

	mock.ExpectQuery("SELECT user_id FROM answer").WithArgs(2).WillReturnError(errors.New("error"))
	_, err = GetAnswerAuthorId(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestMarkAnswerAsDeleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectExec("UPDATE answer").WithArgs(1).WillReturnError(errors.New("error"))
	err = MarkAnswerAsDeleted(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	mock.ExpectExec("UPDATE answer").WithArgs(2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = MarkAnswerAsDeleted(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error should be nil")
	}
}

func TestGetAnswersErrors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// test with error on first query
	mock.ExpectQuery("SELECT id, question_id, text, created_at FROM answer").WithArgs(1, 0, 30).WillReturnError(errors.New("error for test"))
	_, err = getAnswers(1, 0, 30, 0, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// test with wrong type
	mock.ExpectQuery("SELECT id, question_id, text, created_at FROM answer").WithArgs(1, 0, 30).WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "text", "created_at"}).AddRow(1, "error", "text", time.Now()))
	_, err = getAnswers(1, 0, 30, 0, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	// test with error on GetQuestionById query and GetLikeCountForAnswer query
	mock.ExpectQuery("SELECT id, question_id, text, created_at FROM answer").WithArgs(1, 0, 30).WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "text", "created_at"}).AddRow(1, 1, "text", time.Now()))
	mock.ExpectQuery("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question").WithArgs(1).WillReturnError(errors.New("error for test"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1).WillReturnError(errors.New("error for test"))
	answers, err := getAnswers(1, 0, 30, 0, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error should be nil")
	}
	if len(answers) != 1 {
		t.Errorf("Length of answers should be 1")
	}
	if answers[0].Author.Id != 0 {
		t.Errorf("Author id should be 0")
	}
	if answers[0].LikeCount != 0 {
		t.Errorf("Like count should be 0")
	}
}

func TestGetAnswers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, question_id, text, created_at FROM answer").WithArgs(1, 0, 30).WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "text", "created_at"}).AddRow(1, 1, "text", time.Now()))
	mock.ExpectQuery("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "created_at"}).AddRow(1, "text", 0, true, 1, time.Now()))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	answers, err := getAnswers(1, 0, 30, 0, db)
	if err != nil {
		t.Errorf("Error while getting answers: %s", err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	if len(answers) != 1 {
		t.Errorf("Length of answers should be 1")
	}
	if answers[0].Author.Id != 0 {
		t.Errorf("Author id should be 0")
	}
	if answers[0].LikeCount != 1 {
		t.Errorf("Like count should be 1")
	}
	if answers[0].IsAuthorAnonymous != true {
		t.Errorf("IsAuthorAnonymous should be true")
	}
}

func TestGetAnswersNotLikedByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, question_id, text, created_at FROM answer").WithArgs(1, 0, 30).WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "text", "created_at"}).AddRow(1, 1, "text", time.Now()))
	mock.ExpectQuery("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "created_at"}).AddRow(1, "text", 0, true, 1, time.Now()))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	answers, err := getAnswers(1, 1, 30, 0, db)
	if err != nil {
		t.Errorf("Error while getting answers: %s", err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	if len(answers) != 1 {
		t.Errorf("Length of answers should be 1")
	}
	if answers[0].Author.Id != 0 {
		t.Errorf("Author id should be 0")
	}
	if answers[0].LikeCount != 1 {
		t.Errorf("Like count should be 1")
	}
	if answers[0].IsAuthorAnonymous != true {
		t.Errorf("IsAuthorAnonymous should be true")
	}
	if answers[0].LikedByRequester {
		t.Errorf("LikedByRequester should be false")
	}
}

func TestGetAnswersLikedByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, question_id, text, created_at FROM answer").WithArgs(1, 0, 30).WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "text", "created_at"}).AddRow(1, 1, "text", time.Now()))
	mock.ExpectQuery("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "created_at"}).AddRow(1, "text", 0, true, 1, time.Now()))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	answers, err := getAnswers(1, 1, 30, 0, db)
	if err != nil {
		t.Errorf("Error while getting answers: %s", err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	if len(answers) != 1 {
		t.Errorf("Length of answers should be 1")
	}
	if answers[0].Author.Id != 0 {
		t.Errorf("Author id should be 0")
	}
	if answers[0].LikeCount != 1 {
		t.Errorf("Like count should be 1")
	}
	if answers[0].IsAuthorAnonymous != true {
		t.Errorf("IsAuthorAnonymous should be true")
	}
	if !answers[0].LikedByRequester {
		t.Errorf("LikedByRequester should be true")
	}
}

func TestGetAnswersErrorCheckLike(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, question_id, text, created_at FROM answer").WithArgs(1, 0, 30).WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "text", "created_at"}).AddRow(1, 1, "text", time.Now()))
	mock.ExpectQuery("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "created_at"}).AddRow(1, "text", 0, true, 1, time.Now()))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1, 1).WillReturnError(errors.New("error"))
	answers, err := getAnswers(1, 1, 30, 0, db)
	if err != nil {
		t.Errorf("Error while getting answers: %s", err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	if len(answers) != 1 {
		t.Errorf("Length of answers should be 1")
	}
	if answers[0].Author.Id != 0 {
		t.Errorf("Author id should be 0")
	}
	if answers[0].LikeCount != 1 {
		t.Errorf("Like count should be 1")
	}
	if answers[0].IsAuthorAnonymous != true {
		t.Errorf("IsAuthorAnonymous should be true")
	}
	if answers[0].LikedByRequester {
		t.Errorf("LikedByRequester should be false")
	}
}

func TestGetAnswerIdByQuestionId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Test success case
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT id FROM answer WHERE question_id = \\? AND has_been_deleted = 0").WithArgs(1).WillReturnRows(rows)

	answerId, err := GetAnswerIdByQuestionId(1, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if answerId != 1 {
		t.Errorf("unexpected answer id, want %d, got %d", 1, answerId)
	}

	// Test error case
	mock.ExpectQuery("SELECT id FROM answer WHERE question_id = \\? AND has_been_deleted = 0").WithArgs(2).WillReturnError(sql.ErrNoRows)

	answerId, err = GetAnswerIdByQuestionId(2, db)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if answerId != 0 {
		t.Errorf("unexpected answer id, want %d, got %d", 0, answerId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
