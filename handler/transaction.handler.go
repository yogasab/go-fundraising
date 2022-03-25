package handler

import (
	"fmt"
	"go-fundraising/dto"
	"go-fundraising/entity"
	"go-fundraising/helper"
	"go-fundraising/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler interface {
	GetTransactionsByCampaignID(ctx *gin.Context)
	GetTransactionsByUserID(ctx *gin.Context)
	CreateTransaction(ctx *gin.Context)
	GetNotification(ctx *gin.Context)
}

type transactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) TransactionHandler {
	return &transactionHandler{
		transactionService: transactionService,
	}
}

func (h *transactionHandler) GetTransactionsByCampaignID(ctx *gin.Context) {
	var request dto.TransactionGetRequestID
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		response := helper.APIResponse("Failed to process request", http.StatusUnprocessableEntity, "failed", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	user := ctx.MustGet("user").(entity.User)
	request.User = user

	transactions, err := h.transactionService.GetTransactionsByCampaignID(request)
	if err != nil {
		response := helper.APIResponse("Failed to get transactions", http.StatusBadRequest, "failed", err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Transactions fetched successfully", http.StatusOK, "success", helper.FormatCampaignTransactions(transactions))
	ctx.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetTransactionsByUserID(ctx *gin.Context) {
	user := ctx.MustGet("user").(entity.User)
	fmt.Println(user)
	userID := user.ID
	transactions, err := h.transactionService.GetTransactionsByUserID(int(userID))
	if err != nil {
		response := helper.APIResponse("Failed to process request", http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Transactions fetched succesfully", http.StatusOK, "success", helper.FormatUserTransactions(transactions))
	ctx.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(ctx *gin.Context) {
	var request dto.TransactionCreateRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errors := helper.FormatValidationErrors(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to process request", http.StatusUnprocessableEntity, "failed", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	user := ctx.MustGet("user").(entity.User)
	request.User = user
	transaction, err := h.transactionService.CreateTransaction(request)
	if err != nil {
		response := helper.APIResponse("Failed to create new transaction", http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Transaction created successfully",
		http.StatusCreated, "success",
		helper.FormatTransaction(transaction))
	ctx.JSON(http.StatusCreated, response)
}

func (h *transactionHandler) GetNotification(ctx *gin.Context) {
	var request dto.TransactionNotificationRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		response := helper.APIResponse("Failed to process request", http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	err = h.transactionService.ProcessPayment(request)
	if err != nil {
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	ctx.JSON(http.StatusOK, request)
}
