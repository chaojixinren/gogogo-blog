package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name  string `gorm:"size:64;uniqueIndex"`
	Slug  string `gorm:"size:64;uniqueIndex"`
	Posts []Post `gorm:"many2many:post_tags" json:"-"`
}
