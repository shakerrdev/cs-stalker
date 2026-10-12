package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/shakerrdev/cs-stalker/internal/db/models/leetify"
)

type Match struct {
	gorm.Model

	UUID            string `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`
	Map             string `gorm:"size:64;not null"`
	MatchType       string `gorm:"size:32"` // "competitive", "wingman", etc.
	ScoreTeam1      int
	ScoreTeam2      int
	MatchDate       time.Time
	DurationSeconds int
	// TODO?: Hasbanned player or/and replay / share code, iscs2?

	// Relationships
	LeetifyMatch leetify.LeetifyMatch
}
