package handler

import (
	"github.com/gin-gonic/gin"
	"go-fundraising/dto"
	"go-fundraising/service"
	"net/http"
)

type UserHandler interface {
	Index(ctx *gin.Context)
	Add(ctx *gin.Context)
	Store(ctx *gin.Context)
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

func (h *userHandler) Add(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "add_user.html", nil)
}

func (h *userHandler) Store(ctx *gin.Context) {
	var requestForm dto.FormStoreUserRequest
	err := ctx.ShouldBind(&requestForm)
	if err != nil {
		requestForm.Error = err.Error()
		ctx.HTML(http.StatusBadRequest, "add_user.html", requestForm)
		return
	}
	// Convert requestUserForm to registerRequestJSON to use RegisterUser service
	registerRequest := dto.RegisterRequest{}
	registerRequest.Name = requestForm.Name
	registerRequest.Email = requestForm.Email
	registerRequest.Occupation = requestForm.Occupation
	registerRequest.Password = requestForm.Password
	_, err = h.userService.RegisterUser(registerRequest)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	ctx.Redirect(http.StatusFound, "/users")
}
