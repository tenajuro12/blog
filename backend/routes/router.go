package routes

import (
	"blog-app/handlers"
	"blog-app/middleware"
	"net/http"
	"strings"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("/api/auth/register", method("POST", handlers.Register))
	mux.HandleFunc("/api/auth/login", method("POST", handlers.Login))
	mux.HandleFunc("/api/auth/me", method("GET", middleware.Auth(handlers.GetMe)))

	// Posts
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

		// /api/posts/{slug}/comments or /api/posts/{slug}/comments/{id}
		if strings.Contains(path, "/comments") {
			parts := strings.Split(strings.Trim(path, "/"), "/")
			// ["api", "posts", "{slug}", "comments"] or [..., "{id}"]
			if len(parts) == 4 {
				// /api/posts/{slug}/comments
				switch r.Method {
				case http.MethodGet:
					handlers.ListComments(w, r)
				case http.MethodPost:
					middleware.Auth(handlers.CreateComment)(w, r)
				default:
					http.NotFound(w, r)
				}
			} else if len(parts) == 5 {
				// /api/posts/{slug}/comments/{id}
				if r.Method == http.MethodDelete {
					middleware.Auth(handlers.DeleteComment)(w, r)
				} else {
					http.NotFound(w, r)
				}
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

	// Profiles
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
