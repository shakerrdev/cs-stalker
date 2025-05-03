package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Match struct {
	gorm.Model

	UUID       datatypes.UUID `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex"`
	Map        string         `gorm:"size:64"`
	MatchType  string         `gorm:"size:32"` // "competitive", "wingman", etc.
	Team1Score int
	Team2Score int
	IsCS2      bool
	Duration   int // in seconds
	FinishedAt time.Time
	// TODO?: Hasbanned player or/and replay / share code

	// Relationships
	Teams       []MatchTeam
	Players     []MatchPlayer
	LeetifyData *LeetifyMatch `gorm:"foreignKey:MatchID"`
}
