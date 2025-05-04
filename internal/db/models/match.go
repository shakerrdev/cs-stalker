package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/finnbechinka/cs-stalker/internal/db/models/leetify"
)

type Match struct {
	gorm.Model

	UUID       datatypes.UUID `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex"`
	Map        string         `gorm:"size:64;not null"`
	MatchType  string         `gorm:"size:32"` // "competitive", "wingman", etc.
	ScoreTeam1 int
	ScoreTeam2 int
	MatchDate  string
	Duration   int // in seconds
	// TODO?: Hasbanned player or/and replay / share code, iscs2?

	// Relationships
	LeetifyMatch leetify.LeetifyMatch
}
