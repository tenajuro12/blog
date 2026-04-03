package tests

import (
	"blog-app/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ── BookmarkPost ──────────────────────────────────────────────────────────────

func TestBookmarkPost_Success(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")

	req := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/bookmark", nil, userID)
	req.URL.Path = "/api/posts/" + post.Slug + "/bookmark"
	req.Method = http.MethodPost

	rr := httptest.NewRecorder()
	handlers.BookmarkPost(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestBookmarkPost_NotFound(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPost, "/api/posts/no-slug/bookmark", nil, userID)
	req.URL.Path = "/api/posts/no-slug/bookmark"
	req.Method = http.MethodPost

	rr := httptest.NewRecorder()
	handlers.BookmarkPost(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestBookmarkPost_Duplicate(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")

	// bookmark once
	req1 := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/bookmark", nil, userID)
	req1.URL.Path = "/api/posts/" + post.Slug + "/bookmark"
	req1.Method = http.MethodPost
	handlers.BookmarkPost(httptest.NewRecorder(), req1)

	// bookmark again — should return 200 not error
	req2 := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/bookmark", nil, userID)
	req2.URL.Path = "/api/posts/" + post.Slug + "/bookmark"
	req2.Method = http.MethodPost
	rr := httptest.NewRecorder()
	handlers.BookmarkPost(rr, req2)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 on duplicate bookmark, got %d", rr.Code)
	}
}

func TestRemoveBookmark_Success(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")

	// bookmark first
	req1 := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/bookmark", nil, userID)
	req1.URL.Path = "/api/posts/" + post.Slug + "/bookmark"
	req1.Method = http.MethodPost
	handlers.BookmarkPost(httptest.NewRecorder(), req1)

	// then remove
	req2 := authRequest(http.MethodDelete, "/api/posts/"+post.Slug+"/bookmark", nil, userID)
	req2.URL.Path = "/api/posts/" + post.Slug + "/bookmark"
	req2.Method = http.MethodDelete

	rr := httptest.NewRecorder()
	handlers.BookmarkPost(rr, req2)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}
}

// ── GetBookmarks ──────────────────────────────────────────────────────────────

func TestGetBookmarks_Empty(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodGet, "/api/bookmarks", nil, userID)
	rr := httptest.NewRecorder()
	handlers.GetBookmarks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var posts []interface{}
	json.NewDecoder(rr.Body).Decode(&posts)
	if len(posts) != 0 {
		t.Errorf("expected 0 bookmarks, got %d", len(posts))
	}
}

func TestGetBookmarks_WithData(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")

	// bookmark it
	req1 := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/bookmark", nil, userID)
	req1.URL.Path = "/api/posts/" + post.Slug + "/bookmark"
	req1.Method = http.MethodPost
	handlers.BookmarkPost(httptest.NewRecorder(), req1)

	// get bookmarks
	req2 := authRequest(http.MethodGet, "/api/bookmarks", nil, userID)
	rr := httptest.NewRecorder()
	handlers.GetBookmarks(rr, req2)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var posts []interface{}
	json.NewDecoder(rr.Body).Decode(&posts)
	if len(posts) != 1 {
		t.Errorf("expected 1 bookmark, got %d", len(posts))
	}
}
