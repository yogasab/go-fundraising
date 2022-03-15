package handler

import (
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
