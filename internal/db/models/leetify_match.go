package models

import (
	"time"

	"gorm.io/gorm"
)

type LeetifyMatch struct {
	gorm.Model

	MatchID     uint   `gorm:"uniqueIndex"`
	LeetifyID   string `gorm:"size:64"`
	CreatedAt   time.Time
	Status      string `gorm:"size:32"`
	DataSource  string `gorm:"size:32"` // "matchmaking", "faceit", etc.
	Recalculate bool

	// Relationship
	Match Match `gorm:"foreignKey:MatchID"`
}
