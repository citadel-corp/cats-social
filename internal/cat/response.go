package cat

import "time"

type CreateCatResponse struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CatResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Race        string    `json:"race"`
	Sex         string    `json:"sex"`
	AgeInMonth  int       `json:"ageInMonth"`
	ImageUrls   []string  `json:"imageUrls"`
	Description string    `json:"description"`
	HasMatched  bool      `json:"hasMatched"`
	CreatedAt   time.Time `json:"createdAt"`
}
