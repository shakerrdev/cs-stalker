package models

import "gorm.io/gorm"

type MatchTeam struct {
	gorm.Model

	MatchID    uint
	TeamNumber int // Typically 2 (T) and 3 (CT)
	Score      int
	Players    []MatchPlayer `gorm:"foreignKey:TeamID"`
}
