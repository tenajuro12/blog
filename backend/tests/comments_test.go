package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog-app/handlers"
	"blog-app/models"
)

func createPost(t *testing.T, userID uint, title, body string) models.Post {
	req := authRequest(http.MethodPost, "/api/posts", map[string]string{
		"title": title, "body": body,
	}, userID)
	rr := httptest.NewRecorder()
	handlers.CreatePost(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("createPost failed: %s", rr.Body.String())
	}
	var post models.Post
	json.NewDecoder(rr.Body).Decode(&post)
	return post
}

// ── Create Comment ────────────────────────────────────────────────────────────

func TestCreateComment_Success(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Test Post", "Post body content here.")

	req := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/comments", map[string]string{
		"body": "Great post!",
	}, userID)
	req.URL.Path = "/api/posts/" + post.Slug + "/comments"

	rr := httptest.NewRecorder()
	handlers.CreateComment(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", rr.Code, rr.Body.String())
	}

	var comment models.Comment
	json.NewDecoder(rr.Body).Decode(&comment)
	if comment.Body != "Great post!" {
		t.Errorf("expected 'Great post!', got '%s'", comment.Body)
	}
}

func TestCreateComment_EmptyBody(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Test Post", "Post body content here.")

	req := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/comments", map[string]string{
		"body": "",
	}, userID)
	req.URL.Path = "/api/posts/" + post.Slug + "/comments"

	rr := httptest.NewRecorder()
	handlers.CreateComment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCreateComment_PostNotFound(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPost, "/api/posts/nonexistent-slug/comments", map[string]string{
		"body": "comment",
	}, userID)
	req.URL.Path = "/api/posts/nonexistent-slug/comments"

	rr := httptest.NewRecorder()
	handlers.CreateComment(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

// ── Delete Comment ────────────────────────────────────────────────────────────

func TestDeleteComment_NotOwner(t *testing.T) {
	setupTestDB(t)
	ownerID, _ := createUser(t, "owner", "owner@example.com", "pass")
	otherID, _ := createUser(t, "other", "other@example.com", "pass")
	post := createPost(t, ownerID, "Post", "Body of the post content.")

	// create comment as owner
	req := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/comments", map[string]string{
		"body": "my comment",
	}, ownerID)
	req.URL.Path = "/api/posts/" + post.Slug + "/comments"
	rr := httptest.NewRecorder()
	handlers.CreateComment(rr, req)

	var comment models.Comment
	json.NewDecoder(rr.Body).Decode(&comment)

	// try delete as other user
	delReq := authRequest(http.MethodDelete, "/api/posts/"+post.Slug+"/comments/1", nil, otherID)
	delReq.URL.Path = "/api/posts/" + post.Slug + "/comments/1"
	delRR := httptest.NewRecorder()
	handlers.DeleteComment(delRR, delReq)

	if delRR.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", delRR.Code)
	}
}
