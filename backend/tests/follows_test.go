package tests

import (
	"blog-app/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ── FollowUser ────────────────────────────────────────────────────────────────

func TestFollowUser_Success(t *testing.T) {
	setupTestDB(t)
	followerID, _ := createUser(t, "alice", "alice@example.com", "pass")
	targetID, _ := createUser(t, "bob", "bob@example.com", "pass")

	req := authRequest(http.MethodPost, "/api/users/"+itoa(targetID)+"/follow", nil, followerID)
	req.URL.Path = "/api/users/" + itoa(targetID) + "/follow"
	req.Method = http.MethodPost

	rr := httptest.NewRecorder()
	handlers.FollowUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestFollowUser_TargetNotFound(t *testing.T) {
	setupTestDB(t)
	followerID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPost, "/api/users/9999/follow", nil, followerID)
	req.URL.Path = "/api/users/9999/follow"
	req.Method = http.MethodPost

	rr := httptest.NewRecorder()
	handlers.FollowUser(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestFollowUser_CannotFollowSelf(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodPost, "/api/users/"+itoa(userID)+"/follow", nil, userID)
	req.URL.Path = "/api/users/" + itoa(userID) + "/follow"
	req.Method = http.MethodPost

	rr := httptest.NewRecorder()
	handlers.FollowUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 when following self, got %d", rr.Code)
	}
}

func TestUnfollowUser_Success(t *testing.T) {
	setupTestDB(t)
	followerID, _ := createUser(t, "alice", "alice@example.com", "pass")
	targetID, _ := createUser(t, "bob", "bob@example.com", "pass")

	// follow first
	req1 := authRequest(http.MethodPost, "/api/users/"+itoa(targetID)+"/follow", nil, followerID)
	req1.URL.Path = "/api/users/" + itoa(targetID) + "/follow"
	req1.Method = http.MethodPost
	handlers.FollowUser(httptest.NewRecorder(), req1)

	// then unfollow
	req2 := authRequest(http.MethodDelete, "/api/users/"+itoa(targetID)+"/follow", nil, followerID)
	req2.URL.Path = "/api/users/" + itoa(targetID) + "/follow"
	req2.Method = http.MethodDelete

	rr := httptest.NewRecorder()
	handlers.FollowUser(rr, req2)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestFollowUser_Duplicate(t *testing.T) {
	setupTestDB(t)
	followerID, _ := createUser(t, "alice", "alice@example.com", "pass")
	targetID, _ := createUser(t, "bob", "bob@example.com", "pass")

	// follow once
	req1 := authRequest(http.MethodPost, "/api/users/"+itoa(targetID)+"/follow", nil, followerID)
	req1.URL.Path = "/api/users/" + itoa(targetID) + "/follow"
	req1.Method = http.MethodPost
	handlers.FollowUser(httptest.NewRecorder(), req1)

	// follow again — should return 200 not error
	req2 := authRequest(http.MethodPost, "/api/users/"+itoa(targetID)+"/follow", nil, followerID)
	req2.URL.Path = "/api/users/" + itoa(targetID) + "/follow"
	req2.Method = http.MethodPost
	rr := httptest.NewRecorder()
	handlers.FollowUser(rr, req2)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 on duplicate follow, got %d", rr.Code)
	}
}

// ── GetFollowers ──────────────────────────────────────────────────────────────

func TestGetFollowers_Empty(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := httptest.NewRequest(http.MethodGet, "/api/users/"+itoa(userID)+"/followers", nil)
	req.URL.Path = "/api/users/" + itoa(userID) + "/followers"
	rr := httptest.NewRecorder()
	handlers.GetFollowers(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var users []interface{}
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 0 {
		t.Errorf("expected 0 followers, got %d", len(users))
	}
}

func TestGetFollowers_WithData(t *testing.T) {
	setupTestDB(t)
	followerID, _ := createUser(t, "alice", "alice@example.com", "pass")
	targetID, _ := createUser(t, "bob", "bob@example.com", "pass")

	// alice follows bob
	req1 := authRequest(http.MethodPost, "/api/users/"+itoa(targetID)+"/follow", nil, followerID)
	req1.URL.Path = "/api/users/" + itoa(targetID) + "/follow"
	req1.Method = http.MethodPost
	handlers.FollowUser(httptest.NewRecorder(), req1)

	// check bob's followers
	req2 := httptest.NewRequest(http.MethodGet, "/api/users/"+itoa(targetID)+"/followers", nil)
	req2.URL.Path = "/api/users/" + itoa(targetID) + "/followers"
	rr := httptest.NewRecorder()
	handlers.GetFollowers(rr, req2)

	var users []interface{}
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 {
		t.Errorf("expected 1 follower, got %d", len(users))
	}
}

// ── GetFollowing ──────────────────────────────────────────────────────────────

func TestGetFollowing_Empty(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := httptest.NewRequest(http.MethodGet, "/api/users/"+itoa(userID)+"/following", nil)
	req.URL.Path = "/api/users/" + itoa(userID) + "/following"
	rr := httptest.NewRecorder()
	handlers.GetFollowing(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetFollowing_WithData(t *testing.T) {
	setupTestDB(t)
	followerID, _ := createUser(t, "alice", "alice@example.com", "pass")
	targetID, _ := createUser(t, "bob", "bob@example.com", "pass")

	// alice follows bob
	req1 := authRequest(http.MethodPost, "/api/users/"+itoa(targetID)+"/follow", nil, followerID)
	req1.URL.Path = "/api/users/" + itoa(targetID) + "/follow"
	req1.Method = http.MethodPost
	handlers.FollowUser(httptest.NewRecorder(), req1)

	// check alice's following list
	req2 := httptest.NewRequest(http.MethodGet, "/api/users/"+itoa(followerID)+"/following", nil)
	req2.URL.Path = "/api/users/" + itoa(followerID) + "/following"
	rr := httptest.NewRecorder()
	handlers.GetFollowing(rr, req2)

	var users []interface{}
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 {
		t.Errorf("expected 1 following, got %d", len(users))
	}
}

// ── GetFeed ───────────────────────────────────────────────────────────────────

func TestGetFeed_Empty(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := authRequest(http.MethodGet, "/api/feed", nil, userID)
	rr := httptest.NewRecorder()
	handlers.GetFeed(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var posts []interface{}
	json.NewDecoder(rr.Body).Decode(&posts)
	if len(posts) != 0 {
		t.Errorf("expected 0 posts in empty feed, got %d", len(posts))
	}
}

func TestGetFeed_WithFollowedPosts(t *testing.T) {
	setupTestDB(t)
	followerID, _ := createUser(t, "alice", "alice@example.com", "pass")
	authorID, _ := createUser(t, "bob", "bob@example.com", "pass")

	// alice follows bob
	req1 := authRequest(http.MethodPost, "/api/users/"+itoa(authorID)+"/follow", nil, followerID)
	req1.URL.Path = "/api/users/" + itoa(authorID) + "/follow"
	req1.Method = http.MethodPost
	handlers.FollowUser(httptest.NewRecorder(), req1)

	// bob creates a post
	createPost(t, authorID, "Bob's Post", "Content from bob for the feed test.")

	// alice checks feed
	req2 := authRequest(http.MethodGet, "/api/feed", nil, followerID)
	rr := httptest.NewRecorder()
	handlers.GetFeed(rr, req2)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var posts []interface{}
	json.NewDecoder(rr.Body).Decode(&posts)
	if len(posts) != 1 {
		t.Errorf("expected 1 post in feed, got %d", len(posts))
	}
}
