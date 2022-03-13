package handler

import (
	"github.com/gin-gonic/gin"
	"go-fundraising/helper"
	"go-fundraising/service"
	"net/http"
	"strconv"
)

type CampaignHandler interface {
	GetCampaigns(ctx *gin.Context)
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
