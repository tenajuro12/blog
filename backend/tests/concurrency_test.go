package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"blog-app/handlers"
	"blog-app/models"
	"blog-app/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupSilentDB(t *testing.T) {
	var err error
	models.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := models.DB.AutoMigrate(
		&models.User{}, &models.Post{}, &models.Comment{},
		&models.Like{}, &models.Follow{}, &models.Tag{}, &models.Bookmark{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	utils.InitJWT("test-secret")
}

func TestConcurrency_SimultaneousRegistrations(t *testing.T) {
	setupSilentDB(t)
	db := models.DB

	const n = 10
	results := make([]int, n)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			models.DB = db
			rr := postJSON(handlers.Register, "/api/auth/register", map[string]string{
				"username": fmt.Sprintf("user%d", i),
				"email":    fmt.Sprintf("user%d@example.com", i),
				"password": "password123",
			})
			mu.Lock()
			results[i] = rr.Code
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	// No 500 errors — server must not panic
	for i, code := range results {
		if code == http.StatusInternalServerError {
			t.Errorf("goroutine %d: got 500 — server panic", i)
		}
	}

	// Count successes
	successes := 0
	for _, code := range results {
		if code == http.StatusCreated {
			successes++
		}
	}
	t.Logf("Concurrent registrations: %d/%d succeeded (rest = lock contention)", successes, n)

	// At least some must succeed
	if successes == 0 {
		t.Error("expected at least 1 successful registration")
	}
}

func TestConcurrency_DoubleLike(t *testing.T) {
	setupSilentDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")
	db := models.DB

	var wg sync.WaitGroup
	results := make([]int, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			models.DB = db
			req := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/like", nil, userID)
			req.URL.Path = "/api/posts/" + post.Slug + "/like"
			req.Method = http.MethodPost
			rr := httptest.NewRecorder()
			handlers.LikePost(rr, req)
			results[i] = rr.Code
		}(i)
	}
	wg.Wait()

	for i, code := range results {
		if code == http.StatusInternalServerError {
			t.Errorf("goroutine %d: double like caused 500", i)
		}
	}

	// DB must have at most 1 like record
	var count int64
	models.DB.Model(&models.Like{}).Where("user_id = ? AND post_id = ?", userID, post.ID).Count(&count)
	if count > 1 {
		t.Errorf("expected max 1 like in DB, got %d — UNIQUE constraint not working", count)
	}
	t.Logf("Double like: results=%v DB_count=%d", results, count)
}

func TestConcurrency_DoubleBookmark(t *testing.T) {
	setupSilentDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")
	post := createPost(t, userID, "Post", "Body content for the post here.")
	db := models.DB

	var wg sync.WaitGroup
	results := make([]int, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			models.DB = db
			req := authRequest(http.MethodPost, "/api/posts/"+post.Slug+"/bookmark", nil, userID)
			req.URL.Path = "/api/posts/" + post.Slug + "/bookmark"
			req.Method = http.MethodPost
			rr := httptest.NewRecorder()
			handlers.BookmarkPost(rr, req)
			results[i] = rr.Code
		}(i)
	}
	wg.Wait()

	for i, code := range results {
		if code == http.StatusInternalServerError {
			t.Errorf("goroutine %d: double bookmark caused 500", i)
		}
	}

	var count int64
	models.DB.Model(&models.Bookmark{}).Where("user_id = ? AND post_id = ?", userID, post.ID).Count(&count)
	if count > 1 {
		t.Errorf("expected max 1 bookmark in DB, got %d", count)
	}
	t.Logf("Double bookmark: results=%v DB_count=%d", results, count)
}

func TestConcurrency_DoubleFollow(t *testing.T) {
	setupSilentDB(t)
	followerID, _ := createUser(t, "alice", "alice@example.com", "pass")
	targetID, _ := createUser(t, "bob", "bob@example.com", "pass")
	db := models.DB

	var wg sync.WaitGroup
	results := make([]int, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			models.DB = db
			req := authRequest(http.MethodPost, fmt.Sprintf("/api/users/%d/follow", targetID), nil, followerID)
			req.URL.Path = fmt.Sprintf("/api/users/%d/follow", targetID)
			req.Method = http.MethodPost
			rr := httptest.NewRecorder()
			handlers.FollowUser(rr, req)
			results[i] = rr.Code
		}(i)
	}
	wg.Wait()

	for i, code := range results {
		if code == http.StatusInternalServerError {
			t.Errorf("goroutine %d: double follow caused 500", i)
		}
	}

	var count int64
	models.DB.Model(&models.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, targetID).Count(&count)
	if count > 1 {
		t.Errorf("expected max 1 follow in DB, got %d", count)
	}
	t.Logf("Double follow: results=%v DB_count=%d", results, count)
}
