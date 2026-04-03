package tests

import (
	"blog-app/handlers"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ── GetProfile ────────────────────────────────────────────────────────────────

func TestGetProfile_Success(t *testing.T) {
	setupTestDB(t)
	createUser(t, "alice", "alice@example.com", "pass")

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/alice", nil)
	req.URL.Path = "/api/profiles/alice"
	rr := httptest.NewRecorder()
	handlers.GetProfile(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/nobody", nil)
	req.URL.Path = "/api/profiles/nobody"
	rr := httptest.NewRecorder()
	handlers.GetProfile(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

// ── UpdateProfile ─────────────────────────────────────────────────────────────

func TestUpdateProfile_Success(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPut, "/api/profile", map[string]string{
		"username": "alice",
		"bio":      "I love Go",
		"avatar":   "",
	}, userID)

	rr := httptest.NewRecorder()
	handlers.UpdateProfile(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestUpdateProfile_UsernameConflict(t *testing.T) {
	setupTestDB(t)
	createUser(t, "bob", "bob@example.com", "pass")
	aliceID, _ := createUser(t, "alice", "alice@example.com", "pass")

	// try to change alice's username to bob (already taken)
	req := authRequest(http.MethodPut, "/api/profile", map[string]string{
		"username": "bob",
		"bio":      "",
		"avatar":   "",
	}, aliceID)

	rr := httptest.NewRecorder()
	handlers.UpdateProfile(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", rr.Code)
	}
}

func TestUpdateProfile_InvalidBody(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPut, "/api/profile", "not-json", userID)
	rr := httptest.NewRecorder()
	handlers.UpdateProfile(rr, req)

	// invalid JSON — expect 400
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

// ── GetUserPosts ──────────────────────────────────────────────────────────────

func TestGetUserPosts_Empty(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/"+itoa(userID)+"/posts", nil)
	req.URL.Path = "/api/profiles/" + itoa(userID) + "/posts"
	rr := httptest.NewRecorder()
	handlers.GetUserPosts(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetUserPosts_WithPosts(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	createPost(t, userID, "Post 1", "Body one content here for testing.")
	createPost(t, userID, "Post 2", "Body two content here for testing.")

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/"+itoa(userID)+"/posts", nil)
	req.URL.Path = "/api/profiles/" + itoa(userID) + "/posts"
	rr := httptest.NewRecorder()
	handlers.GetUserPosts(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var posts []interface{}
	json.NewDecoder(rr.Body).Decode(&posts)
	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}
}

func TestGetUserPosts_InvalidID(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/abc/posts", nil)
	req.URL.Path = "/api/profiles/abc/posts"
	rr := httptest.NewRecorder()
	handlers.GetUserPosts(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

// helper: uint to string
func itoa(id uint) string {
	return fmt.Sprintf("%d", id)
}
