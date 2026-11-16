package models

import advModels "api-tests-template/internal/managers/advertisements/models"

type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

type FavoriteItem struct {
	AddedAt       string                          `json:"added_at"`
	Advertisement advModels.AdvertisementResponse `json:"advertisement"`
	FavoriteID    string                          `json:"favorite_id"`
}

type FavoriteListResponse struct {
	Items      []FavoriteItem `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}

type FavoriteResponse struct {
	AddedAt       string                          `json:"added_at"`
	Advertisement advModels.AdvertisementResponse `json:"advertisement"`
	FavoriteID    string                          `json:"favorite_id"`
}
