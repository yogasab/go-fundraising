package dto

import "go-fundraising/entity"

type TransactionGetRequestID struct {
	ID   int `uri:"id" binding:"required"`
	User entity.User
}

type TransactionCreateRequest struct {
	Amount     int `json:"amount" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
	User       entity.User
}
