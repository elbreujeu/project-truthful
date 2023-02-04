package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"project_truthful/client/basicfuncs"
	"project_truthful/client/database"
	"project_truthful/helpunittesting"
	"project_truthful/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// this function tests the SetupRoutes function
func TestSetupRoutes(t *testing.T) {
	router := gin.Default()
	SetupRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello_world", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, `{"message":"Hello world !"}`, w.Body.String())
}

// func TestMiddleware(t *testing.T) {
// 	router := mux.NewRouter()
// 	SetupRoutes(router)
// 	SetMiddleware(router)
// 	r, _ := http.NewRequest("GET", "/hello_world", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, r)
// 	if w.Code != http.StatusOK {
// 		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
// 	}
// 	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
// 	assert.Equal(t, "GET", w.Header().Get("Access-Control-Allow-Methods"))
// 	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Headers"))
// 	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
// }

func TestRegister(t *testing.T) {
	router := gin.Default()
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
	assert.Equal(t, `{"error":"invalid character '\u003c' looking for beginning of value","message":"invalid request body"}`, w.Body.String())

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
	assert.Equal(t, `{"error":"missing fields","message":"invalid request body"}`, w.Body.String())

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
	assert.Equal(t, `{"error":"username must be between 3 and 20 characters","message":"error while creating user"}`, w.Body.String())

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
	assert.Equal(t, `{"id":1,"message":"User created"}`, w.Body.String())
	os.Setenv("IS_TEST", "false")
}

func TestLogin(t *testing.T) {
	router := gin.Default()
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
	assert.Equal(t, `{"error":"invalid character '\u003c' looking for beginning of value","message":"invalid request body"}`, w.Body.String())

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
	assert.Equal(t, `{"error":"missing fields","message":"invalid request body"}`, w.Body.String())

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
	assert.Equal(t, `{"error":"user not found","message":"error while logging in"}`, w.Body.String())

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
	assert.Equal(t, `{"message":"User logged in","token":"test"}`, w.Body.String())
	os.Setenv("IS_TEST", "false")
}

