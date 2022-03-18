package dto

import "go-fundraising/entity"

type TransactionGetRequestID struct {
	ID   int `uri:"id" binding:"required"`
	User entity.User
}
