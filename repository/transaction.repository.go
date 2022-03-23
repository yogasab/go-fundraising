package repository

import (
	"go-fundraising/entity"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetByCampaignID(campaignID int) ([]entity.Transaction, error)
	GetByUserID(userID int) ([]entity.Transaction, error)
	Save(transaction entity.Transaction) (entity.Transaction, error)
	Update(transaction entity.Transaction) (entity.Transaction, error)
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

func (r *transactionRepository) GetByUserID(userID int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.connection.
		// Accessing indirect relationship
		Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").
		Where("user_id = ?", userID).
		Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (r *transactionRepository) Save(transaction entity.Transaction) (entity.Transaction, error) {
	err := r.connection.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *transactionRepository) Update(transaction entity.Transaction) (entity.Transaction, error) {
	err := r.connection.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
