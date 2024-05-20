package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/koshikovme/assignment3-greenlight/internal/data"
	"github.com/koshikovme/assignment3-greenlight/internal/jsonlog"
	"github.com/koshikovme/assignment3-greenlight/internal/mailer"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Connect to the test database.
	db, err := sql.Open("postgres", "your_test_dsn_here")
	if err != nil {
		t.Fatal(err)
	}

	// Clean the database before running tests.
	_, err = db.Exec("TRUNCATE TABLE movies RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

// newTestApplication creates a new application struct with mocked dependencies for testing.
func newTestApplication(t *testing.T) *application {
	t.Helper()

	db := setupTestDB(t)

	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	return &application{
		config: config{},
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.Mailer{},
	}
}

func TestRegisterUserHandler(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
	}{
		{
			name: "Valid User Registration",
			payload: map[string]string{
				"name":     "John Doe",
				"email":    "john.doe@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name: "Missing Fields",
			payload: map[string]string{
				"name":  "Jane Doe",
				"email": "jane.doe@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid Email",
			payload: map[string]string{
				"name":     "Invalid Email",
				"email":    "invalid-email",
				"password": "password123",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(app.registerUserHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestActivateUserHandler(t *testing.T) {
	app := newTestApplication(t)

	// Mock user and token for the test
	user := &data.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Activated: false,
	}
	token := "valid-token"

	// Add mocked user and token to the app's models for the test
	app.models.Users.Insert(user)
	app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
	}{
		{
			name: "Valid Token",
			payload: map[string]string{
				"token": token,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Token",
			payload: map[string]string{
				"token": "invalid-token",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/v1/users/activate", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(app.activateUserHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
