package tests

import (
	"blog-app/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ── LikePost ──────────────────────────────────────────────────────────────────

func TestLikePost_Success(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")

	req := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/like", nil, userID)
	req.URL.Path = "/api/posts/" + post.Slug + "/like"
	req.Method = http.MethodPost

	rr := httptest.NewRecorder()
	handlers.LikePost(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestLikePost_NotFound(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPost, "/api/posts/no-such-slug/like", nil, userID)
	req.URL.Path = "/api/posts/no-such-slug/like"
	req.Method = http.MethodPost

	rr := httptest.NewRecorder()
	handlers.LikePost(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestUnlikePost_Success(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")

	// like first
	likeReq := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/like", nil, userID)
	likeReq.URL.Path = "/api/posts/" + post.Slug + "/like"
	likeReq.Method = http.MethodPost
	handlers.LikePost(httptest.NewRecorder(), likeReq)

	// then unlike
	req := authRequest(http.MethodDelete, "/api/posts/"+post.Slug+"/like", nil, userID)
	req.URL.Path = "/api/posts/" + post.Slug + "/like"
	req.Method = http.MethodDelete

	rr := httptest.NewRecorder()
	handlers.LikePost(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestLikePost_Duplicate(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")

	// like once
	req1 := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/like", nil, userID)
	req1.URL.Path = "/api/posts/" + post.Slug + "/like"
	req1.Method = http.MethodPost
	handlers.LikePost(httptest.NewRecorder(), req1)

	// like again — should return 200, not error
	req2 := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/like", nil, userID)
	req2.URL.Path = "/api/posts/" + post.Slug + "/like"
	req2.Method = http.MethodPost
	rr := httptest.NewRecorder()
	handlers.LikePost(rr, req2)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 on duplicate like, got %d", rr.Code)
	}
}

// ── GetLikes ──────────────────────────────────────────────────────────────────

func TestGetLikes_Zero(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")

	req := httptest.NewRequest(http.MethodGet, "/api/posts/"+post.Slug+"/likes", nil)
	req.URL.Path = "/api/posts/" + post.Slug + "/likes"
	rr := httptest.NewRecorder()
	handlers.GetLikes(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp map[string]int64
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp["likes"] != 0 {
		t.Errorf("expected 0 likes, got %d", resp["likes"])
	}
}

func TestGetLikes_AfterLike(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")

	// like the post
	likeReq := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/like", nil, userID)
	likeReq.URL.Path = "/api/posts/" + post.Slug + "/like"
	likeReq.Method = http.MethodPost
	handlers.LikePost(httptest.NewRecorder(), likeReq)

	// check count
	req := httptest.NewRequest(http.MethodGet, "/api/posts/"+post.Slug+"/likes", nil)
	req.URL.Path = "/api/posts/" + post.Slug + "/likes"
	rr := httptest.NewRecorder()
	handlers.GetLikes(rr, req)

	var resp map[string]int64
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp["likes"] != 1 {
		t.Errorf("expected 1 like, got %d", resp["likes"])
	}
}

func TestGetLikes_PostNotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/posts/no-such/likes", nil)
	req.URL.Path = "/api/posts/no-such/likes"
	rr := httptest.NewRecorder()
	handlers.GetLikes(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}
