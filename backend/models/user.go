package models

import "gorm.io/gorm"

type UserProfile struct {
	gorm.Model
	FirebaseUID string `gorm:"uniqueIndex;not null"`
	Username    string
	Bio         string
}
