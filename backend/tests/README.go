package tests

// This file documents the test setup strategy.
//
// LOCAL development:
//   Uses SQLite in-memory DB — no Postgres needed.
//   Run: cd backend && go test ./tests/... -v
//
// CI (GitHub Actions):
//   Uses real PostgreSQL via Docker service.
//   The DATABASE_URL env var is set in the workflow.
//
// The setupTestDB() helper in auth_test.go uses SQLite (:memory:)
// which is fast and requires no external setup.
