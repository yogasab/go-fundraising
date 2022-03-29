package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-fundraising/dto"
	"go-fundraising/service"
	"net/http"
)

type CampaignHandler interface {
	Index(ctx *gin.Context)
	Add(ctx *gin.Context)
	Store(ctx *gin.Context)
}

type campaignHandler struct {
	campaignService service.CampaignService
	userService     service.UserService
}

func NewCampaignHandler(campaignService service.CampaignService, userService service.UserService) CampaignHandler {
	return &campaignHandler{
		campaignService: campaignService,
		userService:     userService,
	}
}

func (h *campaignHandler) Index(ctx *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		fmt.Println(err.Error())
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	ctx.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

func (h *campaignHandler) Add(ctx *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	request := dto.FormStoreCampaignRequest{}
	request.Users = users
	ctx.HTML(http.StatusOK, "add_campaign.html", request)
}

func (h *campaignHandler) Store(ctx *gin.Context) {
	var request dto.FormStoreCampaignRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		users, e := h.userService.GetAllUsers()
		if e != nil {
			ctx.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
		request.Error = err
		request.Users = users
		ctx.HTML(http.StatusInternalServerError, "error.html", request)
		return
	}
	user, err := h.userService.GetUserByID(request.UserID)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	formRequest := dto.CreateCampaignRequest{}
	formRequest.Name = request.Name
	formRequest.ShortDescription = request.ShortDescription
	formRequest.Description = request.Description
	formRequest.Perks = request.Perks
	formRequest.GoalAmount = request.GoalAmount
	formRequest.User = user
	_, err = h.campaignService.CreateCampaign(formRequest)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
	}
	ctx.Redirect(http.StatusFound, "/campaigns")
}