func TestRefreshToken(t *testing.T) {
	router := gin.Default()
	SetupRoutes(router)
	SetMiddleware(router)

	// With invalid format token
	r, _ := http.NewRequest("GET", "/refresh_token", nil)
	r.Header.Set("Authorization", "invalid_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"missing fields","message":"error while parsing token"}`, w.Body.String())

	// With valid format token, but invalid token
	r, _ = http.NewRequest("GET", "/refresh_token", nil)
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"error":"token contains an invalid number of segments","message":"error while checking token"}`, w.Body.String())

	// With is_test env variable set to true (token is impossible to check due to the fact that the key is not the same)
	os.Setenv("IS_TEST", "true")
	r, _ = http.NewRequest("GET", "/refresh_token", nil)
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, `{"message":"Token refreshed","token":"test"}`, w.Body.String())
	os.Setenv("IS_TEST", "false")
}

func TestGetUserProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db
	router := gin.Default()
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
	assert.Equal(t, `{"error":"error","message":"error while getting user"}`, w.Body.String())

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
	assert.Equal(t, `{"id":1,"username":"username","display_name":"display_name","follower_count":1,"following_count":1,"answer_count":1,"answers":null}`, w.Body.String())
}

func TestFollowUser(t *testing.T) {
	// With invalid format token
	router := gin.Default()
	SetupRoutes(router)
	SetMiddleware(router)
	r, _ := http.NewRequest("POST", "/follow_user", nil)
	r.Header.Set("Authorization", "invalid_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"missing fields","message":"error while parsing token"}`, w.Body.String())

	os.Setenv("IS_TEST", "true")

	// Test with nil request body
	r, _ = http.NewRequest("POST", "/follow_user", nil)
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"invalid request","message":"invalid request body"}`, w.Body.String())

	// Test with invalid request body
	requestBody := []byte(`{"username": "toto"}`)
	r, _ = http.NewRequest("POST", "/follow_user", bytes.NewBuffer(requestBody))
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"missing fields","message":"invalid request body"}`, w.Body.String())

	// Test with invalid request body
	requestBody = []byte("<invalid json>")
	r, _ = http.NewRequest("POST", "/follow_user", bytes.NewBuffer(requestBody))
	r.Header.Set("Authorization", "Bearer 123456789")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"invalid character '\u003c' looking for beginning of value","message":"invalid request body"}`, w.Body.String())

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
	assert.Equal(t, `{"error":"error","message":"error while following user"}`, w.Body.String())

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
	assert.Equal(t, `{"message":"User followed"}`, w.Body.String())

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
	assert.Equal(t, `{"message":"User unfollowed"}`, w.Body.String())

	os.Setenv("IS_TEST", "false")
}

func TestAskQuestion(t *testing.T) {
	// With valid format token but invalid token
	router := gin.Default()
	SetupRoutes(router)
	SetMiddleware(router)
	r, _ := http.NewRequest("POST", "/ask_question", nil)
	r.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"error":"token contains an invalid number of segments","message":"error while checking token"}`, w.Body.String())

	// With nil body
	r, _ = http.NewRequest("POST", "/ask_question", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"invalid request","message":"invalid request body"}`, w.Body.String())

	// With invalid body
	body := []byte(`<invalid body>`)
	r, _ = http.NewRequest("POST", "/ask_question", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"invalid character '\u003c' looking for beginning of value","message":"invalid request body"}`, w.Body.String())

	// Test for error when asking question
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB = db

	body = []byte(`{"user_id": 1, "text":"question"}`)
	r, _ = http.NewRequest("POST", "/ask_question", bytes.NewBuffer(body))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectExec("INSERT INTO question").WithArgs("question", "", 1).WillReturnError(errors.New("error"))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"error":"error","message":"error while asking question"}`, w.Body.String())

	// Success
	body = []byte(`{"user_id": 1, "text":"question"}`)
	r, _ = http.NewRequest("POST", "/ask_question", bytes.NewBuffer(body))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectExec("INSERT INTO question").WithArgs("question", "", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
	assert.Equal(t, `{"id":1,"message":"Question asked"}`, w.Body.String())
}

func TestGetQuestionsFailQueryParameters(t *testing.T) {
	router := gin.Default()
	SetupRoutes(router)
	SetMiddleware(router)
	r, _ := http.NewRequest("GET", "/get_questions?count=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"strconv.Atoi: parsing \"abc\": invalid syntax","message":"invalid count"}`, w.Body.String())

	r, _ = http.NewRequest("GET", "/get_questions?count=1&start=abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"strconv.Atoi: parsing \"abc\": invalid syntax","message":"invalid count"}`, w.Body.String())
}

func TestGetQuestionsWithParameters(t *testing.T) {
	// With valid format token but invalid token
	router := gin.Default()
	SetupRoutes(router)
	SetMiddleware(router)
	r, _ := http.NewRequest("GET", "/get_questions", nil)
	r.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"error":"token contains an invalid number of segments","message":"error while checking token"}`, w.Body.String())

	// Test for error while getting questions
	os.Setenv("IS_TEST", "true")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB = db
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnError(errors.New("error for checking user"))
	r, _ = http.NewRequest("GET", "/get_questions", nil)
	r.Header.Set("Authorization", "Bearer token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"error":"error for checking user","message":"error while getting questions"}`, w.Body.String())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test for success while getting questions
	// Generate questions
	questionTime := time.Now()
	questions := helpunittesting.GenerateTestQuestions(10, 1, questionTime)
	rows := sqlmock.NewRows([]string{"id", "text", "author_id", "is_author_anonymous", "receiver_id", "creation_date"})
	for _, question := range questions {
		rows.AddRow(question.Id, question.Text, question.AuthorId, question.IsAuthorAnonymous, question.ReceiverId, question.CreatedAt)
	}
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT").WithArgs(1, 0, 10).WillReturnRows(rows)

	for _, question := range questions {
		mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(question.Id).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
	}
	r, _ = http.NewRequest("GET", "/get_questions", nil)
	r.Header.Set("Authorization", "Bearer token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	var returnedQuestions []models.Question
	if err := json.Unmarshal(w.Body.Bytes(), &returnedQuestions); err != nil {
		t.Errorf("Error while unmarshalling response: %s", err)
	}
	if len(returnedQuestions) != 10 {
		t.Errorf("Expected 10 questions, got %d", len(returnedQuestions))
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
	}

	os.Setenv("IS_TEST", "false")
}

func TestAnswerQuestion(t *testing.T) {
	// With valid format token but invalid token
	router := gin.Default()
	SetupRoutes(router)
	SetMiddleware(router)
	r, _ := http.NewRequest("POST", "/answer_question", nil)
	r.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"error":"token contains an invalid number of segments","message":"error while checking token"}`, w.Body.String())

	os.Setenv("IS_TEST", "true")

	// Test with invalid request body
	requestBody := bytes.NewBuffer([]byte(`<invalid json>`))
	r, _ = http.NewRequest("POST", "/answer_question", requestBody)
	r.Header.Set("Authorization", "Bearer valid_token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"invalid character '\u003c' looking for beginning of value","message":"error while parsing request body"}`, w.Body.String())

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db

	// Test with error when answering question
	str := basicfuncs.GenerateRandomString(1500)
	requestBody = bytes.NewBuffer([]byte(`{"question_id": 1, "text": "` + str + `"}`))
	r, _ = http.NewRequest("POST", "/answer_question", requestBody)
	r.Header.Set("Authorization", "Bearer valid_token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"answer is too long","message":"error while answering question"}`, w.Body.String())

	// Test success
	requestBody = bytes.NewBuffer([]byte(`{"question_id": 1, "text": "answer"}`))
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT receiver_id FROM question").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO answer").WithArgs(1, 1, "answer", "").WillReturnResult(sqlmock.NewResult(1, 1))
	r, _ = http.NewRequest("POST", "/answer_question", requestBody)
	r.Header.Set("Authorization", "Bearer valid_token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
	assert.Equal(t, `{"id":1,"message":"question answered"}`, w.Body.String())

	os.Setenv("IS_TEST", "false")
}

