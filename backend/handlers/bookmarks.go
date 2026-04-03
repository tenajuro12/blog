package handlers

import (
	"blog-app/middleware"
	"blog-app/models"
	"blog-app/utils"
	"net/http"
	"strings"
)

// POST   /api/posts/:slug/bookmark — bookmark a post
// DELETE /api/posts/:slug/bookmark — remove bookmark
func BookmarkPost(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// ["api","posts","{slug}","bookmark"]
	slug := ""
	if len(parts) >= 3 {
		slug = parts[2]
	}

	var post models.Post
	if err := models.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "post not found")
		return
	}

	if r.Method == http.MethodDelete {
		models.DB.Where("user_id = ? AND post_id = ?", userID, post.ID).Delete(&models.Bookmark{})
		utils.JSON(w, http.StatusOK, map[string]string{"message": "bookmark removed"})
		return
	}

	bm := models.Bookmark{UserID: userID, PostID: post.ID}
	if err := models.DB.Create(&bm).Error; err != nil {
		utils.JSON(w, http.StatusOK, map[string]string{"message": "already bookmarked"})
		return
	}
	utils.JSON(w, http.StatusCreated, map[string]string{"message": "bookmarked"})
}

// GET /api/bookmarks — get current user's bookmarked posts
func GetBookmarks(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var bookmarks []models.Bookmark
	models.DB.Where("user_id = ?", userID).
		Preload("Post.Author").
		Order("created_at desc").
		Find(&bookmarks)

	var posts []models.Post
	for _, bm := range bookmarks {
		if bm.Post != nil {
			p := *bm.Post
			enrichPost(&p, userID)
			posts = append(posts, p)
		}
	}
	if posts == nil {
		posts = []models.Post{}
	}
	utils.JSON(w, http.StatusOK, posts)
}
