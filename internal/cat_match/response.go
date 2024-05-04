package catmatch

import (
	"time"

	"github.com/citadel-corp/cats-social/internal/cat"
)

type CatMatchResponse struct {
	ID             string          `json:"id"`
	IssuedBy       Issuer          `json:"issuedBy"`
	MatchCatDetail cat.CatResponse `json:"matchCatDetail"`
	UserCatDetail  cat.CatResponse `json:"userCatDetail"`
	Message        string          `json:"message"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type Issuer struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
