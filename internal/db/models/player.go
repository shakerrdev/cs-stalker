package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model

	SteamID64    string   `gorm:"uniqueIndex;not null;size:17"`
	Name         string   `gorm:"size:128"`
	AvatarURL    string   `gorm:"size:256"`
	Region       string   `gorm:"size:64"`
	PlatformBans []string `gorm:"type:jsonb"`

	SteamProfile   SteamProfile
	FaceitProfile  FaceitProfile
	LeetifyProfile LeetifyProfile
}
