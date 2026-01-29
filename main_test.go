package main

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

func setupRouter(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/humans", getUsers(db)).Methods("GET")
	r.HandleFunc("/humans/{id}", getUser(db)).Methods("GET")
	r.HandleFunc("/humans", createUser(db)).Methods("POST")
	r.HandleFunc("/humans/{id}", updateUser(db)).Methods("PUT")
	r.HandleFunc("/humans/{id}", deleteUser(db)).Methods("DELETE")
	return r
}

func TestNewRouter_Health(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	handler := newRouter(db)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}


func TestRootHandler_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler := root()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if body == "" {
		t.Fatal("response body should not be empty")
	}
}


func TestJSONContentTypeMiddleware(t *testing.T) {
	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware := jsonContentTypeMiddleware(next)
	middleware.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next handler was not called")
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}
}


func TestGetUsers_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "f_name", "l_name"}).
		AddRow(1, "John", "Doe")

	mock.ExpectQuery("SELECT \\* FROM humans").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/humans", nil)
	rr := httptest.NewRecorder()

	router := setupRouter(db)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}



func TestGetUser_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "f_name", "l_name"}).
		AddRow(1, "John", "Doe")

	mock.ExpectQuery("SELECT \\* FROM humans WHERE id =").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/humans/1", nil)
	rr := httptest.NewRecorder()

	router := setupRouter(db)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200")
	}
}

func TestGetUser_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT \\* FROM humans WHERE id =").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest("GET", "/humans/99", nil)
	rr := httptest.NewRecorder()

	router := setupRouter(db)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404")
	}
}

func TestCreateUser_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("INSERT INTO humans").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	body := `{"F_name":"Jane","L_name":"Doe"}`
	req := httptest.NewRequest("POST", "/humans", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter(db)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200")
	}
}

func TestUpdateUser_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("UPDATE humans").
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"F_name":"Updated","L_name":"Name"}`
	req := httptest.NewRequest("PUT", "/humans/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter(db)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200")
	}
}

func TestDeleteUser_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "f_name", "l_name"}).
		AddRow(1, "John", "Doe")

	mock.ExpectQuery("SELECT \\* FROM humans WHERE id =").
		WillReturnRows(rows)

	mock.ExpectExec("DELETE FROM humans").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest("DELETE", "/humans/1", nil)
	rr := httptest.NewRecorder()

	router := setupRouter(db)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT \\* FROM humans WHERE id =").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest("DELETE", "/humans/99", nil)
	rr := httptest.NewRecorder()

	router := setupRouter(db)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404")
	}
}






