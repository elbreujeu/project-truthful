package routes

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"project_truthful/client/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// this function tests the SetupRoutes function
func TestSetupRoutes(t *testing.T) {
	router := mux.NewRouter()
	SetupRoutes(router)
	r, _ := http.NewRequest("GET", "/hello_world", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, `{"message": "Hello world !"}`, w.Body.String())
}

func TestMiddleware(t *testing.T) {
	router := mux.NewRouter()
	SetupRoutes(router)
	SetMiddleware(router)
	r, _ := http.NewRequest("GET", "/hello_world", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestRegister(t *testing.T) {
	router := mux.NewRouter()
	SetupRoutes(router)
	SetMiddleware(router)
	// tests for json decoder error, todo
	body := "<invalid json>"
	r, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"message": "Invalid request body", "error": "invalid character '<' looking for beginning of value"}`, w.Body.String())

	// checks with empty username
	r, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"username": "", "password": "Toto123@", "email_address": "toto@toto.fr", "birthdate": "1990-01-01"}`)))
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"message": "Invalid request body", "error": "missing fields"}`, w.Body.String())

	// checks with register failure (invalid username)
	r, err = http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"username": "to", "password": "Toto123@", "email_address": "toto@toto.fr", "birthdate": "1990-01-01"}`)))
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"message": "error while creating user", "error": "username must be between 3 and 20 characters"}`, w.Body.String())

	// checks with register success
	var mock sqlmock.Sqlmock
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer database.DB.Close()
	// sets the env variable for password hashing
	os.Setenv("IS_TEST", "true")
	// expect a query to check if the username is already taken
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	// expect a query to check if the email address is already taken
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs("toto@toto.fr").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	// expect a query to insert the user
	mock.ExpectExec("INSERT INTO user").WithArgs("toto", "toto", "Toto123@", "toto@toto.fr", "1990-01-01").WillReturnResult(sqlmock.NewResult(1, 1))
	r, err = http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"username": "toto", "password": "Toto123@", "email_address": "toto@toto.fr", "birthdate": "1990-01-01"}`)))
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
	assert.Equal(t, `{"message": "User created", "id": 1}`, w.Body.String())
	os.Setenv("IS_TEST", "false")
}

func TestLogin(t *testing.T) {
	router := mux.NewRouter()
	SetupRoutes(router)
	SetMiddleware(router)
	// tests for json decoder error, todo
	body := "<invalid json>"
	r, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"message": "Invalid request body", "error": "invalid character '<' looking for beginning of value"}`, w.Body.String())

	// checks with empty username
	r, err := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username": "", "password": "Toto123@"}`)))
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"message": "Invalid request body", "error": "missing fields"}`, w.Body.String())

	var mock sqlmock.Sqlmock
	database.DB, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer database.DB.Close()

	// checks with login failure (invalid username)
	mock.ExpectQuery("SELECT (.+) FROM user").WithArgs("toto").WillReturnError(sql.ErrNoRows)
	r, err = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username": "toto", "password": "Toto123@"}`)))
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
	assert.Equal(t, `{"message": "error while logging in", "error": "user not found"}`, w.Body.String())

	// checks with login success
	os.Setenv("IS_TEST", "true")
	mock.ExpectQuery("SELECT id FROM user").WithArgs("toto").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT password FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("Toto123@"))
	r, err = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username": "toto", "password": "Toto123@"}`)))
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, `{"message": "User logged in", "token": "test"}`, w.Body.String())
	os.Setenv("IS_TEST", "false")
}

func TestRefreshToken(t *testing.T) {
	// With invalid format token
	router := mux.NewRouter()
	SetupRoutes(router)
	SetMiddleware(router)
	r, _ := http.NewRequest("GET", "/refresh_token", nil)
	r.Header.Set("Authorization", "invalid_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"message": "error while parsing token", "error": "missing fields"}`, w.Body.String())

	// With valid format token, but invalid token
	r, _ = http.NewRequest("GET", "/refresh_token", nil)
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"message": "error while checking token", "error": "token contains an invalid number of segments"}`, w.Body.String())

	// With is_test env variable set to true (token is impossible to check due to the fact that the key is not the same)
	os.Setenv("IS_TEST", "true")
	r, _ = http.NewRequest("GET", "/refresh_token", nil)
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, `{"message": "Token refreshed", "token": "test"}`, w.Body.String())
	os.Setenv("IS_TEST", "false")
}

func TestGetUserProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db
	router := mux.NewRouter()
	SetupRoutes(router)
	SetMiddleware(router)
	// tests for error when getting user profile
	mock.ExpectQuery("SELECT (.+) FROM user").WithArgs("toto").WillReturnError(errors.New("error"))
	r, _ := http.NewRequest("GET", "/get_user_profile/toto", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"message": "error while getting user", "error": "error"}`, w.Body.String())

	// tests for success
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT id FROM user WHERE username = \\?").WithArgs("username").WillReturnRows(rows)
	r, _ = http.NewRequest("GET", "/get_user_profile/username", nil)
	mock.ExpectQuery("SELECT username, display_name FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"username", "display_name"}).AddRow("username", "display_name"))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, `{"id":1,"username":"username","display_name":"display_name","follower_count":1,"following_count":1,"answer_count":1,"answers":null}`+"\n", w.Body.String())
}

func TestFollowUser(t *testing.T) {
	// With invalid format token
	router := mux.NewRouter()
	SetupRoutes(router)
	SetMiddleware(router)
	r, _ := http.NewRequest("POST", "/follow_user", nil)
	r.Header.Set("Authorization", "invalid_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"message": "error while parsing token", "error": "missing fields"}`, w.Body.String())

	os.Setenv("IS_TEST", "true")

	// Test with nil request body
	r, _ = http.NewRequest("POST", "/follow_user", nil)
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"message": "Invalid request body", "error": "missing fields"}`, w.Body.String())

	// Test with invalid request body
	requestBody := []byte(`{"username": "toto"}`)
	r, _ = http.NewRequest("POST", "/follow_user", bytes.NewBuffer(requestBody))
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"message": "Invalid request body", "error": "missing fields"}`, w.Body.String())

	// Test with invalid request body
	requestBody = []byte("<invalid json>")
	r, _ = http.NewRequest("POST", "/follow_user", bytes.NewBuffer(requestBody))
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"message": "Invalid request body", "error": "invalid character '<' looking for beginning of value"}`, w.Body.String())

	// tests for error when following
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO follow").WillReturnError(errors.New("error"))
	requestBody = []byte(`{"user_id":2, "follow":true}`)
	r, _ = http.NewRequest("POST", "/follow_user", bytes.NewBuffer(requestBody))
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"message": "error while following user", "error": "error"}`, w.Body.String())

	// tests for following success
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO follow").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	requestBody = []byte(`{"user_id":2, "follow":true}`)
	r, _ = http.NewRequest("POST", "/follow_user", bytes.NewBuffer(requestBody))
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, `{"message": "User followed"}`, w.Body.String())

	// tests for unfollowing success
	mock.ExpectQuery("SELECT COUNT(.+) FROM follow").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectExec("DELETE FROM follow").WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	requestBody = []byte(`{"user_id":2, "follow":false}`)
	r, _ = http.NewRequest("POST", "/follow_user", bytes.NewBuffer(requestBody))
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, `{"message": "User unfollowed"}`, w.Body.String())

	os.Setenv("IS_TEST", "false")
}
