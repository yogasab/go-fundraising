package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler interface {
	Index(ctx *gin.Context)
}

type userHandler struct {
}

func NewUserHandler() UserHandler {
	return &userHandler{}
}

func (h *userHandler) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index_user.html", nil)
}
