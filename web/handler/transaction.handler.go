package handler

import (
	"github.com/gin-gonic/gin"
	"go-fundraising/service"
	"net/http"
)

type TransactionHandler interface {
	Index(ctx *gin.Context)
}

type transactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(tansactionService service.TransactionService) TransactionHandler {
	return &transactionHandler{
		transactionService: tansactionService,
	}
}

func (h *transactionHandler) Index(ctx *gin.Context) {
	transactions, err := h.transactionService.GetAllTransactions()
	if err != nil {
		ctx.Redirect(http.StatusInternalServerError, "/transactions")
	}
	ctx.HTML(http.StatusOK, "transaction_index.html", gin.H{"transactions": transactions})
}
