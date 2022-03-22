package helper

import (
	"go-fundraising/entity"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type UserFormatter struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
}

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

type CampaignDetailFormatter struct {
	ID               int                            `json:"id"`
	Name             string                         `json:"name"`
	ShortDescription string                         `json:"short_description"`
	Description      string                         `json:"description"`
	ImageURL         string                         `json:"image_url"`
	GoalAmount       int                            `json:"goal_amount"`
	CurrentAmount    int                            `json:"current_amount"`
	UserID           int                            `json:"user_id"`
	Slug             string                         `json:"slug"`
	Perks            []string                       `json:"perks"`
	User             CampaignDetailUserFormatter    `json:"user"`
	Images           []CampaignDetailImageFormatter `json:"images"`
}

type CampaignDetailUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignDetailImageFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTransactionFormatter struct {
	ID        int                   `json:"id"`
	Amount    int                   `json:"amount"`
	Status    string                `json:"status"`
	CreatedAt time.Time             `json:"created_at"`
	Campaign  UserCampaignFormatter `json:"campaign"`
}

type UserCampaignFormatter struct {
	Name     string `json::"name"`
	ImageURL string `json:"image_url"`
}

func FormatUser(user entity.User, token string) UserFormatter {
	userFormatter := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		Token:      token,
	}
	return userFormatter
}

func FormatCampaign(campaign entity.Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserId
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.ImageURL = ""
	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}
	return campaignFormatter
}

func FormatCampaigns(campaigns []entity.Campaign) []CampaignFormatter {
	campaignsFormatter := []CampaignFormatter{}
	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}
	return campaignsFormatter
}

func FormatValidationErrors(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}

func FormatCampaignDetail(campaign entity.Campaign) CampaignDetailFormatter {
	campaignDetailFormatter := CampaignDetailFormatter{}
	campaignDetailFormatter.ID = campaign.ID
	campaignDetailFormatter.UserID = campaign.UserId
	campaignDetailFormatter.Name = campaign.Name
	campaignDetailFormatter.ShortDescription = campaign.ShortDescription
	campaignDetailFormatter.Description = campaign.Description
	campaignDetailFormatter.ImageURL = ""
	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}
	var perks []string
	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk))
	}
	campaignDetailFormatter.Perks = perks

	// CampaignDetailUserFormatter
	user := campaign.User
	campaignDetailUserFormatter := CampaignDetailUserFormatter{}
	campaignDetailUserFormatter.Name = user.Name
	campaignDetailUserFormatter.ImageURL = user.AvatarFileName
	campaignDetailFormatter.User = campaignDetailUserFormatter

	// CampaignDetailImageFormatter
	images := []CampaignDetailImageFormatter{}
	for _, image := range campaign.CampaignImages {
		campaignDetailImageFormatter := CampaignDetailImageFormatter{}
		campaignDetailImageFormatter.ImageURL = image.FileName
		isPrimary := false
		if image.IsPrimary == 1 {
			isPrimary = true
		}
		campaignDetailImageFormatter.IsPrimary = isPrimary
		images = append(images, campaignDetailImageFormatter)
	}
	campaignDetailFormatter.Images = images

	return campaignDetailFormatter
}

func FormatCampaignTransaction(transaction entity.Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	return formatter
}

func FormatCampaignTransactions(transactions []entity.Transaction) []CampaignTransactionFormatter {
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}
	var transactionsFormatter []CampaignTransactionFormatter
	for _, transaction := range transactions {
		transactionFormatter := FormatCampaignTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, transactionFormatter)
	}
	return transactionsFormatter
}

func FormatUserTransaction(transaction entity.Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	campaignFormatter := UserCampaignFormatter{}
	campaignFormatter.Name = transaction.User.Name
	campaignFormatter.ImageURL = ""
	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}
	formatter.Campaign = campaignFormatter

	return formatter
}

func FormatUserTransactions(transactions []entity.Transaction) []UserTransactionFormatter {
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}

	var userTransactionFormatter []UserTransactionFormatter
	for _, transaction := range transactions {
		formatter := FormatUserTransaction(transaction)
		userTransactionFormatter = append(userTransactionFormatter, formatter)
	}

	return userTransactionFormatter
}
