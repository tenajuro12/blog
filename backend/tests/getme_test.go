package tests

import (
	"blog-app/handlers"
	"blog-app/middleware"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ── GetMe ─────────────────────────────────────────────────────────────────────

func TestGetMe_Success(t *testing.T) {
	setupTestDB(t)
	userID, _ := createUser(t, "alice", "alice@example.com", "pass")

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handlers.GetMe(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — %s", rr.Code, rr.Body.String())
	}
}

func TestGetMe_UserNotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(9999))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handlers.GetMe(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}
