package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `gorm:"size:64;uniqueIndex"`
	Slug        string `gorm:"size:64;uniqueIndex"`
	Description string `gorm:"size:255"`
	Posts       []Post `json:"-"`
}
