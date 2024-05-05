package catmatch

import (
	"time"

	"github.com/citadel-corp/cats-social/internal/cat"
)

type CatMatchList struct {
	ID        string
	IssuedBy  Issuer
	MatchCat  cat.CatResponse
	IssuerCat cat.CatResponse
	Message   string
	CreatedAt time.Time
}

type CatMatchResponse struct {
	ID             string          `json:"id"`
	IssuedBy       Issuer          `json:"issuedBy"`
	MatchCatDetail cat.CatResponse `json:"matchCatDetail"`
	UserCatDetail  cat.CatResponse `json:"userCatDetail"`
	Message        string          `json:"message"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type Issuer struct {
	ID        int64     `json:"-"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func MakeCatMatchResponse(list []CatMatchList, userId int64) []CatMatchResponse {
	res := []CatMatchResponse{}
	for _, match := range list {
		userCat := cat.CatResponse{}
		matchCat := cat.CatResponse{}
		if match.IssuedBy.ID == userId {
			userCat = match.IssuerCat
			matchCat = match.MatchCat
		} else {
			userCat = match.MatchCat
			matchCat = match.IssuerCat
		}

		res = append(res, CatMatchResponse{
			ID:             match.ID,
			IssuedBy:       match.IssuedBy,
			MatchCatDetail: matchCat,
			UserCatDetail:  userCat,
			Message:        match.Message,
			CreatedAt:      match.CreatedAt,
		})
	}

	return res
}
