package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-fundraising/dto"
	"go-fundraising/entity"
	"go-fundraising/helper"
	"go-fundraising/service"
	"net/http"
	"strconv"
)

type CampaignHandler interface {
	GetCampaigns(ctx *gin.Context)
	GetCampaignByID(ctx *gin.Context)
	GetCampaignBySlug(ctx *gin.Context)
	CreateCampaign(ctx *gin.Context)
	UpdateCampaign(ctx *gin.Context)
	UploadCampaignImage(ctx *gin.Context)
}

type campaignHandler struct {
	campaignService service.CampaignService
}

func NewCampaignHandler(campaignService service.CampaignService) CampaignHandler {
	return &campaignHandler{
		campaignService: campaignService,
	}
}

func (h *campaignHandler) GetCampaigns(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.Query("user_id"))

	campaigns, err := h.campaignService.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Failed to find campaigns", http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign fetched successfully", http.StatusOK, "success", helper.FormatCampaigns(campaigns))
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaignByID(ctx *gin.Context) {
	var request dto.CampaignGetRequestID
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		response := helper.APIResponse("Failed to process request", http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadGateway, response)
		return
	}
	campaign, err := h.campaignService.GetCampaignByID(request)
	if err != nil {
		response := helper.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign fetched successfully", http.StatusOK, "success", helper.FormatCampaignDetail(campaign))
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaignBySlug(ctx *gin.Context) {
	var request dto.CampaignGetRequestSlug
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		response := helper.APIResponse("Failed to process request", http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	campaign, err := h.campaignService.GetCampaignBySlug(request)
	if err != nil {
		response := helper.APIResponse("Failed to process request", http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign fetched successfully", http.StatusOK, "success", helper.FormatCampaignDetail(campaign))
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(ctx *gin.Context) {
	var request dto.CreateCampaignRequest
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
	newCampaign, err := h.campaignService.CreateCampaign(request)
	if err != nil {
		response := helper.APIResponse("Failed to process request", http.StatusBadRequest, "failed", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign created successfully", http.StatusCreated, "success", helper.FormatCampaign(newCampaign))
	ctx.JSON(http.StatusCreated, response)
}

func (h *campaignHandler) UpdateCampaign(ctx *gin.Context) {
	var requestID dto.CampaignGetRequestID
	err := ctx.ShouldBindUri(&requestID)
	if err != nil {
		response := helper.APIResponse("Failed to process request", http.StatusBadRequest, "failed", err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var requestCampaign dto.CreateCampaignRequest
	err = ctx.ShouldBindJSON(&requestCampaign)
	if err != nil {
		errors := helper.FormatValidationErrors(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to process request", http.StatusBadRequest, "failed", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	user := ctx.MustGet("user").(entity.User)
	requestCampaign.User = user
	campaign, err := h.campaignService.UpdateCampaign(requestID, requestCampaign)
	if err != nil {
		response := helper.APIResponse("Failed to process request", http.StatusBadRequest, "failed", err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign updated successfully", http.StatusOK, "success", campaign)
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadCampaignImage(ctx *gin.Context) {
	var request dto.CreateCampaignImageRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		errors := helper.FormatValidationErrors(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to process request", http.StatusUnprocessableEntity, "failed", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user := ctx.MustGet("user").(entity.User)
	request.User = user
	userID := int(user.ID)
	file, err := ctx.FormFile("campaign_image")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "failed", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	filename := fmt.Sprintf("images/campaigns/%d-%s", userID, file.Filename)
	err = ctx.SaveUploadedFile(file, filename)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "failed", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.campaignService.CreateCampaignImage(request, filename)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "failed", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image uploaded successfully", http.StatusCreated, "success", data)
	ctx.JSON(http.StatusCreated, response)
}
