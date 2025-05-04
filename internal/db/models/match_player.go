package models

import "gorm.io/gorm"

type MatchPlayer struct {
	gorm.Model

	MatchID  uint `gorm:"index;not null"`
	PlayerID uint `gorm:"uniqueIndex;not null"`
	TeamID   uint `gorm:"not null"`

	// Basic stats
	Kills   int     `gorm:"default:0"`
	Deaths  int     `gorm:"default:0"`
	Assists int     `gorm:"default:0"`
	Score   int     `gorm:"default:0"`
	MVPs    int     `gorm:"default:0"`
	ADR     float64 `gorm:"default:0"`

	// Relationships
	Player       Player `gorm:"foreignKey:PlayerID"`
	Team         MatchTeam
	LeetifyStats *LeetifyPlayerMatchStats `gorm:"foreignKey:MatchPlayerID"`
}
