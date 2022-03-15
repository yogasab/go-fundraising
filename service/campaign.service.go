package service

import (
	"fmt"
	"github.com/gosimple/slug"
	"go-fundraising/dto"
	"go-fundraising/entity"
	"go-fundraising/repository"
)

type CampaignService interface {
	GetCampaigns(userID int) ([]entity.Campaign, error)
	GetCampaignByID(request dto.CampaignGetRequestID) (entity.Campaign, error)
	GetCampaignBySlug(request dto.CampaignGetRequestSlug) (entity.Campaign, error)
	CreateCampaign(request dto.CreateCampaignRequest) (entity.Campaign, error)
}

type campaignService struct {
	campaignRepository repository.CampaignRepository
}

func NewCampaignService(campaignRepository repository.CampaignRepository) CampaignService {
	return &campaignService{campaignRepository: campaignRepository}
}

func (s *campaignService) GetCampaigns(userID int) ([]entity.Campaign, error) {
	if userID != 0 {
		campaigns, err := s.campaignRepository.FindCampaignByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
	campaigns, err := s.campaignRepository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s *campaignService) GetCampaignByID(request dto.CampaignGetRequestID) (entity.Campaign, error) {
	campaign, err := s.campaignRepository.FindCampaignByID(request.ID)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *campaignService) GetCampaignBySlug(request dto.CampaignGetRequestSlug) (entity.Campaign, error) {
	campaign, err := s.campaignRepository.FindCampaignBySlug(request.Slug)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *campaignService) CreateCampaign(request dto.CreateCampaignRequest) (entity.Campaign, error) {
	campaign := entity.Campaign{}
	campaign.Name = request.Name
	campaign.ShortDescription = request.ShortDescription
	campaign.Description = campaign.Description
	campaign.GoalAmount = request.GoalAmount
	campaign.Perks = request.Perks
	campaign.UserId = int(request.User.ID)
	combinedSlug := fmt.Sprintf("%s %d", campaign.Name, int(request.User.ID))
	campaign.Slug = slug.Make(combinedSlug)

	newCampaign, err := s.campaignRepository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}
