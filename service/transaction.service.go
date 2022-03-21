package service

import (
	"errors"
	"go-fundraising/dto"
	"go-fundraising/entity"
	"go-fundraising/repository"
)

type TransactionService interface {
	GetTransactionsByCampaignID(request dto.TransactionGetRequestID) ([]entity.Transaction, error)
	GetTransactionsByUserID(userID int) ([]entity.Transaction, error)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	campaignRepository    repository.CampaignRepository
}

func NewTransactionService(transactionRepository repository.TransactionRepository, campaignRepository repository.CampaignRepository) TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
		campaignRepository:    campaignRepository,
	}
}

func (s *transactionService) GetTransactionsByCampaignID(request dto.TransactionGetRequestID) ([]entity.Transaction, error) {
	// ID referred to CampaignID from params
	campaign, err := s.campaignRepository.FindCampaignByID(request.ID)
	if err != nil {
		return []entity.Transaction{}, err
	}

	if campaign.UserId != int(request.User.ID) {
		return []entity.Transaction{}, errors.New("You are not be able to perform this route")
	}

	transactions, err := s.transactionRepository.GetByCampaignID(request.ID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *transactionService) GetTransactionsByUserID(userID int) ([]entity.Transaction, error) {
	transactions, err := s.transactionRepository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
