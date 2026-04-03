package handlers

import (
	"blog-app/models"
	"blog-app/utils"
	"net/http"
	"strings"
)

// GET /api/tags — list all tags with post counts
func GetTags(w http.ResponseWriter, r *http.Request) {
	// Collect all unique tags from posts
	var posts []models.Post
	models.DB.Select("tags").Where("tags != ''").Find(&posts)

	tagCount := map[string]int{}
	for _, p := range posts {
		for _, t := range strings.Split(p.Tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tagCount[t]++
			}
		}
	}

	type TagResult struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	var result []TagResult
	for name, count := range tagCount {
		result = append(result, TagResult{Name: name, Count: count})
	}
	if result == nil {
		result = []TagResult{}
	}
	utils.JSON(w, http.StatusOK, result)
}

// GET /api/tags/:name — get posts for a specific tag
func GetPostsByTag(w http.ResponseWriter, r *http.Request) {
	// path: /api/tags/{name}
	path := strings.TrimSuffix(r.URL.Path, "/")
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	// parts: ["api", "tags", "{name}"]
	if len(parts) < 3 || parts[2] == "" {
		utils.Error(w, http.StatusBadRequest, "tag name required")
		return
	}
	tagName := parts[2]

	var posts []models.Post
	models.DB.Where("tags LIKE ?", "%"+tagName+"%").
		Preload("Author").
		Order("created_at desc").
		Find(&posts)

	if posts == nil {
		posts = []models.Post{}
	}
	utils.JSON(w, http.StatusOK, posts)
}
