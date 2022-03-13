package service

import (
	"go-fundraising/entity"
	"go-fundraising/repository"
)

type CampaignService interface {
	FindCampaigns(userID int) ([]entity.Campaign, error)
}

type campaignService struct {
	campaignRepository repository.CampaignRepository
}

func NewCampaignService(campaignRepository repository.CampaignRepository) CampaignService {
	return &campaignService{campaignRepository: campaignRepository}
}

func (s *campaignService) FindCampaigns(userID int) ([]entity.Campaign, error) {
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
