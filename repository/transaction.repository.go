package repository

import (
	"go-fundraising/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetByCampaignID(campaignID int) ([]entity.Transaction, error)
}

type transactionRepository struct {
	connection *gorm.DB
}

func NewTransactionRepository(connection *gorm.DB) TransactionRepository {
	return &transactionRepository{
		connection: connection,
	}
}

func (r *transactionRepository) GetByCampaignID(campaignID int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.connection.
		Preload("User").
		Where("campaign_id = ?", campaignID).
		Order("id desc").
		Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
