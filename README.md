# BlogApp — QA Setup Guide

## Project Structure

```
blog-app/
├── backend/
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # JWT auth middleware
│   ├── models/               # GORM models (User, Post, Comment, Like, Follow, Tag, Bookmark)
│   ├── tests/                # All Go test files
│   │   ├── auth_test.go          # TC01–TC06: Register, Login, GetMe
│   │   ├── posts_test.go         # TC07–TC17: Posts CRUD + Authorization
│   │   ├── comments_test.go      # TC18–TC20: Comments
│   │   ├── utils_test.go         # TC21–TC23: JWT, Slugify
│   │   ├── getme_test.go         # TC33–TC34: GetMe endpoint
│   │   ├── posts_extra_test.go   # TC35–TC42: GetPost, UpdatePost, ListComments
│   │   ├── profile_test.go       # TC43–TC50: Profile management
│   │   ├── likes_test.go         # TC51–TC57: Likes
│   │   ├── bookmarks_test.go     # TC58–TC63: Bookmarks
│   │   ├── follows_test.go       # TC64–TC74: Follows / Feed
│   │   ├── tags_test.go          # TC75–TC79: Tags
│   │   ├── integration_test.go   # TC-INT-01–04: Cross-module chains
│   │   ├── edge_cases_test.go    # TC-EDGE-*: Empty inputs, large payloads, injection
│   │   └── concurrency_test.go   # TC-CONC-*: Race conditions, invalid behavior
│   ├── utils/
│   ├── main.go
│   └── go.mod
├── frontend/
│   ├── e2e/                  # Playwright E2E tests
│   │   ├── auth.spec.js          # TC-E2E-01,02,05,06,07
│   │   ├── posts.spec.js         # TC-E2E-03,08,09,10
│   │   ├── social.spec.js        # TC-E2E-04,11,12,13
│   │   └── helpers.js
│   ├── playwright.config.js
│   └── src/
├── postman/
│   ├── BlogApp.postman_collection.json   # Auth, Posts, Comments, Profile folders
│   └── BlogApp.postman_environment.json
├── k6_tests/                 # Performance tests
│   ├── normal_load.js            # 10 VUs, 2 min
│   ├── peak_load.js              # 50 VUs, 3 min
│   ├── spike_load.js             # 100 VUs, 1m 40s
│   └── endurance_load.js         # 5 VUs, 10 min
├── .github/workflows/
│   └── ci.yml                # 8-job GitHub Actions pipeline
└── docker-compose.yml
```

---

## 1. Run Unit Tests (Go)

Tests use **SQLite in-memory** — no PostgreSQL needed.

```bash
cd backend

# Run all tests
go test ./tests/... -v

# Run with coverage
go test ./tests/... -coverprofile=coverage.out -covermode=atomic \
  -coverpkg=blog-app/handlers,blog-app/utils

# View coverage summary
go tool cover -func coverage.out | tail -5

# View coverage in browser
go tool cover -html coverage.out -o coverage.html
```

### Run specific test categories

```bash
# Unit tests only
go test ./tests/... -v -run "^Test[^(Integration|Edge|Concurrency|InvalidBehavior)]"

# Integration tests only
go test ./tests/... -v -run "^TestIntegration_"

# Edge case tests only
go test ./tests/... -v -run "^TestEdge_"

# Concurrency + Invalid behavior (with race detector)
go test ./tests/... -v -run "^(TestConcurrency_|TestInvalidBehavior_)" -race -count=1
```

### Expected output
```
ok  blog-app/tests  6.016s  coverage: 88.1% of statements in blog-app/handlers, blog-app/utils
```

---

## 2. Run API Tests (Newman/Postman)

**Start backend first:**
```bash
cd backend
go run main.go
```

### Option A: Postman GUI
1. Open Postman
2. Import `postman/BlogApp.postman_collection.json`
3. Import `postman/BlogApp.postman_environment.json`
4. Run the collection

