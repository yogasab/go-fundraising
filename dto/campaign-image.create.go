package dto

import "go-fundraising/entity"

type CreateCampaignImageRequest struct {
	CampaignID int  `form:"campaign_id" binding:"required"`
	IsPrimary  bool `form:"is_primary"`
	User       entity.User
}
