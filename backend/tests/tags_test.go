package tests

import (
	"blog-app/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ── GetTags ───────────────────────────────────────────────────────────────────

func TestGetTags_Empty(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/tags/", nil)
	req.URL.Path = "/api/tags/"
	rr := httptest.NewRecorder()
	handlers.GetTags(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}

	var tags []interface{}
	json.NewDecoder(rr.Body).Decode(&tags)
	if len(tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(tags))
	}
}

func TestGetTags_WithData(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	// create posts with tags
	req1 := authRequest(http.MethodPost, "/api/posts", map[string]string{
		"title": "Go Post",
		"body":  "Body content for the go post here.",
		"tags":  "go,testing",
	}, userID)
	rr1 := httptest.NewRecorder()
	handlers.CreatePost(rr1, req1)

	req2 := authRequest(http.MethodPost, "/api/posts", map[string]string{
		"title": "Another Go Post",
		"body":  "Body content for another post here.",
		"tags":  "go,web",
	}, userID)
	rr2 := httptest.NewRecorder()
	handlers.CreatePost(rr2, req2)

	// get tags
	req := httptest.NewRequest(http.MethodGet, "/api/tags/", nil)
	req.URL.Path = "/api/tags/"
	rr := httptest.NewRecorder()
	handlers.GetTags(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var tags []map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&tags)

	// should have 3 unique tags: go, testing, web
	if len(tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(tags))
	}

	// find "go" tag and check count = 2
	for _, tag := range tags {
		if tag["name"] == "go" {
			if tag["count"].(float64) != 2 {
				t.Errorf("expected go tag count 2, got %v", tag["count"])
			}
		}
	}
}

// ── GetPostsByTag ─────────────────────────────────────────────────────────────

func TestGetPostsByTag_Found(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req1 := authRequest(http.MethodPost, "/api/posts", map[string]string{
		"title": "Go Post",
		"body":  "Body content for the go post here.",
		"tags":  "go,testing",
	}, userID)
	handlers.CreatePost(httptest.NewRecorder(), req1)

	req := httptest.NewRequest(http.MethodGet, "/api/tags/go", nil)
	req.URL.Path = "/api/tags/go"
	rr := httptest.NewRecorder()
	handlers.GetPostsByTag(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}

	var posts []interface{}
	json.NewDecoder(rr.Body).Decode(&posts)
	if len(posts) != 1 {
		t.Errorf("expected 1 post for tag 'go', got %d", len(posts))
	}
}

func TestGetPostsByTag_NoResults(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/tags/python", nil)
	req.URL.Path = "/api/tags/python"
	rr := httptest.NewRecorder()
	handlers.GetPostsByTag(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var posts []interface{}
	json.NewDecoder(rr.Body).Decode(&posts)
	if len(posts) != 0 {
		t.Errorf("expected 0 posts for unknown tag, got %d", len(posts))
	}
}

func TestGetPostsByTag_MissingTagName(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/tags/", nil)
	req.URL.Path = "/api/tags/"
	rr := httptest.NewRecorder()
	handlers.GetPostsByTag(rr, req)

	// empty tag name → 400
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty tag name, got %d", rr.Code)
	}
}
