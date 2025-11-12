package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string    `gorm:"size:64;uniqueIndex"`
	Email       *string   `gorm:"size:128;uniqueIndex"`
	Password    string    `json:"-"`
	DisplayName string    `gorm:"size:128"`
	Bio         string    `gorm:"type:text"`
	AvatarURL   string    `gorm:"size:255"`
	Posts       []Post    `gorm:"foreignKey:AuthorID" json:"-"`
	Comments    []Comment `gorm:"foreignKey:UserID" json:"-"`
}
