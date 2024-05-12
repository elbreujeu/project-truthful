package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBanUserError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db

	// test with error while checking moderator status
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New("error while checking moderator status"))
	_, code, err := BanUser(1, 1, 1, "reason")
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with error while checking admin status
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(0))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New("error while checking admin status"))
	_, code, err = BanUser(1, 1, 1, "reason")
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with user not being a moderator or admin
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(0))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(0))
	_, code, err = BanUser(1, 1, 1, "reason")
	if code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with error while checking user id exists
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New("error while checking user id exists"))
	_, code, err = BanUser(1, 1, 1, "reason")
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with user not found
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(0))
	_, code, err = BanUser(1, 1, 1, "reason")
	if code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with user being self
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
	_, code, err = BanUser(1, 1, 1, "reason")
	if code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with error while checking admin status of user
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(2).WillReturnError(errors.New("error while checking admin status of user"))
	_, code, err = BanUser(2, 1, 1, "reason")
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// // test with user being an admin
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	_, code, err = BanUser(2, 1, 1, "reason")
	if code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with error while banning user
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(0))
	mock.ExpectExec("INSERT INTO ban").WithArgs(2, 1, "reason").WillReturnError(errors.New("error while banning user"))
	_, code, err = BanUser(2, 1, 1, "reason")
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestBanUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db

	// test with success
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(0))
	mock.ExpectExec("INSERT INTO ban").WithArgs(2, 1, "ban reason").WillReturnResult(sqlmock.NewResult(1, 1))
	_, code, err := BanUser(2, 1, 0, "ban reason")
	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}
	if err != nil {
		t.Errorf("expected nil, got error")
	}
}

func TestPardonUserError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db

	// test with error while checking moderator status
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New("error while checking moderator status"))
	_, code, err := PardonUser(1, 1)
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with error while checking admin status
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(0))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New("error while checking admin status"))
	_, code, err = PardonUser(1, 1)
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with user not being a moderator or admin
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(0))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(0))
	_, code, err = PardonUser(1, 1)
	if code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with error while checking ban exists by ban id
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New("error while checking ban exists by ban id"))
	_, code, err = PardonUser(1, 1)
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with ban not found
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	_, code, err = PardonUser(1, 1)
	if code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with error while checking pardon exists
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New("error while checking pardon exists"))
	_, code, err = PardonUser(1, 1)
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// test with error while pardoning user
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO pardon").WithArgs(1, 1).WillReturnError(errors.New("error while pardoning user"))

	_, code, err = PardonUser(1, 1)
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestPardonUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db

	// test with success
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_moderator"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO pardon").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	_, code, err := PardonUser(1, 1)
	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}
	if err != nil {
		t.Errorf("expected nil, got error")
	}
}
