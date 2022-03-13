package entity

import "time"

type Campaign struct {
	ID               int    `gorm:"primaryKey:auto_increment"`
	UserId           int    `gorm:"not null"`
	Name             string `gorm:"type:varchar(100)"`
	ShortDescription string `gorm:"type:varchar(100)"`
	Description      string `gorm:"type:varchar(100)"`
	Perks            string `gorm:"type:varchar(100)"`
	BackerCount      int    `gorm:"bigint"`
	GoalAmount       int    `gorm:"bigint"`
	CurrentAmount    int    `gorm:"bigint"`
	Slug             string `gorm:"type:varchar(100)"`
	CampaignImages   []CampaignImage
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CampaignImage struct {
	ID         int
	CampaignID int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
