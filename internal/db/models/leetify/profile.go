package leetify

import "gorm.io/gorm"

type LeetifyProfile struct {
	gorm.Model
	PlayerID       uint   `gorm:"index;uniqueIndex"`
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
	LeetifyProfileID uint
	ID               string
	Name             string
	Tag              string
}

type LeetifyTeammate struct {
	gorm.Model
	LeetifyProfileID      uint
	Steam64Id             string
	SteamNickname         string
	SteamAvatarURL        string
	MatchesPlayedTogether int
	WinRateTogether       float64

	Rank LeetifyRank
}

type LeetifyRank struct {
	gorm.Model
	LeetifyTeammateID uint
	Type              string
	DataSource        string
	SkillLevel        int
}

type LeetifyHighlight struct {
	gorm.Model
	LeetifyProfileID uint
	Description      string
	GameID           string
	RoundNumber      int
	ThumbnailURL     string
	URL              string
	CreatedAt        string
}

type LeetifyPersonalBest struct {
	gorm.Model
	LeetifyProfileID uint
	GameID           string
	SkillID          string
	Value            string
	IsCS2            bool `gorm:"default:true"`
}
