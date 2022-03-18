package handler

import (
	"github.com/gin-gonic/gin"
	"go-fundraising/dto"
	"go-fundraising/entity"
	"go-fundraising/helper"
	"go-fundraising/service"
	"net/http"
)

type TransactionHandler interface {
	GetTransactionsByCampaignID(ctx *gin.Context)
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
