package routes

import (
	"blog-app/handlers"
	"blog-app/middleware"
	"net/http"
	"strings"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// ── Auth ──────────────────────────────────────────────────────────────────
	mux.HandleFunc("/api/auth/register", method("POST", handlers.Register))
	mux.HandleFunc("/api/auth/login", method("POST", handlers.Login))
	mux.HandleFunc("/api/auth/me", method("GET", middleware.Auth(handlers.GetMe)))

	// ── Posts ─────────────────────────────────────────────────────────────────
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListPosts(w, r)
		case http.MethodPost:
			middleware.Auth(handlers.CreatePost)(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// /api/posts/{slug}/comments[/{id}]
		if strings.Contains(path, "/comments") {
			parts := strings.Split(strings.Trim(path, "/"), "/")
			if len(parts) == 4 {
				switch r.Method {
				case http.MethodGet:
					handlers.ListComments(w, r)
				case http.MethodPost:
					middleware.Auth(handlers.CreateComment)(w, r)
				default:
					http.NotFound(w, r)
				}
			} else if len(parts) == 5 {
				if r.Method == http.MethodDelete {
					middleware.Auth(handlers.DeleteComment)(w, r)
				} else {
					http.NotFound(w, r)
				}
			}
			return
		}

		// /api/posts/{slug}/like
		if strings.HasSuffix(path, "/like") {
			switch r.Method {
			case http.MethodPost:
				middleware.Auth(handlers.LikePost)(w, r)
			case http.MethodDelete:
				middleware.Auth(handlers.LikePost)(w, r)
			default:
				http.NotFound(w, r)
			}
			return
		}

		// /api/posts/{slug}/likes
		if strings.HasSuffix(path, "/likes") {
			if r.Method == http.MethodGet {
				handlers.GetLikes(w, r)
			}
			return
		}

		// /api/posts/{slug}/bookmark
		if strings.HasSuffix(path, "/bookmark") {
			switch r.Method {
			case http.MethodPost:
				middleware.Auth(handlers.BookmarkPost)(w, r)
			case http.MethodDelete:
				middleware.Auth(handlers.BookmarkPost)(w, r)
			default:
				http.NotFound(w, r)
			}
			return
		}

		// /api/posts/{slug}
		switch r.Method {
		case http.MethodGet:
			handlers.GetPost(w, r)
		case http.MethodPut:
			middleware.Auth(handlers.UpdatePost)(w, r)
		case http.MethodDelete:
			middleware.Auth(handlers.DeletePost)(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	// ── Feed ──────────────────────────────────────────────────────────────────
	mux.HandleFunc("/api/feed", method("GET", middleware.Auth(handlers.GetFeed)))
	mux.HandleFunc("/api/feed/", method("GET", middleware.Auth(handlers.GetFeed)))

	// ── Bookmarks ─────────────────────────────────────────────────────────────
	mux.HandleFunc("/api/bookmarks", method("GET", middleware.Auth(handlers.GetBookmarks)))
	mux.HandleFunc("/api/bookmarks/", method("GET", middleware.Auth(handlers.GetBookmarks)))

	// ── Tags ──────────────────────────────────────────────────────────────────
	mux.HandleFunc("/api/tags/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		// /api/tags/ or /api/tags → list all tags
		path := strings.TrimSuffix(r.URL.Path, "/")
		if path == "/api/tags" {
			handlers.GetTags(w, r)
			return
		}
		// /api/tags/{name} → posts by tag
		handlers.GetPostsByTag(w, r)
	})

	// ── Users (follow/unfollow) ───────────────────────────────────────────────
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/follow") {
			switch r.Method {
			case http.MethodPost:
				middleware.Auth(handlers.FollowUser)(w, r)
			case http.MethodDelete:
				middleware.Auth(handlers.FollowUser)(w, r)
			default:
				http.NotFound(w, r)
			}
			return
		}
		if strings.HasSuffix(path, "/followers") {
			handlers.GetFollowers(w, r)
			return
		}
		if strings.HasSuffix(path, "/following") {
			handlers.GetFollowing(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// ── Profile ───────────────────────────────────────────────────────────────
	mux.HandleFunc("/api/profile", method("PUT", middleware.Auth(handlers.UpdateProfile)))
	mux.HandleFunc("/api/profiles/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/posts") {
			handlers.GetUserPosts(w, r)
		} else {
			handlers.GetProfile(w, r)
		}
	})

	return cors(mux)
}

func method(m string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != m {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
