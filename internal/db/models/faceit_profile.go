package models

import "gorm.io/gorm"

type FaceitProfile struct {
	gorm.Model

	PlayerID uint   `gorm:"uniqueIndex;not null"`
	Nickname string `gorm:"size:128"`
	ELO      int
	Region   string `gorm:"size:64"`
}
