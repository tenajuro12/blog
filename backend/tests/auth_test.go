package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog-app/handlers"
	"blog-app/models"
	"blog-app/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB initializes an in-memory SQLite DB for testing
func setupTestDB(t *testing.T) {
	var err error
	models.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := models.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	utils.InitJWT("test-secret")
}

func postJSON(handler http.HandlerFunc, path string, body interface{}) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

// ── Register ──────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	setupTestDB(t)

	rr := postJSON(handlers.Register, "/api/auth/register", map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
		"password": "password123",
	})

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp["token"] == nil {
		t.Error("expected token in response")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	setupTestDB(t)

	payload := map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
		"password": "password123",
	}
	postJSON(handlers.Register, "/api/auth/register", payload)

	payload["username"] = "alice2"
	rr := postJSON(handlers.Register, "/api/auth/register", payload)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", rr.Code)
	}
}

func TestRegister_MissingFields(t *testing.T) {
	setupTestDB(t)

	rr := postJSON(handlers.Register, "/api/auth/register", map[string]string{
		"email": "noname@example.com",
	})

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

// ── Login ─────────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	setupTestDB(t)

	postJSON(handlers.Register, "/api/auth/register", map[string]string{
		"username": "bob",
		"email":    "bob@example.com",
		"password": "secret",
	})

	rr := postJSON(handlers.Login, "/api/auth/login", map[string]string{
		"email":    "bob@example.com",
		"password": "secret",
	})

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp["token"] == nil {
		t.Error("expected token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	setupTestDB(t)

	postJSON(handlers.Register, "/api/auth/register", map[string]string{
		"username": "bob",
		"email":    "bob@example.com",
		"password": "secret",
	})

	rr := postJSON(handlers.Login, "/api/auth/login", map[string]string{
		"email":    "bob@example.com",
		"password": "wrongpassword",
	})

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	setupTestDB(t)

	rr := postJSON(handlers.Login, "/api/auth/login", map[string]string{
		"email":    "ghost@example.com",
		"password": "whatever",
	})

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}