func TestLikeAnswerErrors(t *testing.T) {
	// With valid format token but invalid token
	router := gin.Default()
	SetupRoutes(router)
	SetMiddleware(router)
	r, _ := http.NewRequest("POST", "/like_answer", nil)
	r.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"error":"token contains an invalid number of segments","message":"error while checking token"}`, w.Body.String())

	os.Setenv("IS_TEST", "true")

	// Test with invalid request body
	requestBody := bytes.NewBuffer([]byte(`<invalid json>`))
	r, _ = http.NewRequest("POST", "/like_answer", requestBody)
	r.Header.Set("Authorization", "Bearer valid_token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	assert.Equal(t, `{"error":"invalid character '\u003c' looking for beginning of value","message":"error while parsing request body"}`, w.Body.String())
	os.Setenv("IS_TEST", "false")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db

	// Test with error when liking answer
	os.Setenv("IS_TEST", "true")
	requestBody = bytes.NewBuffer([]byte(`{"answer_id": 1, "like": true}`))
	r, _ = http.NewRequest("POST", "/like_answer", requestBody)
	r.Header.Set("Authorization", "Bearer valid_token")
	w = httptest.NewRecorder()
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnError(errors.New("error when getting user"))
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"error":"error when getting user","message":"error while liking answer"}`, w.Body.String())

	// Test with error when unliking answer
	requestBody = bytes.NewBuffer([]byte(`{"answer_id": 1, "like": false}`))
	r, _ = http.NewRequest("POST", "/like_answer", requestBody)
	r.Header.Set("Authorization", "Bearer valid_token")
	w = httptest.NewRecorder()
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1, 1).WillReturnError(errors.New("error when getting user"))
	router.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	assert.Equal(t, `{"error":"error when getting user","message":"error while unliking answer"}`, w.Body.String())

	os.Setenv("IS_TEST", "false")
}

func TestLikeAnswerSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db
	router := gin.Default()
	SetupRoutes(router)
	SetMiddleware(router)
	os.Setenv("IS_TEST", "true")

	// Test success
	requestBody := bytes.NewBuffer([]byte(`{"answer_id": 1, "like": true}`))
	r, _ := http.NewRequest("POST", "/like_answer", requestBody)
	r.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()
	mock.ExpectQuery("SELECT COUNT(.+) FROM user").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO answer_like").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	router.ServeHTTP(w, r)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
	assert.Equal(t, `{"message":"answer liked"}`, w.Body.String())

	os.Setenv("IS_TEST", "false")
}

func TestUnlikeAnswerSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.DB = db
	router := gin.Default()
	SetupRoutes(router)
	SetMiddleware(router)
	os.Setenv("IS_TEST", "true")

	// Test success
	requestBody := bytes.NewBuffer([]byte(`{"answer_id": 1, "like": false}`))
	r, _ := http.NewRequest("POST", "/like_answer", requestBody)
	r.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()
	mock.ExpectQuery("SELECT COUNT(.+) FROM answer_like").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectExec("DELETE FROM answer_like").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	router.ServeHTTP(w, r)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	assert.Equal(t, `{"message":"answer unliked"}`, w.Body.String())

	os.Setenv("IS_TEST", "false")
}
