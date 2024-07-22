package client

import (
	"database/sql"
	"errors"
	"project_truthful/client/database"
	"project_truthful/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserProfileError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	database.DB = db

	// tests with error when getting user id
	mock.ExpectQuery("SELECT id FROM user").WithArgs("toto").WillReturnError(errors.New("error for test"))
	_, code, err := GetUserProfile("toto", 0, 30, 0)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if code != 500 {
		t.Errorf("Expected code 500, got %d", code)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// tests with profile not found
	mock.ExpectQuery("SELECT id FROM user").WithArgs("toto").WillReturnError(sql.ErrNoRows)
	_, code, err = GetUserProfile("toto", 0, 30, 0)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if code != 404 {
		t.Errorf("Expected code 404, got %d", code)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// tests with error when getting user profile infos
	mock.ExpectQuery("SELECT id FROM user").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(1).WillReturnError(errors.New("error for test"))
	_, code, err = GetUserProfile("toto", 0, 30, 0)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if code != 500 {
		t.Errorf("Expected code 500, got %d", code)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
}

func TestGetUserProfileSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	database.DB = db

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
	mock.ExpectQuery("SELECT id FROM user").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
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

	profile, code, err := GetUserProfile("toto", 0, 30, 0)

	if err != nil {
		t.Errorf("Error while getting user profile: %s", err.Error())
	}
	if code != 200 {
		t.Errorf("Expected code 200, got %d", code)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	// compares each field of the struct
	if profile.Id != expected.Id {
		t.Errorf("Database error: expected %d, got %d", expected.Id, profile.Id)
	}
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

func TestGetUserProfileSuccessUserRequestingSelf(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	database.DB = db

	creationTime := time.Now()
	mock.ExpectQuery("SELECT id FROM user").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
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

	profile, code, err := GetUserProfile("toto", 1, 30, 0)

	if err != nil {
		t.Errorf("Error while getting user profile: %s", err.Error())
	}
	if code != 200 {
		t.Errorf("Expected code 200, got %d", code)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	if profile.IsRequestingSelf != true {
		t.Errorf("Expected true, got false")
	}
}

func TestGetUserProfileSuccessUserNotFollowing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	database.DB = db

	creationTime := time.Now()
	mock.ExpectQuery("SELECT id FROM user").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
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

	// checks if user is followed by the requestern, should return false
	mock.ExpectQuery("SELECT COUNT").WithArgs(2, 1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))

	profile, code, err := GetUserProfile("toto", 2, 30, 0)

	if err != nil {
		t.Errorf("Error while getting user profile: %s", err.Error())
	}
	if code != 200 {
		t.Errorf("Expected code 200, got %d", code)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	if profile.IsRequestingSelf != false {
		t.Errorf("Expected false, got true")
	}
	if profile.IsFollowedByRequester != false {
		t.Errorf("Expected false, got true")
	}
}

func TestGetUserProfileSuccessUserFollowing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	database.DB = db

	creationTime := time.Now()
	mock.ExpectQuery("SELECT id FROM user").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
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

	// checks if user is followed by the requestern, should return false
	mock.ExpectQuery("SELECT COUNT").WithArgs(2, 1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))

	profile, code, err := GetUserProfile("toto", 2, 30, 0)

	if err != nil {
		t.Errorf("Error while getting user profile: %s", err.Error())
	}
	if code != 200 {
		t.Errorf("Expected code 200, got %d", code)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}

	if profile.IsRequestingSelf != false {
		t.Errorf("Expected false, got true")
	}
	if profile.IsFollowedByRequester != true {
		t.Errorf("Expected true, got false")
	}
}
