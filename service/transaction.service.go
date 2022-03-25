package service

import (
	"errors"
	"fmt"
	"go-fundraising/dto"
	"go-fundraising/entity"
	"go-fundraising/repository"
	"math/rand"
	"strconv"
)

type TransactionService interface {
	GetTransactionsByCampaignID(request dto.TransactionGetRequestID) ([]entity.Transaction, error)
	GetTransactionsByUserID(userID int) ([]entity.Transaction, error)
	CreateTransaction(request dto.TransactionCreateRequest) (entity.Transaction, error)
	ProcessPayment(request dto.TransactionNotificationRequest) error
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	campaignRepository    repository.CampaignRepository
	paymentService        PaymentService
}

func NewTransactionService(transactionRepository repository.TransactionRepository, campaignRepository repository.CampaignRepository, paymentService PaymentService) TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
		campaignRepository:    campaignRepository,
		paymentService:        paymentService,
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

func (s *transactionService) CreateTransaction(request dto.TransactionCreateRequest) (entity.Transaction, error) {
	transaction := entity.Transaction{}
	transaction.CampaignID = request.CampaignID
	transaction.Amount = request.Amount
	transaction.UserID = int(request.User.ID)
	transaction.Status = "pending"
	transaction.Code = fmt.Sprintf("TRX-%d", rand.Intn(100000))
	newTransaction, err := s.transactionRepository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentURL, err := s.paymentService.GetPaymentURL(request.User, newTransaction)
	if err != nil {
		return newTransaction, err
	}
	newTransaction.PaymentURL = paymentURL
	newTransaction, err = s.transactionRepository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (s *transactionService) ProcessPayment(request dto.TransactionNotificationRequest) error {
	transaction_id, _ := strconv.Atoi(request.OrderID)
	transaction, err := s.transactionRepository.GetByID(transaction_id)
	if err != nil {
		return err
	}

	if request.PaymentType == "credit_card" && request.TransactionStatus == "capture" && request.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if request.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if request.TransactionStatus == "deny" || request.TransactionStatus == "expire" || request.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}
	updatedTransaction, err := s.transactionRepository.Update(transaction)
	if err != nil {
		return err
	}
	campaign, err := s.campaignRepository.FindCampaignByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}
	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount
		_, err = s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}
