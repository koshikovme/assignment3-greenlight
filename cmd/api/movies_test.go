package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql" // MySQL driver import
	"github.com/koshikovme/assignment3-greenlight/internal/data"
	"github.com/koshikovme/assignment3-greenlight/internal/jsonlog"
	"github.com/koshikovme/assignment3-greenlight/internal/mailer"
)

// setupTestDB initializes a new test database connection.
func setupTestDB2(t *testing.T) *sql.DB {
	t.Helper()

	// Connect to the test database.
	db, err := sql.Open("mysql", "admin:1234@tcp(localhost:3306)/electronicsstore")
	if err != nil {
		t.Fatal(err)
	}

	// Clean the database before running tests.
	_, err = db.Exec("TRUNCATE TABLE movies")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

// newTestApplication creates a new application struct with mocked dependencies for testing.
func newTestApplication2(t *testing.T) *application {
	t.Helper()

	db := setupTestDB2(t)

	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	return &application{
		config: config{},
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.Mailer{},
	}
}

// TestCreateMovieHandler tests the createMovieHandler.
func TestCreateMovieHandler(t *testing.T) {
	app := newTestApplication2(t)

	movie := map[string]interface{}{
		"title":   "Test Movie",
		"year":    2023,
		"runtime": 120,
		"genres":  []string{"Drama"},
	}

	body, err := json.Marshal(movie)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/movies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.createMovieHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created but got %d", rr.Code)
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if response["movie"].(map[string]interface{})["title"] != "Test Movie" {
		t.Errorf("expected title to be 'Test Movie' but got '%s'", response["movie"].(map[string]interface{})["title"])
	}
}

// Add other test functions similarly...



// TestShowMovieHandler tests the showMovieHandler.
func TestShowMovieHandler(t *testing.T) {
	app := newTestApplication2(t)

	// Insert a test movie into the database.
	movie := &data.Movie{
		Title:   "Test Movie",
		Year:    2023,
		Runtime: 120,
		Genres:  []string{"Drama"},
	}

	err := app.models.Movies.Insert(movie)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/movies/"+strconv.Itoa(int(movie.ID)), nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.showMovieHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK but got %d", rr.Code)
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if response["movie"].(map[string]interface{})["title"] != "Test Movie" {
		t.Errorf("expected title to be 'Test Movie' but got '%s'", response["movie"].(map[string]interface{})["title"])
	}
}

// TestUpdateMovieHandler tests the updateMovieHandler.
func TestUpdateMovieHandler(t *testing.T) {
	app := newTestApplication2(t)

	// Insert a test movie into the database.
	movie := &data.Movie{
		Title:   "Old Title",
		Year:    2022,
		Runtime: 100,
		Genres:  []string{"Comedy"},
	}

	err := app.models.Movies.Insert(movie)
	if err != nil {
		t.Fatal(err)
	}

	updatedMovie := map[string]interface{}{
		"title": "New Title",
	}

	body, err := json.Marshal(updatedMovie)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/v1/movies/"+strconv.Itoa(int(movie.ID)), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.updateMovieHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK but got %d", rr.Code)
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if response["movie"].(map[string]interface{})["title"] != "New Title" {
		t.Errorf("expected title to be 'New Title' but got '%s'", response["movie"].(map[string]interface{})["title"])
	}
}

// TestDeleteMovieHandler tests the deleteMovieHandler.
func TestDeleteMovieHandler(t *testing.T) {
	app := newTestApplication2(t)

	// Insert a test movie into the database.
	movie := &data.Movie{
		Title:   "Test Movie",
		Year:    2023,
		Runtime: 120,
		Genres:  []string{"Drama"},
	}

	err := app.models.Movies.Insert(movie)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/v1/movies/"+strconv.Itoa(int(movie.ID)), nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.deleteMovieHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK but got %d", rr.Code)
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if response["message"] != "movie successfully deleted" {
		t.Errorf("expected message to be 'movie successfully deleted' but got '%s'", response["message"])
	}
}

// TestListMoviesHandler tests the listMoviesHandler.
func TestListMoviesHandler(t *testing.T) {
	app := newTestApplication2(t)

	// Insert a few test movies into the database.
	movies := []*data.Movie{
		{
			Title:   "Movie 1",
			Year:    2021,
			Runtime: 100,
			Genres:  []string{"Drama"},
		},
		{
			Title:   "Movie 2",
			Year:    2022,
			Runtime: 110,
			Genres:  []string{"Comedy"},
		},
	}

	for _, movie := range movies {
		err := app.models.Movies.Insert(movie)
		if err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/movies", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.listMoviesHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK but got %d", rr.Code)
	}

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	moviesResponse := response["movies"].([]interface{})
	if len(moviesResponse) != len(movies) {
		t.Errorf("expected %d movies but got %d", len(movies), len(moviesResponse))
	}
}
