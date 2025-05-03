package models

import "gorm.io/gorm"

type LeetifyProfile struct {
	gorm.Model

	PlayerID  uint   `gorm:"uniqueIndex"`
	LeetifyID string // internal ID if any
}
