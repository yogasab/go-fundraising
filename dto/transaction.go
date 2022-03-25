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

type TransactionNotificationRequest struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
}
