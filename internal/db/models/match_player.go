package models

import "gorm.io/gorm"

type MatchPlayer struct {
	gorm.Model

	MatchID  uint
	PlayerID uint // References Player.SteamID64
	TeamID   uint // References MatchTeam

	// Basic stats
	Kills   int
	Deaths  int
	Assists int
	Score   int
	MVPs    int
	ADR     float64

	// Relationships
	Player       Player `gorm:"foreignKey:SteamID64;references:PlayerID"`
	Team         MatchTeam
	LeetifyStats *LeetifyPlayerMatchStats `gorm:"foreignKey:MatchPlayerID"`
}
