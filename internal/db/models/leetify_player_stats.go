package models

import "gorm.io/gorm"

type LeetifyPlayerMatchStats struct {
	gorm.Model

	MatchPlayerID uint   // References MatchPlayer
	LeetifyUserID string `gorm:"size:64"`

	// Ratings
	LeetifyRating             float64
	CTLeetifyRating           float64
	TLeetifyRating            float64
	PersonalPerformanceRating float64

	// Detailed stats
	AimRating          float64
	UtilityRating      float64
	ClutchRating       float64
	OpeningRating      float64
	FlashbangRating    float64
	PreaimRating       float64
	ReactionTime       float64
	CrosshairPlacement float64

	// Relationship
	MatchPlayer MatchPlayer `gorm:"foreignKey:MatchPlayerID"`
}
