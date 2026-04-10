package tests

// edge_cases_test.go
// Tests boundary conditions: empty inputs, extremely large payloads,
// special characters, and injection-like inputs.

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"blog-app/handlers"
)

// ── Empty input fields ────────────────────────────────────────────────────────

func TestEdge_Register_EmptyUsername(t *testing.T) {
	setupTestDB(t)
	rr := postJSON(handlers.Register, "/api/auth/register", map[string]string{
		"username": "",
		"email":    "alice@example.com",
		"password": "password123",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty username, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestEdge_Register_EmptyEmail(t *testing.T) {
	setupTestDB(t)
	rr := postJSON(handlers.Register, "/api/auth/register", map[string]string{
		"username": "alice",
		"email":    "",
		"password": "password123",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty email, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestEdge_Register_EmptyPassword(t *testing.T) {
	setupTestDB(t)
	rr := postJSON(handlers.Register, "/api/auth/register", map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
		"password": "",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty password, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestEdge_CreatePost_EmptyTitle(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPost, "/api/posts", map[string]string{
		"title": "",
		"body":  "Some body content here.",
	}, userID)
	rr := httptest.NewRecorder()
	handlers.CreatePost(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty title, got %d", rr.Code)
	}
}

func TestEdge_CreatePost_EmptyBody(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPost, "/api/posts", map[string]string{
		"title": "Some Title",
		"body":  "",
	}, userID)
	rr := httptest.NewRecorder()
	handlers.CreatePost(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty body, got %d", rr.Code)
	}
}

func TestEdge_CreateComment_EmptyBody(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")

	req := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/comments",
		map[string]string{"body": ""}, userID)
	req.URL.Path = "/api/posts/" + post.Slug + "/comments"
	rr := httptest.NewRecorder()
	handlers.CreateComment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty comment body, got %d", rr.Code)
	}
}
