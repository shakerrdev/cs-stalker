package models

import (
	"time"

	"gorm.io/gorm"
)

type LeetifyMatch struct {
	gorm.Model

	MatchID     uint   `gorm:"uniqueIndex;not null"`
	LeetifyID   string `gorm:"size:64;index"`
	CreatedAt   time.Time
	Status      string `gorm:"size:32;not null"`
	DataSource  string `gorm:"size:32"`
	Recalculate bool

	// Relationship
	Match Match `gorm:"foreignKey:MatchID"`
}
