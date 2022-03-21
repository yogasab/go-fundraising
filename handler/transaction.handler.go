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
