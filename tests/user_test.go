package tests

import (
	"backend/internal/users/delivery"
	"backend/internal/users/repository"
	"backend/internal/users/usecase"
	"backend/middleware"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// SetupRouter sets up the Gin router for tests
func setupRouter(mockDB *sql.DB) *gin.Engine {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Setup User Repository and Usecase
	userRepo := repository.NewUserRepository(mockDB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := delivery.NewUserHandler(userUsecase)

	// Routes with authentication
	router.POST("/api/login", userHandler.Login)
	router.POST("/users", userHandler.CreateUser)
	auth := router.Group("/api")
	auth.Use(middleware.JWTMiddleware(mockDB))
	{
		auth.GET("/users", userHandler.GetAllUsers)
		auth.POST("/users/delete", userHandler.DeleteUser)
	}

	return router
}

// TestCreateUser tests creating a new user
func TestCreateUser(t *testing.T) {
	// Create mock database connection and mock statements
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Mock database query for creating a user
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("john_doe", "john_doe@example.com", "hashed_password", 1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

	// Setup Gin router with the mock DB
	router := setupRouter(db)

	// Prepare request payload for creating a user
	reqBody := map[string]interface{}{
		"username":  "password123",
		"email":     "password123@example.com",
		"password":  "password123",
		"role_id":   121,
		"client_id": 121,
	}
	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/api/users", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	// Perform the request
	resp := performRequest(router, req)

	// Assert response code and message
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User successfully created")
}

// TestGetAllUsers tests getting all users
func TestGetAllUsers(t *testing.T) {
	// Create mock database connection and mock statements
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Mock database query for getting users
	mock.ExpectQuery("SELECT user_id, username, email").WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "email"}).
		AddRow(1, "password123", "password123@example.com").
		AddRow(2, "password123", "password123@example.com"))

	// Setup Gin router with the mock DB
	router := setupRouter(db)

	// Perform the request
	req, err := http.NewRequest("GET", "/api/users", nil)
	assert.NoError(t, err)

	// Perform the request
	resp := performRequest(router, req)

	// Assert response code and message
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "password123")
	assert.Contains(t, resp.Body.String(), "password123")
}

// TestDeleteUser_Success tests deleting a user by ID
func TestDeleteUser_Success(t *testing.T) {
	// Create mock database connection and mock statements
	db, mock, err := sqlmock.New() // Create mock database
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close() // Ensure the database is closed after the test

	// Mock database query for getting user
	mock.ExpectQuery("SELECT user_id, username, email").WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "email"}).
		AddRow(1, "john_doe", "password123@example.com"))

	// Mock database query for deleting user
	mock.ExpectExec("DELETE FROM users WHERE user_id = $1").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock database query for adding token to blacklist
	mock.ExpectExec("INSERT INTO blacklisted_tokens (token) VALUES ($1)").
		WithArgs("valid_token_string").WillReturnResult(sqlmock.NewResult(1, 1))

	// Setup Gin router with the mock DB
	router := setupRouter(db) // Pass *sql.DB here

	// Prepare request payload for deleting a user
	reqBody := map[string]int{"id": 1}
	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("DELETE", "/api/users/delete", bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "Bearer valid_token_string")
	assert.NoError(t, err)

	// Perform the request
	resp := performRequest(router, req)

	// Assert response code and message
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User successfully deleted")
}

// TestDeleteUser_UserNotFound tests if user is not found
func TestDeleteUser_UserNotFound(t *testing.T) {
	// Create mock database connection and mock statements
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Mock database query for getting user (return no rows)
	mock.ExpectQuery("SELECT user_id, username, email").WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "email"}))

	// Setup Gin router with the mock DB
	router := setupRouter(db)

	// Prepare request payload for deleting a user
	reqBody := map[string]int{"id": 999} // non-existent user ID
	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("DELETE", "/api/users/delete", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	// Perform the request
	resp := performRequest(router, req)

	// Assert response code and message
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "User not found")
}

// Helper function to perform HTTP request
func performRequest(router http.Handler, req *http.Request) *httptest.ResponseRecorder {
	// Record the response
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
