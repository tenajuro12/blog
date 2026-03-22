package handlers

import (
	"blog-app/middleware"
	"blog-app/models"
	"blog-app/utils"
	"net/http"
	"strconv"
	"strings"
)

type commentRequest struct {
	Body string `json:"body"`
}

func ListComments(w http.ResponseWriter, r *http.Request) {
	slug := extractSlug(r.URL.Path)

	var post models.Post
	if err := models.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "post not found")
		return
	}

	var comments []models.Comment
	models.DB.Where("post_id = ?", post.ID).Preload("Author").Order("created_at asc").Find(&comments)
	utils.JSON(w, http.StatusOK, comments)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	slug := extractSlug(r.URL.Path)

	var post models.Post
	if err := models.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "post not found")
		return
	}

	var req commentRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Body == "" {
		utils.Error(w, http.StatusBadRequest, "comment body is required")
		return
	}

	comment := models.Comment{
		Body:     req.Body,
		PostID:   post.ID,
		AuthorID: userID,
	}

	models.DB.Create(&comment)
	models.DB.Preload("Author").First(&comment, comment.ID)
	utils.JSON(w, http.StatusCreated, comment)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	// URL: /api/posts/{slug}/comments/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := parts[len(parts)-1]
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid comment id")
		return
	}

	var comment models.Comment
	if err := models.DB.First(&comment, commentID).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "comment not found")
		return
	}

	if comment.AuthorID != userID {
		utils.Error(w, http.StatusForbidden, "you can only delete your own comments")
		return
	}

	models.DB.Delete(&comment)
	utils.JSON(w, http.StatusOK, map[string]string{"message": "comment deleted"})
}
