package dto

type CampaignGetRequestID struct {
	ID int `uri:"id"`
}

type CampaignGetRequestSlug struct {
	Slug string `uri:"slug" binding:"required"`
}
