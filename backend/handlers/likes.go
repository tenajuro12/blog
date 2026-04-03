package handlers

import (
	"blog-app/middleware"
	"blog-app/models"
	"blog-app/utils"
	"net/http"
	"strings"
)

// POST /api/posts/:slug/like   — like a post
// DELETE /api/posts/:slug/like — unlike a post
func LikePost(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	slug := extractSlug(r.URL.Path) // reuse from posts.go

	var post models.Post
	if err := models.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "post not found")
		return
	}

	if r.Method == http.MethodDelete {
		models.DB.Where("user_id = ? AND post_id = ?", userID, post.ID).Delete(&models.Like{})
		utils.JSON(w, http.StatusOK, map[string]string{"message": "unliked"})
		return
	}

	like := models.Like{UserID: userID, PostID: post.ID}
	if err := models.DB.Create(&like).Error; err != nil {
		// already liked — just return ok
		utils.JSON(w, http.StatusOK, map[string]string{"message": "already liked"})
		return
	}
	utils.JSON(w, http.StatusCreated, map[string]string{"message": "liked"})
}

// GET /api/posts/:slug/likes — get like count for a post
func GetLikes(w http.ResponseWriter, r *http.Request) {
	slug := extractSlug(r.URL.Path)

	var post models.Post
	if err := models.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "post not found")
		return
	}

	var count int64
	models.DB.Model(&models.Like{}).Where("post_id = ?", post.ID).Count(&count)

	utils.JSON(w, http.StatusOK, map[string]int64{"likes": count})
}

// helper used by posts handlers to enrich posts with like/bookmark info
func enrichPost(post *models.Post, userID uint) {
	var likeCount int64
	models.DB.Model(&models.Like{}).Where("post_id = ?", post.ID).Count(&likeCount)
	post.LikeCount = int(likeCount)

	if userID > 0 {
		var like models.Like
		post.Liked = models.DB.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&like).Error == nil
		var bm models.Bookmark
		post.Bookmarked = models.DB.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&bm).Error == nil
	}
}

// extractSlugFromLikePath: /api/posts/{slug}/like -> slug
func extractSlugFromSubpath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}
