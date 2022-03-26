package handler

import (
	"github.com/gin-gonic/gin"
	"go-fundraising/service"
	"net/http"
)

type UserHandler interface {
	Index(ctx *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) Index(ctx *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	ctx.HTML(http.StatusOK, "index_user.html", gin.H{"users": users})
}
