package faceit

import "gorm.io/gorm"

type FaceitProfile struct {
	gorm.Model

	FaceitID string `gorm:"uniqueIndex;not null;size:36"`
	PlayerID uint   `gorm:"uniqueIndex;not null"`
	Nickname string `gorm:"size:128"`
	ELO      int
	Region   string `gorm:"size:64"`
}
