.PHONY: test test-unit test-api lint build run

# ── Backend ──────────────────────────────────────────────────────────────────

## Run all Go unit tests
test-unit:
	cd backend && go test ./tests/... -v -coverprofile=coverage.out
	cd backend && go tool cover -func=coverage.out

## Run Go tests with coverage HTML report
test-coverage:
	cd backend && go test ./tests/... -coverprofile=coverage.out
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: backend/coverage.html"

## Build backend binary
build-backend:
	cd backend && go build -o bin/blog-app ./main.go
	@echo "Binary: backend/bin/blog-app"

## Run backend locally (requires Postgres)
run-backend:
	cd backend && go run main.go

# ── API Tests ─────────────────────────────────────────────────────────────────

## Run Newman API tests (requires backend running on :8080)
test-api:
	newman run postman/BlogApp.postman_collection.json \
		--environment postman/BlogApp.postman_environment.json \
		--reporters cli

## Run Newman with HTML report
test-api-report:
	newman run postman/BlogApp.postman_collection.json \
		--environment postman/BlogApp.postman_environment.json \
		--reporters cli,htmlextra \
		--reporter-htmlextra-export newman-report.html
	@echo "Report: newman-report.html"

# ── Frontend ──────────────────────────────────────────────────────────────────

## Install frontend deps
install-frontend:
	cd frontend && npm install --legacy-peer-deps

## Build frontend
build-frontend:
	cd frontend && npm run build

## Start frontend dev server
run-frontend:
	cd frontend && npm start

# ── Docker ────────────────────────────────────────────────────────────────────

## Start only PostgreSQL
db-up:
	docker-compose up -d postgres

## Start full stack
up:
	docker-compose up -d

## Stop everything
down:
	docker-compose down

## Run all tests (unit + api)
test: test-unit
	@echo "\n✓ All unit tests passed"
	@echo "Run 'make test-api' to run API tests (requires backend running)"
