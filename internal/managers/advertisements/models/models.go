package models

type Advertisement struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserId      string `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}
