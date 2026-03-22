package tests

import (
	"testing"

	"blog-app/utils"
)

// ── JWT ───────────────────────────────────────────────────────────────────────

func TestJWT_GenerateAndParse(t *testing.T) {
	utils.InitJWT("test-secret")

	token, err := utils.GenerateToken(42)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("expected userID 42, got %d", claims.UserID)
	}
}

func TestJWT_InvalidToken(t *testing.T) {
	utils.InitJWT("test-secret")

	_, err := utils.ParseToken("this.is.not.valid")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestJWT_WrongSecret(t *testing.T) {
	utils.InitJWT("secret-one")
	token, _ := utils.GenerateToken(1)

	utils.InitJWT("secret-two")
	_, err := utils.ParseToken(token)
	if err == nil {
		t.Error("expected error when secret changed")
	}
}

// ── Slug ──────────────────────────────────────────────────────────────────────

func TestSlugify_Basic(t *testing.T) {
	slug := utils.Slugify("Hello World")
	if slug == "" {
		t.Fatal("expected non-empty slug")
	}
	// should be lowercase and contain "hello-world"
	for _, ch := range slug {
		if ch >= 'A' && ch <= 'Z' {
			t.Errorf("slug should be lowercase, got: %s", slug)
		}
	}
}

func TestSlugify_SpecialChars(t *testing.T) {
	slug := utils.Slugify("Hello, World! How are you?")
	// should not contain commas, exclamation marks, or question marks
	forbidden := []rune{',', '!', '?', ' '}
	for _, f := range forbidden {
		for _, ch := range slug {
			if ch == f {
				t.Errorf("slug contains forbidden char '%c': %s", f, slug)
			}
		}
	}
}
