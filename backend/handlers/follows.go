package handlers

import (
	"blog-app/middleware"
	"blog-app/models"
	"blog-app/utils"
	"net/http"
	"strconv"
	"strings"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, _ := middleware.GetUserID(r)

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// ["api","users","{id}","follow"]
	if len(parts) < 3 {
		utils.Error(w, http.StatusBadRequest, "invalid path")
		return
	}
	targetID, err := strconv.Atoi(parts[2])
	if err != nil || uint(targetID) == followerID {
		utils.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// ensure target exists
	var target models.User
	if err := models.DB.First(&target, targetID).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "user not found")
		return
	}

	if r.Method == http.MethodDelete {
		models.DB.Where("follower_id = ? AND following_id = ?", followerID, targetID).Delete(&models.Follow{})
		utils.JSON(w, http.StatusOK, map[string]string{"message": "unfollowed"})
		return
	}

	follow := models.Follow{FollowerID: followerID, FollowingID: uint(targetID)}
	if err := models.DB.Create(&follow).Error; err != nil {
		utils.JSON(w, http.StatusOK, map[string]string{"message": "already following"})
		return
	}
	utils.JSON(w, http.StatusCreated, map[string]string{"message": "followed"})
}

// GET /api/users/:id/followers
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	userID, _ := strconv.Atoi(parts[2])

	var follows []models.Follow
	models.DB.Where("following_id = ?", userID).Find(&follows)

	var users []models.User
	for _, f := range follows {
		var u models.User
		if models.DB.First(&u, f.FollowerID).Error == nil {
			users = append(users, u)
		}
	}
	if users == nil {
		users = []models.User{}
	}
	utils.JSON(w, http.StatusOK, users)
}

// GET /api/users/:id/following
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	userID, _ := strconv.Atoi(parts[2])

	var follows []models.Follow
	models.DB.Where("follower_id = ?", userID).Find(&follows)

	var users []models.User
	for _, f := range follows {
		var u models.User
		if models.DB.First(&u, f.FollowingID).Error == nil {
			users = append(users, u)
		}
	}
	if users == nil {
		users = []models.User{}
	}
	utils.JSON(w, http.StatusOK, users)
}

// GET /api/feed — posts from followed users
func GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var follows []models.Follow
	models.DB.Where("follower_id = ?", userID).Find(&follows)

	var followingIDs []uint
	for _, f := range follows {
		followingIDs = append(followingIDs, f.FollowingID)
	}

	var posts []models.Post
	if len(followingIDs) > 0 {
		models.DB.Where("author_id IN ?", followingIDs).
			Preload("Author").
			Order("created_at desc").
			Find(&posts)
	}

	for i := range posts {
		enrichPost(&posts[i], userID)
	}

	if posts == nil {
		posts = []models.Post{}
	}
	utils.JSON(w, http.StatusOK, posts)
}
