package leetify

import (
	"time"

	"gorm.io/gorm"
)

type LeetifyProfile struct {
	gorm.Model
	PlayerID       uint   `gorm:"uniqueIndex"`
	LeetifyUserID  string `gorm:"uniqueIndex"`
	Steam64ID      string `gorm:"uniqueIndex;size:17"`
	Name           string
	SteamNickname  string
	SteamAvatarURL string
	IsProPlan      bool

	// Ratings
	Aim         float64
	Positioning float64
	Utility     float64
	Leetify     float64
	CtLeetify   float64
	TLeetify    float64
	Opening     float64
	Clutch      float64

	// Relationships
	Club          LeetifyClub
	Teammates     []LeetifyTeammate
	Highlights    []LeetifyHighlight
	PersonalBests []LeetifyPersonalBest
	Matches       []LeetifyMatch `gorm:"many2many:leetify_match_profiles;"`
}

type LeetifyClub struct {
	gorm.Model
	LeetifyProfileID uint   `gorm:"uniqueIndex"`
	ClubID           string `gorm:"index"`
	Name             string
	Tag              string
}

type LeetifyTeammate struct {
	gorm.Model
	LeetifyProfileID      uint `gorm:"index"`
	Steam64Id             string
	SteamNickname         string
	SteamAvatarURL        string
	MatchesPlayedTogether int
	WinRateTogether       float64

	Rank LeetifyRank
}

type LeetifyRank struct {
	gorm.Model
	LeetifyTeammateID uint `gorm:"uniqueIndex"`
	Type              string
	DataSource        string
	SkillLevel        int
}

type LeetifyHighlight struct {
	gorm.Model
	LeetifyProfileID uint `gorm:"index"`
	Description      string
	GameID           string
	RoundNumber      int
	ThumbnailURL     string
	URL              string
	LeetifyCreatedAt time.Time
}

type LeetifyPersonalBest struct {
	gorm.Model
	LeetifyProfileID uint `gorm:"index"`
	GameID           string
	SkillID          string
	Value            string
	IsCS2            bool `gorm:"default:true"`
}
