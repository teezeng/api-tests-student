package models

type CreateAdvertisementRequest struct {
	Title       string
	Description string
	Price       float64
	Quantity    int
	Photos      [][]byte
}

type Photo struct {
	ID              string `json:"id"`
	AdvertisementID string `json:"advertisement_id,omitempty"`
	URL             string `json:"url"`
	SortOrder       int    `json:"sort_order"`
	CreatedAt       string `json:"created_at"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AdvertisementResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Photos      []Photo `json:"photos"`
	UserID      string  `json:"user_id"`
	User        User    `json:"user"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type AdvertisementsList struct {
	Items      []AdvertisementResponse `json:"items"`
	HasNext    bool                    `json:"has_next"`
	NextCursor string                  `json:"next_cursor"`
	Total      int                     `json:"total"`
}

func PhotoIDs(photos []Photo) []string {
	ids := make([]string, 0, len(photos))
	for _, photo := range photos {
		ids = append(ids, photo.ID)
	}
	return ids
}
