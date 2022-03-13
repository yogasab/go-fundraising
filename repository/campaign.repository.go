package repository

import (
	"go-fundraising/entity"
	"gorm.io/gorm"
)

type CampaignRepository interface {
	FindCampaignByUserID(userID int) ([]entity.Campaign, error)
	FindAll() ([]entity.Campaign, error)
}

type campaignRepository struct {
	connection *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) CampaignRepository {
	return &campaignRepository{
		connection: db,
	}
}

func (r *campaignRepository) FindCampaignByUserID(userID int) ([]entity.Campaign, error) {
	var campaigns []entity.Campaign
	err := r.connection.
		Where("user_id = ?", userID).
		Preload("CampaignImages", "is_primary = 1").
		Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *campaignRepository) FindAll() ([]entity.Campaign, error) {
	var campaigns []entity.Campaign
	err := r.connection.
		Preload("CampaignImages", "is_primary = 1").
		Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}
