package database

import (
	"errors"
	"project_truthful/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserProfileInfosError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error on query to get name and display name
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(1).WillReturnError(errors.New("error for db test"))
	_, err = GetUserProfileInfos(1, 30, 0, db)
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// Test with an error when getting follower count
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username", "display_name"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(2).WillReturnError(errors.New("error for db test"))
	_, err = GetUserProfileInfos(2, 30, 0, db)
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// Test with an error when getting following count
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username", "display_name"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(3).WillReturnError(errors.New("error for db test"))
	_, err = GetUserProfileInfos(3, 30, 0, db)
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// Test with an error when getting answer count
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username", "display_name"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(4).WillReturnError(errors.New("error for db test"))
	_, err = GetUserProfileInfos(4, 30, 0, db)
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// Test with an error when getting answers
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username", "display_name"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT id, question_id, text, created_at FROM answer").WithArgs(5, 0, 30).WillReturnError(errors.New("error for db test"))
	_, err = GetUserProfileInfos(5, 30, 0, db)
	if err == nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
}

func TestGetProfileInfos(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	creationTime := time.Now()

	var expected models.UserProfileInfos
	expected.Id = 1
	expected.Username = "username"
	expected.DisplayName = "display_name"
	expected.FollowerCount = 1
	expected.FollowingCount = 1
	expected.AnswerCount = 1
	expected.Answers = []models.Answer{
		{
			Id:                1,
			IsAuthorAnonymous: false,
			Author: models.UserPreview{
				Id:          2,
				Username:    "username_author",
				DisplayName: "display_name_author",
			},
			QuestionText: "question_text",
			AnswerText:   "answer_text",
			CreatedAt:    creationTime,
			LikeCount:    1,
		},
	}
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username", "display_name"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT id, question_id, text, created_at FROM answer").WithArgs(1, 0, 30).WillReturnRows(sqlmock.NewRows([]string{"id", "question_id", "text", "created_at"}).AddRow(1, 1, "answer_text", creationTime))

	questionRows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "created_at"}).
		AddRow(1, "question_text", 2, false, 1, creationTime)
	mock.ExpectQuery("SELECT id, text, author_id, is_author_anonymous, receiver_id, created_at FROM question").WithArgs(1).WillReturnRows(questionRows)
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username_author", "display_name_author"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	profile, err := GetUserProfileInfos(1, 30, 0, db)
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
	if profile.Answers[0].Id != expected.Answers[0].Id {
		t.Errorf("Database error: expected %d, got %d", expected.Answers[0].Id, profile.Answers[0].Id)
	}
	if profile.Answers[0].QuestionText != expected.Answers[0].QuestionText {
		t.Errorf("Database error: expected %s, got %s", expected.Answers[0].QuestionText, profile.Answers[0].QuestionText)
	}
	if profile.Answers[0].AnswerText != expected.Answers[0].AnswerText {
		t.Errorf("Database error: expected %s, got %s", expected.Answers[0].AnswerText, profile.Answers[0].AnswerText)
	}
	if profile.Answers[0].CreatedAt != expected.Answers[0].CreatedAt {
		t.Errorf("Database error: expected %s, got %s", expected.Answers[0].CreatedAt, profile.Answers[0].CreatedAt)
	}
	if profile.Answers[0].LikeCount != expected.Answers[0].LikeCount {
		t.Errorf("Database error: expected %d, got %d", expected.Answers[0].LikeCount, profile.Answers[0].LikeCount)
	}
	if profile.Answers[0].Author.Id != expected.Answers[0].Author.Id {
		t.Errorf("Database error: expected %d, got %d", expected.Answers[0].Author.Id, profile.Answers[0].Author.Id)
	}
	if profile.Answers[0].Author.Username != expected.Answers[0].Author.Username {
		t.Errorf("Database error: expected %s, got %s", expected.Answers[0].Author.Username, profile.Answers[0].Author.Username)
	}
	if profile.Answers[0].Author.DisplayName != expected.Answers[0].Author.DisplayName {
		t.Errorf("Database error: expected %s, got %s", expected.Answers[0].Author.DisplayName, profile.Answers[0].Author.DisplayName)
	}
	if profile.Answers[0].IsAuthorAnonymous != expected.Answers[0].IsAuthorAnonymous {
		t.Errorf("Database error: expected %t, got %t", expected.Answers[0].IsAuthorAnonymous, profile.Answers[0].IsAuthorAnonymous)
	}
}
