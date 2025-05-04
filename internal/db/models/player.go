package models

import (
	"github.com/finnbechinka/cs-stalker/internal/db/models/faceit"
	"github.com/finnbechinka/cs-stalker/internal/db/models/leetify"
	"github.com/finnbechinka/cs-stalker/internal/db/models/steam"
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model

	SteamID64 string   `gorm:"uniqueIndex;not null;size:17"`
	Name      string   `gorm:"size:128"`
	AvatarURL string   `gorm:"size:256"`
	Region    string   `gorm:"size:64"`
	Bans      []string `gorm:"type:jsonb"`

	LeetifyProfile leetify.LeetifyProfile
	FaceitProfile  faceit.FaceitProfile
	SteamProfile   steam.SteamProfile
}
