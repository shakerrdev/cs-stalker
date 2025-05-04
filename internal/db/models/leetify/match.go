package leetify

import "gorm.io/gorm"

type LeetifyMatch struct {
	gorm.Model
	MatchID        uint   `gorm:"index;uniqueIndex;not null"`
	LeetifyID      string `gorm:"uniqueIndex"`
	DataSource     string
	IsCS2          bool
	MapName        string
	Team1Score     int
	Team2Score     int
	FinishedAt     string
	Status         string
	SteamShareCode string

	// Relationships
	Details     LeetifyMatchDetails
	PlayerStats []LeetifyPlayerStats
	Agents      []LeetifyAgent
	Parties     []LeetifyParty
	Profiles    []LeetifyProfile `gorm:"many2many:leetify_match_profiles;"`
}

type LeetifyMatchDetails struct {
	gorm.Model
	LeetifyMatchID uint
	Tickrate       int
	Ticks          int
	ServerName     string
}

type LeetifyPlayerStats struct {
	gorm.Model
	LeetifyMatchID   uint
	LeetifyProfileID uint
	Steam64Id        string
	Name             string
	Kills            int
	Deaths           int
	Assists          int
	Score            int
	LeetifyRating    float64
	CtLeetifyRating  float64
	TLeetifyRating   float64

	// Aim stats
	AimRating    float64
	Preaim       float64
	ReactionTime float64

	// Utility stats
	UtilityRating   float64
	FlashbangThrown int
	FlashbangHitFoe int

	// Other stats
	HltvRating float64
	Adr        float64
}

type LeetifyAgent struct {
	gorm.Model
	LeetifyMatchID uint
	Steam64Id      string
	TeamNumber     int
	Skin           string
}

type LeetifyParty struct {
	gorm.Model
	LeetifyMatchID uint
	Steam64Id      string
	Party          int
}
