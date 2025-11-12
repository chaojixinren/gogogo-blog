package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	PostStatusDraft     = "draft"
	PostStatusPublished = "published"
	PostStatusArchived  = "archived"
)

type Post struct {
	gorm.Model
	Title       string     `gorm:"size:200;not null"`
	Summary     string     `gorm:"size:512"`
	Content     string     `gorm:"type:longtext"`
	Slug        string     `gorm:"size:200;uniqueIndex"`
	Status      string     `gorm:"size:32;default:draft"`
	CoverImage  string     `gorm:"size:255"`
	PublishedAt *time.Time `json:"publishedAt"`
	AuthorID    uint       `json:"authorId"`
	Author      User       `json:"author"`
	CategoryID  *uint      `json:"categoryId"`
	Category    *Category  `json:"category"`
	Tags        []Tag      `gorm:"many2many:post_tags" json:"tags"`
	Comments    []Comment  `json:"comments"`
}

func (p *Post) IsPublished() bool {
	return p.Status == PostStatusPublished && p.PublishedAt != nil
}
