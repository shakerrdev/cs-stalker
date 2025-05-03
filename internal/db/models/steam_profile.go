package models

import "gorm.io/gorm"

type SteamProfile struct {
	gorm.Model

	PlayerID        uint `gorm:"uniqueIndex"`
	SteamName       string
	AvatarURL       string
	IsProfilePublic bool
}
