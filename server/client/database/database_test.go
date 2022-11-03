package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO user").WithArgs("username", "username", "password", "email", "birthdate").WillReturnResult(sqlmock.NewResult(4, 1))
	id, err := InsertUser("username", "password", "email", "birthdate", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if id != 4 {
		t.Errorf("id should be 4, but is %d", id)
	}
	if err != nil {
		t.Errorf("Error while inserting user: %s", err.Error())
	}

	mock.ExpectExec("INSERT INTO user").WithArgs("username_error", "username_error", "password", "email", "birthdate").WillReturnError(errors.New("error"))
	_, err = InsertUser("username_error", "password", "email", "birthdate", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestCheckUsernameExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckUsernameExists("username", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if username exists: %s", err.Error())
	}
	if !exists {
		t.Errorf("Username should exist")
	}
	// tests that the username does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckUsernameExists("toto", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if username exists: %s", err.Error())
	}
	if exists {
		t.Errorf("Username should not exist")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs("username_error").WillReturnError(errors.New("error"))
	_, err = CheckUsernameExists("username_error", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestCheckEmailExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs("email@email.fr").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckEmailExists("email@email.fr", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if email exists: %s", err.Error())
	}
	if !exists {
		t.Errorf("Email should exist")
	}
	// tests that the email does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs("toto@toto.fr").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckEmailExists("toto@toto.fr", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if email exists: %s", err.Error())
	}
	if exists {
		t.Errorf("Email should not exist")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs("unexistant@toto.fr").WillReturnError(errors.New("error"))
	_, err = CheckEmailExists("unexistant@toto.fr", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGetUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT id").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	id, err := GetUserId("username", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while getting user id: %s", err.Error())
	}
	if id != 1 {
		t.Errorf("id should be 1, but is %d", id)
	}

	mock.ExpectQuery("SELECT id").WithArgs("username_not_existing").WillReturnRows(sqlmock.NewRows([]string{"id"}))
	id, err = GetUserId("username_not_existing", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while getting user id: %s", err.Error())
	}
	if id != 0 {
		t.Errorf("id should be 0, but is %d", id)
	}

	mock.ExpectQuery("SELECT id").WithArgs("username_error").WillReturnError(errors.New("error"))
	_, err = GetUserId("username_error", db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestGetHashedPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT password").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("password"))
	password, err := GetHashedPassword(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while getting hashed password: %s", err.Error())
	}
	if password != "password" {
		t.Errorf("password should be 'password', but is %s", password)
	}

	mock.ExpectQuery("SELECT password").WithArgs(3).WillReturnError(errors.New("error"))
	_, err = GetHashedPassword(3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestCheckUserIdExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckUserIdExists(1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if user id exists: %s", err.Error())
	}
	if !exists {
		t.Errorf("User id should exist")
	}
	// tests that the user id does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckUserIdExists(2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if user id exists: %s", err.Error())
	}
	if exists {
		t.Errorf("User id should not exist")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(3).WillReturnError(errors.New("error"))
	_, err = CheckUserIdExists(3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestCheckFollowExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckFollowExists(1, 2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if follow exists: %s", err.Error())
	}
	if !exists {
		t.Errorf("Follow should exist")
	}
	// tests that the follow does not exist
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckFollowExists(1, 3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while checking if follow exists: %s", err.Error())
	}
	if exists {
		t.Errorf("Follow should not exist")
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 4).WillReturnError(errors.New("error"))
	_, err = CheckFollowExists(1, 4, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestAddFollow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectExec("INSERT INTO follow").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = AddFollow(1, 2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while adding follow: %s", err.Error())
	}

	mock.ExpectExec("INSERT INTO follow").WithArgs(1, 3).WillReturnError(errors.New("error"))
	err = AddFollow(1, 3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

func TestRemoveFollow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()
	mock.ExpectExec("DELETE FROM follow").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = RemoveFollow(1, 2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Error while removing follow: %s", err.Error())
	}

	mock.ExpectExec("DELETE FROM follow").WithArgs(1, 3).WillReturnError(errors.New("error"))
	err = RemoveFollow(1, 3, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}

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

func TestCheckPostIdExists(t *testing.T) {
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

func TestCheckLikeExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error
	mock.ExpectQuery("SELECT COUNT").WithArgs(1, 1).WillReturnError(errors.New("error for db test"))
	_, err = CheckLikeExists(1, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	// Test with one row returned
	mock.ExpectQuery("SELECT COUNT").WithArgs(2, 2).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	exists, err := CheckLikeExists(2, 2, db)
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
	mock.ExpectQuery("SELECT COUNT").WithArgs(3, 3).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	exists, err = CheckLikeExists(3, 3, db)
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

func TestAddLike(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error
	mock.ExpectExec("INSERT INTO answer_like").WithArgs(1, 1).WillReturnError(errors.New("error for db test"))
	err = AddLike(1, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	//Test with no error
	mock.ExpectExec("INSERT INTO answer_like").WithArgs(2, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = AddLike(2, 2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
}

func TestRemoveLike(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating sqlmock: %s", err.Error())
	}
	defer db.Close()

	// Test with an SQL error
	mock.ExpectExec("DELETE FROM answer_like").WithArgs(1, 1).WillReturnError(errors.New("error for db test"))
	err = RemoveLike(1, 1, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Database error: expected error, got nil")
	}

	//Test with no error
	mock.ExpectExec("DELETE FROM answer_like").WithArgs(2, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = RemoveLike(2, 2, db)
	if mock.ExpectationsWereMet() != nil {
		t.Errorf("Error while checking expectations: %s", err.Error())
	}
	if err != nil {
		t.Errorf("Database error: expected nil, got %s", err.Error())
	}
}
