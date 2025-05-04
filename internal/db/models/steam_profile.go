package models

import "gorm.io/gorm"

type SteamProfile struct {
	gorm.Model

	PlayerID        uint   `gorm:"index;uniqueIndex;not null"`
	SteamName       string `gorm:"size:128"`
	AvatarURL       string `gorm:"size:256"`
	IsProfilePublic bool
}
