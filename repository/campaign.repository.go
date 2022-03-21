package repository

import (
	"go-fundraising/entity"

	"gorm.io/gorm"
)

type CampaignRepository interface {
	FindCampaignByUserID(userID int) ([]entity.Campaign, error)
	FindAll() ([]entity.Campaign, error)
	FindCampaignByID(id int) (entity.Campaign, error)
	FindCampaignBySlug(slug string) (entity.Campaign, error)
	Save(campaign entity.Campaign) (entity.Campaign, error)
	Update(campaign entity.Campaign) (entity.Campaign, error)
	CreateImage(campaignImage entity.CampaignImage) (entity.CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID int) (bool, error)
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

func (r *campaignRepository) FindCampaignByID(id int) (entity.Campaign, error) {
	var campaign entity.Campaign
	err := r.connection.
		Preload("User").
		Preload("CampaignImages").
		Where("id = ?", id).
		Find(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *campaignRepository) FindCampaignBySlug(slug string) (entity.Campaign, error) {
	var campaign entity.Campaign
	err := r.connection.
		Preload("User").
		Preload("CampaignImages").
		Where("slug = ?", slug).
		Find(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *campaignRepository) Save(campaign entity.Campaign) (entity.Campaign, error) {
	err := r.connection.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *campaignRepository) Update(campaign entity.Campaign) (entity.Campaign, error) {
	err := r.connection.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *campaignRepository) CreateImage(campaignImage entity.CampaignImage) (entity.CampaignImage, error) {
	err := r.connection.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}
	return campaignImage, nil
}

func (r *campaignRepository) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	err := r.connection.Model(&entity.CampaignImage{}).
		Where("campaign_id = ?", campaignID).
		Update("is_primary", 0).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
