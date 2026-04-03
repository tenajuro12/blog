package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Bio       string    `json:"bio"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Posts     []Post    `json:"posts,omitempty" gorm:"foreignKey:AuthorID"`
}

type Post struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title" gorm:"not null"`
	Slug       string    `json:"slug" gorm:"uniqueIndex;not null"`
	Body       string    `json:"body" gorm:"not null"`
	Tags       string    `json:"tags"`
	AuthorID   uint      `json:"author_id"`
	Author     *User     `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Comments   []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	LikeCount  int       `json:"like_count" gorm:"-"`
	Liked      bool      `json:"liked" gorm:"-"`
	Bookmarked bool      `json:"bookmarked" gorm:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Body      string    `json:"body" gorm:"not null"`
	PostID    uint      `json:"post_id"`
	AuthorID  uint      `json:"author_id"`
	Author    *User     `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Like — many-to-many between users and posts
type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"uniqueIndex:idx_like_user_post;not null"`
	PostID    uint      `json:"post_id" gorm:"uniqueIndex:idx_like_user_post;not null"`
	CreatedAt time.Time `json:"created_at"`
}

type Follow struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FollowerID  uint      `json:"follower_id" gorm:"uniqueIndex:idx_follow;not null"`
	FollowingID uint      `json:"following_id" gorm:"uniqueIndex:idx_follow;not null"`
	CreatedAt   time.Time `json:"created_at"`
}

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at"`
}

type Bookmark struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"uniqueIndex:idx_bookmark;not null"`
	PostID    uint      `json:"post_id" gorm:"uniqueIndex:idx_bookmark;not null"`
	Post      *Post     `json:"post,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `json:"created_at"`
}
