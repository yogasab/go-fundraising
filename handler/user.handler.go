package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-fundraising/dto"
	"go-fundraising/helper"
	"go-fundraising/service"
	"net/http"
)

type UserHandler interface {
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	CheckEmailAvaibility(ctx *gin.Context)
	UploadAvatar(ctx *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) RegisterUser(ctx *gin.Context) {
	var request dto.RegisterRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errors := helper.FormatValidationErrors(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to process request", http.StatusUnprocessableEntity, "failed", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	user, err := h.userService.RegisterUser(request)
	if err != nil {
		response := helper.APIResponse("Failed to process request", http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	userResponse := helper.FormatUser(user, "fafifu")
	response := helper.APIResponse("User registered successfully", http.StatusCreated, "success", userResponse)
	ctx.JSON(http.StatusCreated, response)
}

func (h *userHandler) LoginUser(ctx *gin.Context) {
	var request dto.LoginRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errors := helper.FormatValidationErrors(err)
		errorMessages := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to process request", http.StatusUnprocessableEntity, "failed", errorMessages)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user, errUser := h.userService.LoginUser(request)
	if errUser != nil {
		errorMessage := gin.H{"errors": errUser.Error()}
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "failed", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("User login successfully", http.StatusOK, "sucess", user)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvaibility(ctx *gin.Context) {
	var request dto.CheckEmailRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errors := helper.FormatValidationErrors(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to process request", http.StatusUnprocessableEntity, "failed", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	isEmailAvailable, err := h.userService.CheckEmailAvailability(request)
	if err != nil {
		errorMessage := gin.H{"errors": "Error on finding user email"}
		response := helper.APIResponse("Failed to process request", http.StatusUnprocessableEntity, "failed", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	data := gin.H{
		"is_available": isEmailAvailable,
	}
	metaMessage := "Email is already registered, please try another"
	if isEmailAvailable {
		metaMessage = "Email is available you can continue the registration"
	}
	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Error on upload avatar", http.StatusBadRequest, "failed", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	userID := 1
	destination := fmt.Sprintf("images/avatars/%d-%s", userID, file.Filename)
	err = ctx.SaveUploadedFile(file, destination)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, "failed", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.userService.SaveAvatar(16, destination)
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar uploaded successfully", http.StatusOK, "success", data)
	ctx.JSON(http.StatusOK, response)
}
