package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog-app/handlers"
	"blog-app/middleware"
	"blog-app/models"
)

// helper: register a user and return their ID + token
func createUser(t *testing.T, username, email, password string) (uint, string) {
	rr := postJSON(handlers.Register, "/api/auth/register", map[string]string{
		"username": username, "email": email, "password": password,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("createUser failed: %s", rr.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	user := resp["user"].(map[string]interface{})
	return uint(user["id"].(float64)), resp["token"].(string)
}

// helper: make authenticated request
func authRequest(method, path string, body interface{}, userID uint) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

// ── Create Post ───────────────────────────────────────────────────────────────

func TestCreatePost_Success(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPost, "/api/posts", map[string]string{
		"title": "Hello World",
		"body":  "This is my first post content here.",
		"tags":  "go,testing",
	}, userID)

	rr := httptest.NewRecorder()
	handlers.CreatePost(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", rr.Code, rr.Body.String())
	}

	var post models.Post
	json.NewDecoder(rr.Body).Decode(&post)
	if post.Title != "Hello World" {
		t.Errorf("expected title 'Hello World', got '%s'", post.Title)
	}
	if post.Slug == "" {
		t.Error("expected slug to be generated")
	}
}

func TestCreatePost_MissingTitle(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPost, "/api/posts", map[string]string{
		"body": "Body without title",
	}, userID)

	rr := httptest.NewRecorder()
	handlers.CreatePost(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

// ── List Posts ────────────────────────────────────────────────────────────────

func TestListPosts_Empty(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	rr := httptest.NewRecorder()
	handlers.ListPosts(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var posts []models.Post
	json.NewDecoder(rr.Body).Decode(&posts)
	if len(posts) != 0 {
		t.Errorf("expected empty list, got %d posts", len(posts))
	}
}

// ── Delete Post ───────────────────────────────────────────────────────────────

func TestDeletePost_NotOwner(t *testing.T) {
	setupTestDB(t)
	ownerID, _ := createUser(t, "owner", "owner@example.com", "pass")
	otherID, _ := createUser(t, "other", "other@example.com", "pass")

	// create post as owner
	req := authRequest(http.MethodPost, "/api/posts", map[string]string{
		"title": "Owner post",
		"body":  "Content of the owner post.",
	}, ownerID)
	rr := httptest.NewRecorder()
	handlers.CreatePost(rr, req)

	var post models.Post
	json.NewDecoder(rr.Body).Decode(&post)

	// try to delete as other user
	delReq := authRequest(http.MethodDelete, "/api/posts/"+post.Slug, nil, otherID)
	delReq.URL.Path = "/api/posts/" + post.Slug
	delRR := httptest.NewRecorder()
	handlers.DeletePost(delRR, delReq)

	if delRR.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", delRR.Code)
	}
}
