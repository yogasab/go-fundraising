package service

import (
	"errors"
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
	UpdateCampaign(requestID dto.CampaignGetRequestID, requestCampaign dto.CreateCampaignRequest) (entity.Campaign, error)
	CreateCampaignImage(request dto.CreateCampaignImageRequest, fileLocation string) (entity.CampaignImage, error)
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

func (s *campaignService) UpdateCampaign(requestID dto.CampaignGetRequestID, requestCampaign dto.CreateCampaignRequest) (entity.Campaign, error) {
	campaign, err := s.campaignRepository.FindCampaignByID(requestID.ID)
	if err != nil {
		return campaign, err
	}
	if campaign.UserId != int(requestCampaign.User.ID) {
		return campaign, errors.New("You are not be able to perform this route")
	}
	campaign.Name = requestCampaign.Name
	campaign.ShortDescription = requestCampaign.ShortDescription
	campaign.Description = requestCampaign.Description
	campaign.GoalAmount = requestCampaign.GoalAmount
	campaign.Perks = requestCampaign.Perks

	updatedCampaign, err := s.campaignRepository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}
	return updatedCampaign, nil
}

func (s *campaignService) CreateCampaignImage(request dto.CreateCampaignImageRequest, filename string) (entity.CampaignImage, error) {
	campaign, err := s.campaignRepository.FindCampaignByID(request.CampaignID)
	if err != nil {
		return entity.CampaignImage{}, err
	}
	if campaign.UserId != int(request.User.ID) {
		fmt.Println(campaign.UserId)
		fmt.Println(request.User.ID)
		return entity.CampaignImage{}, errors.New("You are not be able to perform this route")
	}

	isPrimary := 0
	if request.IsPrimary {
		isPrimary = 1
		_, err := s.campaignRepository.MarkAllImagesAsNonPrimary(request.CampaignID)
		if err != nil {
			return entity.CampaignImage{}, err
		}
	}
	campaignImage := entity.CampaignImage{}
	campaignImage.CampaignID = request.CampaignID
	campaignImage.FileName = filename
	campaignImage.IsPrimary = isPrimary

	newCampaignImage, err := s.campaignRepository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}
	return newCampaignImage, nil
}
