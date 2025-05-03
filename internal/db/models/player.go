package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model

	SteamID64    string `gorm:"uniqueIndex;not null"`
	Name         string
	AvatarURL    string
	Region       string
	PlatformBans []string

	SteamProfile   SteamProfile
	FaceitProfile  FaceitProfile
	LeetifyProfile LeetifyProfile
}
