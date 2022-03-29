package entity

import (
	"github.com/leekchan/accounting"
	"time"
)

type Campaign struct {
	ID               int             `gorm:"primaryKey:auto_increment" json:"id"`
	UserId           int             `gorm:"not null" json:"user_id"`
	Name             string          `gorm:"type:varchar(100)" json:"name"`
	ShortDescription string          `gorm:"type:varchar(100)" json:"short_description"`
	Description      string          `gorm:"type:varchar(100)" json:"description"`
	Perks            string          `gorm:"type:varchar(100)" json:"perks"`
	BackerCount      int             `gorm:"bigint" json:"backer_count"`
	GoalAmount       int             `gorm:"bigint" json:"goal_amount"`
	CurrentAmount    int             `gorm:"bigint" json:"current_amount"`
	Slug             string          `gorm:"type:varchar(100)" json:"slug"`
	CampaignImages   []CampaignImage `json:"campaign_images"`
	User             User            `json:"user"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func (campaign Campaign) FormatGoalAmountIDR() string {
	acc := accounting.Accounting{Symbol: "Rp", Precision: 2, Thousand: ".", Decimal: ","}
	return acc.FormatMoney(campaign.GoalAmount)
}

type CampaignImage struct {
	ID         int       `json:"id"`
	CampaignID int       `json:"campaign_id"`
	FileName   string    `json:"file_name"`
	IsPrimary  int       `json:"is_primary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
