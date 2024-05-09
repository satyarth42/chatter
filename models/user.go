package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UUID         string `gorm:"uuid;uniqueIndex:idx_uuid;size:36"`
	Name         string `gorm:"name;not null;size:128"`
	Email        string `gorm:"email;uniqueIndex:idx_email;not null;size:128"`
	Password     string `gorm:"password;not null;size:128"`
	Verification bool   `gorm:"verification; default:false" `
}

func (u User) TableName() string {
	return "users"
}
