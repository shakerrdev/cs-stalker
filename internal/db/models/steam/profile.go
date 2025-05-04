package steam

import "gorm.io/gorm"

type SteamProfile struct {
	gorm.Model

	Steam64ID       string `gorm:"uniqueIndex;not null;size:17"`
	PlayerID        uint   `gorm:"index;uniqueIndex;not null"`
	SteamName       string `gorm:"size:128"`
	AvatarURL       string `gorm:"size:256"`
	IsProfilePublic bool
}
