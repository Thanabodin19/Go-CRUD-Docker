package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUsers_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "F_name", "L_name"}).
		AddRow(1, "John", "Doe")

	mock.ExpectQuery("SELECT \\* FROM humans").
		WillReturnRows(rows)

	handler := getUsers(db)

	req := httptest.NewRequest(http.MethodGet, "/humans", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}