### Option B: Newman CLI
```bash
# Install Newman
npm install -g newman newman-reporter-htmlextra

# Run all folders
newman run postman/BlogApp.postman_collection.json \
  --environment postman/BlogApp.postman_environment.json \
  --env-var "baseUrl=http://localhost:8080" \
  --reporters cli,htmlextra \
  --reporter-htmlextra-export newman-report.html \
  --bail

# Run specific folder
newman run postman/BlogApp.postman_collection.json \
  --folder Auth \
  --environment postman/BlogApp.postman_environment.json \
  --env-var "baseUrl=http://localhost:8080"
```

### Expected output
```
30 assertions passed, 0 failed
Total run duration: ~1,768ms
```

---

## 3. Run E2E Tests (Playwright)

**Requires both backend and frontend running.**

```bash
# Terminal 1 — backend
cd backend && go run main.go

# Terminal 2 — frontend
cd frontend && npm start

# Terminal 3 — run E2E tests
cd frontend
npm install @playwright/test
npx playwright install chromium
npx playwright test

# Run specific file
npx playwright test e2e/auth.spec.js
npx playwright test e2e/posts.spec.js
npx playwright test e2e/social.spec.js

# With UI mode (interactive)
npx playwright test --ui

# View HTML report
npx playwright show-report
```

### Expected output
```
13 passed (49.1s)
```

---

## 4. Run Performance Tests (k6)

**Requires backend running.**

```bash
# Install k6 (Windows)
winget install k6

# Basic run
cd k6_tests
k6 run normal_load.js

# With Grafana visualization (requires InfluxDB + Grafana running)
docker run -d --name influxdb -p 8086:8086 influxdb:1.8
docker run -d --name grafana  -p 3001:3000 grafana/grafana

k6 run --out influxdb=http://localhost:8086/k6 normal_load.js
k6 run --out influxdb=http://localhost:8086/k6 peak_load.js
k6 run --out influxdb=http://localhost:8086/k6 spike_load.js
k6 run --out influxdb=http://localhost:8086/k6 endurance_load.js
```

### Grafana setup
1. Open `http://localhost:3001` (admin/admin)
2. Connections → Data Sources → Add → InfluxDB
   - URL: `http://host.docker.internal:8086`
   - Name: `k6`
   - Database: `k6`
3. Dashboards → Import → Upload `k6_tests/grafana_dashboard.json`
4. Select `k6` as data source → Import

### Performance thresholds

| Scenario | VUs | Duration | p95 Threshold | Error Threshold |
|---|---|---|---|---|
| Normal Load | 10 | 2 min | < 500ms | < 1% |
| Peak Load | 50 | 3 min | < 1,000ms | < 5% |
| Spike Load | 100 | 1m 40s | < 2,000ms | < 10% |
| Endurance | 5 | 10 min | < 800ms | < 2% |

---

## 5. GitHub Actions CI/CD

Pipeline runs automatically on every push to `main` or `develop` and on all PRs.

### Pipeline overview

```
Push/PR
  │
  ├── Unit Tests (42s)
  │     └── go test -run unit + coverage ≥ 50%
  │
  ├── Integration Tests (18s) ─┐
  ├── Edge Case Tests (15s)    ├── parallel
  ├── Concurrency Tests (69s) ─┘
  │
  ├── Go Build (15s)
  │
  ├── Newman API Tests (53s)
  │     └── 30 assertions, 18 requests
  │
  ├── React Build (22s) ── parallel with all above
  │
  └── Failure Notification
        └── Telegram alert on any failure
```

### Setup

```bash
git init
git add .
git commit -m "initial commit"
git remote add origin https://github.com/YOUR_USERNAME/blog-app.git
git push -u origin main
```

### Required GitHub Secrets

| Secret | Description |
|---|---|
| `TELEGRAM_BOT_TOKEN` | Bot token from @BotFather |
| `TELEGRAM_CHAT_ID` | Your chat ID from @userinfobot |

Add at: **Settings → Secrets and variables → Actions → New repository secret**

---

## 6. Test Cases Summary

### Unit Tests (Go) — 65 tests

