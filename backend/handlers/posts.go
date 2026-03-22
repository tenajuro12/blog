package handlers

import (
	"blog-app/middleware"
	"blog-app/models"
	"blog-app/utils"
	"net/http"
	"strings"
)

type postRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Tags  string `json:"tags"`
}

func ListPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	query := models.DB.Preload("Author").Order("created_at desc")

	if tag := r.URL.Query().Get("tag"); tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	query.Find(&posts)
	utils.JSON(w, http.StatusOK, posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	slug := extractSlug(r.URL.Path)

	var post models.Post
	if err := models.DB.Where("slug = ?", slug).
		Preload("Author").
		Preload("Comments.Author").
		First(&post).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "post not found")
		return
	}

	utils.JSON(w, http.StatusOK, post)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var req postRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" || req.Body == "" {
		utils.Error(w, http.StatusBadRequest, "title and body are required")
		return
	}

	post := models.Post{
		Title:    req.Title,
		Body:     req.Body,
		Tags:     req.Tags,
		Slug:     utils.Slugify(req.Title),
		AuthorID: userID,
	}

	if err := models.DB.Create(&post).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to create post")
		return
	}

	models.DB.Preload("Author").First(&post, post.ID)
	utils.JSON(w, http.StatusCreated, post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	slug := extractSlug(r.URL.Path)

	var post models.Post
	if err := models.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "post not found")
		return
	}

	if post.AuthorID != userID {
		utils.Error(w, http.StatusForbidden, "you can only edit your own posts")
		return
	}

	var req postRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title != "" {
		post.Title = req.Title
		post.Slug = utils.Slugify(req.Title)
	}
	if req.Body != "" {
		post.Body = req.Body
	}
	post.Tags = req.Tags

	models.DB.Save(&post)
	models.DB.Preload("Author").First(&post, post.ID)
	utils.JSON(w, http.StatusOK, post)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	slug := extractSlug(r.URL.Path)

	var post models.Post
	if err := models.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "post not found")
		return
	}

	if post.AuthorID != userID {
		utils.Error(w, http.StatusForbidden, "you can only delete your own posts")
		return
	}

	models.DB.Where("post_id = ?", post.ID).Delete(&models.Comment{})
	models.DB.Delete(&post)

	utils.JSON(w, http.StatusOK, map[string]string{"message": "post deleted"})
}

// extractSlug gets slug from paths like /api/posts/{slug} or /api/posts/{slug}/comments
func extractSlug(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// ["api", "posts", "{slug}"] or ["api", "posts", "{slug}", "comments"]
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}
