package database

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetRateLimitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ip := "192.168.0.1"
	mock.ExpectQuery("SELECT (.+) FROM rate_limit WHERE ip_address = ?").WithArgs(ip).WillReturnError(fmt.Errorf("some error"))

	_, err = GetRateLimit(ip, db)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetRateLimitNoRowsErrorInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ip := "192.168.0.1"
	mock.ExpectQuery("SELECT (.+) FROM rate_limit WHERE ip_address = ?").WithArgs(ip).WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO rate_limit (.+) VALUES (.+)").WithArgs(ip).WillReturnError(fmt.Errorf("some error"))

	_, err = GetRateLimit(ip, db)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetRateLimitNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ip := "192.168.0.1"
	mock.ExpectQuery("SELECT (.+) FROM rate_limit WHERE ip_address = ?").WithArgs(ip).WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO rate_limit (.+) VALUES (.+)").WithArgs(ip).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = GetRateLimit(ip, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestGetRateLimitSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ip := "192.168.0.1"
	rows := sqlmock.NewRows([]string{"ip_address", "request_count", "last_updated"}).AddRow(ip, 1, time.Now())
	mock.ExpectQuery("SELECT (.+) FROM rate_limit WHERE ip_address = ?").WithArgs(ip).WillReturnRows(rows)

	_, err = GetRateLimit(ip, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestResetRateLimitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ip := "192.168.0.1"

	mock.ExpectExec("UPDATE rate_limit SET (.+) WHERE ip_address = ?").WithArgs(ip).WillReturnError(fmt.Errorf("some error"))

	err = ResetRateLimit(ip, db)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestResetRateLimitSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ip := "192.168.0.1"

	mock.ExpectExec("UPDATE rate_limit SET (.+) WHERE ip_address = ?").WithArgs(ip).WillReturnResult(sqlmock.NewResult(1, 1))

	err = ResetRateLimit(ip, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestIncrementRateLimitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ip := "192.168.0.1"

	mock.ExpectExec("UPDATE rate_limit SET (.+) WHERE ip_address = ?").WithArgs(ip).WillReturnError(fmt.Errorf("some error"))

	err = IncrementRateLimit(ip, db)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestIncrementRateLimitSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ip := "192.168.0.1"

	mock.ExpectExec("UPDATE rate_limit SET (.+) WHERE ip_address = ?").WithArgs(ip).WillReturnResult(sqlmock.NewResult(1, 1))

	err = IncrementRateLimit(ip, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
