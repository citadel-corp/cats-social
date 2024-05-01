package cat

import "time"

type CreateCatResponse struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
