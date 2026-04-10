package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog-app/handlers"
	"blog-app/models"
)

func TestIntegration_RegisterThenCreatePost(t *testing.T) {
	setupTestDB(t)

	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Integration Post", "Body of integration test post content.")

	if post.AuthorID != userID {
		t.Errorf("expected author_id %d, got %d", userID, post.AuthorID)
	}
	if post.Slug == "" {
		t.Error("expected slug to be generated")
	}
}

func TestIntegration_CreatePostThenComment(t *testing.T) {
	setupTestDB(t)

	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post for comment", "Post body content for comment chain test.")

	req := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/comments",
		map[string]string{"body": "Great post!"}, userID)
	req.URL.Path = "/api/posts/" + post.Slug + "/comments"

	rr := httptest.NewRecorder()
	handlers.CreateComment(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d — %s", rr.Code, rr.Body.String())
	}

	var comment models.Comment
	json.NewDecoder(rr.Body).Decode(&comment)

	if comment.AuthorID != userID {
		t.Errorf("expected author_id %d, got %d", userID, comment.AuthorID)
	}
	if comment.PostID != post.ID {
		t.Errorf("expected post_id %d, got %d", post.ID, comment.PostID)
	}
}

func TestIntegration_CreatePostThenLike(t *testing.T) {
	setupTestDB(t)

	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post to like", "Body of post to like in integration test.")

	// like the post
	likeReq := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/like", nil, userID)
	likeReq.URL.Path = "/api/posts/" + post.Slug + "/like"
	likeReq.Method = http.MethodPost
	likeRR := httptest.NewRecorder()
	handlers.LikePost(likeRR, likeReq)

	if likeRR.Code != http.StatusCreated {
		t.Fatalf("like failed: %d — %s", likeRR.Code, likeRR.Body.String())
	}

	// verify count = 1
	countReq := httptest.NewRequest(http.MethodGet, "/api/posts/"+post.Slug+"/likes", nil)
	countReq.URL.Path = "/api/posts/" + post.Slug + "/likes"
	countRR := httptest.NewRecorder()
	handlers.GetLikes(countRR, countReq)

	var resp map[string]int64
	json.NewDecoder(countRR.Body).Decode(&resp)
	if resp["likes"] != 1 {
		t.Errorf("expected 1 like, got %d", resp["likes"])
	}
}

func TestIntegration_CreatePostThenBookmark(t *testing.T) {
	setupTestDB(t)

	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post to bookmark", "Body of post to bookmark in integration test.")

	// bookmark
	bmReq := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/bookmark", nil, userID)
	bmReq.URL.Path = "/api/posts/" + post.Slug + "/bookmark"
	bmReq.Method = http.MethodPost
	bmRR := httptest.NewRecorder()
	handlers.BookmarkPost(bmRR, bmReq)

	if bmRR.Code != http.StatusCreated {
		t.Fatalf("bookmark failed: %d — %s", bmRR.Code, bmRR.Body.String())
	}

	// get bookmarks — should contain the post
	listReq := authRequest(http.MethodGet, "/api/bookmarks", nil, userID)
	listRR := httptest.NewRecorder()
	handlers.GetBookmarks(listRR, listReq)

	var posts []models.Post
	json.NewDecoder(listRR.Body).Decode(&posts)

	if len(posts) != 1 {
		t.Errorf("expected 1 bookmarked post, got %d", len(posts))
	}
	if posts[0].ID != post.ID {
		t.Errorf("expected post ID %d in bookmarks, got %d", post.ID, posts[0].ID)
	}
}
