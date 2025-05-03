package models

import "gorm.io/gorm"

type FaceitProfile struct {
	gorm.Model

	PlayerID uint `gorm:"uniqueIndex"`
	Nickname string
	ELO      int
	Region   string
}
