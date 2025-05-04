package models

import "gorm.io/gorm"

type LeetifyProfile struct {
	gorm.Model

	PlayerID  uint   `gorm:"uniqueIndex;not null"`
	LeetifyID string `gorm:"size:64"`
}
