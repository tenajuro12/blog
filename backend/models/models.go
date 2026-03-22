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
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Slug      string    `json:"slug" gorm:"uniqueIndex;not null"`
	Body      string    `json:"body" gorm:"not null"`
	Tags      string    `json:"tags"`
	AuthorID  uint      `json:"author_id"`
	Author    *User     `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Comments  []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