| File | Test Cases | Count |
|---|---|---|
| auth_test.go | Register success/duplicate/missing, Login success/wrong password/unknown email | 6 |
| posts_test.go | Create/List/Delete + ownership authorization | 5 |
| comments_test.go | Create success/empty body/post not found, Delete not owner | 4 |
| utils_test.go | JWT generate/parse/invalid/wrong secret, Slugify basic/special/unique | 6 |
| getme_test.go | GetMe success, UserNotFound | 2 |
| posts_extra_test.go | GetPost, UpdatePost (success/not owner/not found), ListComments | 8 |
| profile_test.go | GetProfile, UpdateProfile (success/conflict/invalid), GetUserPosts | 7 |
| likes_test.go | LikePost, UnlikePost, Duplicate, GetLikes (zero/after/not found) | 7 |
| bookmarks_test.go | BookmarkPost, Remove, Duplicate, GetBookmarks (empty/with data) | 6 |
| follows_test.go | Follow, Unfollow, Duplicate, Self, TargetNotFound, GetFollowers/Following/Feed | 11 |
| tags_test.go | GetTags (empty/with data), GetPostsByTag (found/no results/missing name) | 5 |
| integration_test.go | Auth→Posts, Posts→Comments, Posts→Likes, Posts→Bookmarks chains | 4 |
| edge_cases_test.go | Empty inputs (×6), large payloads (×3), SQL injection, XSS, Unicode | 13 |
| concurrency_test.go | Simultaneous registrations, double like/bookmark/follow, parallel reads, invalid behavior | 16 |
| **TOTAL** | | **65** |

### API Tests (Newman) — 30 assertions, 18 requests

| Folder | Requests | Key Assertions |
|---|---|---|
| Auth | Register, Register duplicate, Login, Login wrong password, Get Me, Get Me no token | Status codes, token presence |
| Posts | Create, Create no auth, List, List by tag, Get not found, Update | Slug generation, auth enforcement |
| Comments | Add comment, List comments, Delete comment | Status codes |
| Profile | Get profile, Get profile 404, Update profile | Username, bio fields |

### E2E Tests (Playwright) — 13 tests

| File | Tests |
|---|---|
| auth.spec.js | Registration flow, Login flow, Wrong password, Protected route redirect, Empty email |
| posts.spec.js | Create post, Create with tags, Publish disabled (no title), Publish disabled (no body) |
| social.spec.js | Like post, Bookmark post, Unlike post, Unauthenticated like blocked |

### Performance Tests (k6) — 4 scenarios

| Script | Scenario | Result |
|---|---|---|
| normal_load.js | 10 VUs, 2 min | p95=77ms, 0% errors ✅ |
| peak_load.js | 50 VUs, 3 min | p95=554ms, 0.45% errors ✅ |
| spike_load.js | 100 VUs, 1m 40s | p95=1,960ms, 0% errors ✅ |
| endurance_load.js | 5 VUs, 10 min | p95=62ms, 0.23% errors ✅ |

---

## 7. Coverage

```
handlers/auth.go:       Register 70%, Login 75%, GetMe 100%
handlers/bookmarks.go:  BookmarkPost 100%, GetBookmarks 100%
handlers/comments.go:   ListComments 100%, CreateComment 88.2%, DeleteComment 62.5%
handlers/follows.go:    FollowUser 90.9%, GetFollowers 100%, GetFeed 100%
handlers/likes.go:      LikePost 100%, GetLikes 100%
handlers/posts.go:      ListPosts 83.3%, GetPost 100%, CreatePost 71.4%, UpdatePost 90.9%, DeletePost 58.3%
handlers/profile.go:    GetProfile 100%, UpdateProfile 88.2%, GetUserPosts 100%
handlers/tags.go:       GetTags 100%, GetPostsByTag 90.9%
utils/jwt.go:           InitJWT 100%, GenerateToken 100%, ParseToken 80%
utils/slug.go:          Slugify 100%

TOTAL: 88.1%
```

### Known coverage gaps
- `DeletePost` — 58.3% (some auth failure branches not covered)
- `DeleteComment` — 62.5% (auth failure branch under specific DB state)
- `extractSlugFromSubpath` — 0% (dead code, never called)
