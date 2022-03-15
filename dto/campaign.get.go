package dto

type CampaignGetRequestID struct {
	ID int `uri:"id" binding:"required"`
}

type CampaignGetRequestSlug struct {
	Slug string `uri:"slug" binding:"required"`
}
