# BlogApp — QA Setup Guide

## Project Structure

```
blog-app/
├── backend/              # Go API (net/http + GORM)
│   ├── tests/            # Unit tests
│   │   ├── auth_test.go
│   │   ├── posts_test.go
│   │   ├── comments_test.go
│   │   └── utils_test.go
│   └── ...
├── frontend/             # React app
├── postman/              # API test collection
│   ├── BlogApp.postman_collection.json
│   └── BlogApp.postman_environment.json
├── .github/workflows/
│   └── ci.yml            # GitHub Actions pipeline
└── Makefile
```

---

## 1. Run Unit Tests (Go)

```bash
# Start Postgres first
make db-up

# Run all unit tests with coverage
make test-unit

# Or directly
cd backend
go test ./tests/... -v
```

Tests use **SQLite in-memory** — no Postgres needed locally.

---

## 2. Run API Tests (Newman/Postman)

### Option A: Postman GUI
1. Open Postman
2. Import `postman/BlogApp.postman_collection.json`
3. Import `postman/BlogApp.postman_environment.json`
4. Run the collection

### Option B: Newman CLI
```bash
# Install Newman
npm install -g newman newman-reporter-htmlextra

# Start backend first
make db-up
make run-backend   # in another terminal

# Run tests
make test-api

# With HTML report
make test-api-report
```

---

## 3. GitHub Actions CI/CD

Pipeline runs automatically on every push to `main` or `develop`.

### Jobs:
| Job | Description |
|-----|-------------|
| `backend-tests` | Go unit tests with coverage |
| `backend-build` | Compiles binary, uploads artifact |
| `api-tests` | Starts server + runs Newman collection |
| `frontend-build` | Builds React app |

### Setup for GitHub:
```bash
# Initialize repo
git init
git add .
git commit -m "initial commit"
git remote add origin https://github.com/YOUR_USERNAME/blog-app.git
git push -u origin main
```

Pipeline starts automatically — check **Actions** tab in GitHub.

### Required Secrets (optional, for deploy):
- None required for basic CI

---

## 4. Coverage

After running `make test-unit`:
```bash
# Terminal summary
go tool cover -func=backend/coverage.out

# HTML report in browser
make test-coverage
open backend/coverage.html
```

---

## Test Cases Summary

### Unit Tests (Go)
| File | Tests |
|------|-------|
| `auth_test.go` | Register success/duplicate/missing fields, Login success/wrong password/unknown email |
| `posts_test.go` | Create success/missing title, List empty/with posts, Delete not owner |
| `comments_test.go` | Create success/empty body/post not found, Delete not owner |
| `utils_test.go` | JWT generate/parse/invalid/wrong secret, Slugify basic/special chars/unique |

### API Tests (Postman/Newman)
| Folder | Requests |
|--------|----------|
| Auth | Register, Register duplicate, Login, Login wrong password, Get Me, Get Me no token |
| Posts | Create, Create no auth, List, List by tag, Get by slug, Get 404, Update |
| Comments | Add, List, Delete |
| Profile | Get, Get 404, Update |
