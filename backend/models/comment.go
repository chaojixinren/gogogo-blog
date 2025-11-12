package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID     uint   `json:"postId"`
	Post       Post   `json:"-"`
	UserID     *uint  `json:"userId"`
	User       *User  `json:"user,omitempty"`
	AuthorName string `gorm:"size:128" json:"authorName"`
	Body       string `gorm:"type:text" json:"body"`
	Approved   bool   `json:"approved"`
}
