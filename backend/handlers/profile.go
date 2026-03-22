package handlers

import (
	"blog-app/middleware"
	"blog-app/models"
	"blog-app/utils"
	"net/http"
	"strconv"
	"strings"
)

type updateProfileRequest struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	// URL: /api/profiles/{username}
	parts := strings.Split(r.URL.Path, "/")
	username := parts[len(parts)-1]

	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "user not found")
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var req updateProfileRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "user not found")
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	user.Bio = req.Bio
	user.Avatar = req.Avatar

	if err := models.DB.Save(&user).Error; err != nil {
		utils.Error(w, http.StatusConflict, "username already taken")
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

func GetUserPosts(w http.ResponseWriter, r *http.Request) {
	// URL: /api/profiles/{id}/posts
	parts := strings.Split(r.URL.Path, "/")
	// parts: ["", "api", "profiles", "{id}", "posts"]
	idStr := parts[len(parts)-2]
	authorID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var posts []models.Post
	models.DB.Where("author_id = ?", authorID).Preload("Author").Order("created_at desc").Find(&posts)

	utils.JSON(w, http.StatusOK, posts)
}
