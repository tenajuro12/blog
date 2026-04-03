package tests

import (
	"blog-app/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ── GetPost ───────────────────────────────────────────────────────────────────

func TestGetPost_Success(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Hello World", "Some body content here.")

	req := httptest.NewRequest(http.MethodGet, "/api/posts/"+post.Slug, nil)
	req.URL.Path = "/api/posts/" + post.Slug
	rr := httptest.NewRecorder()
	handlers.GetPost(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestGetPost_NotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/posts/nonexistent-slug", nil)
	req.URL.Path = "/api/posts/nonexistent-slug"
	rr := httptest.NewRecorder()
	handlers.GetPost(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

// ── UpdatePost ────────────────────────────────────────────────────────────────

func TestUpdatePost_Success(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Original Title", "Original body content here.")

	req := authRequest(http.MethodPut, "/api/posts/"+post.Slug, map[string]string{
		"title": "Updated Title",
		"body":  "Updated body content here.",
	}, userID)
	req.URL.Path = "/api/posts/" + post.Slug

	rr := httptest.NewRecorder()
	handlers.UpdatePost(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestUpdatePost_NotOwner(t *testing.T) {
	setupTestDB(t)
	ownerID, _ := createUser(t, "owner", "owner@example.com", "pass")
	otherID, _ := createUser(t, "other", "other@example.com", "pass")
	post := createPost(t, ownerID, "Owner Post", "Owner body content here.")

	req := authRequest(http.MethodPut, "/api/posts/"+post.Slug, map[string]string{
		"title": "Hacked Title",
		"body":  "Hacked body.",
	}, otherID)
	req.URL.Path = "/api/posts/" + post.Slug

	rr := httptest.NewRecorder()
	handlers.UpdatePost(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rr.Code)
	}
}

func TestUpdatePost_NotFound(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPut, "/api/posts/fake-slug", map[string]string{
		"title": "Title",
		"body":  "Body.",
	}, userID)
	req.URL.Path = "/api/posts/fake-slug"

	rr := httptest.NewRecorder()
	handlers.UpdatePost(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

// ── ListComments ──────────────────────────────────────────────────────────────

func TestListComments_Empty(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content here for the post.")

	req := httptest.NewRequest(http.MethodGet, "/api/posts/"+post.Slug+"/comments", nil)
	req.URL.Path = "/api/posts/" + post.Slug + "/comments"
	rr := httptest.NewRecorder()
	handlers.ListComments(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestListComments_WithData(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content here for the post.")

	// add a comment first
	req := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/comments",
		map[string]string{"body": "Great post!"}, userID)
	req.URL.Path = "/api/posts/" + post.Slug + "/comments"
	rr := httptest.NewRecorder()
	handlers.CreateComment(rr, req)

	// now list
	req2 := httptest.NewRequest(http.MethodGet, "/api/posts/"+post.Slug+"/comments", nil)
	req2.URL.Path = "/api/posts/" + post.Slug + "/comments"
	rr2 := httptest.NewRecorder()
	handlers.ListComments(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr2.Code)
	}

	var comments []interface{}
	json.NewDecoder(rr2.Body).Decode(&comments)
	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}
}

func TestListComments_PostNotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/posts/no-such-post/comments", nil)
	req.URL.Path = "/api/posts/no-such-post/comments"
	rr := httptest.NewRecorder()
	handlers.ListComments(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}
