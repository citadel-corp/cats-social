package catmatch

import validation "github.com/go-ozzo/ozzo-validation/v4"

type PostCatMatch struct {
	MatchCatId string `json:"matchCatId"`
	UserCatId  string `json:"userCatId"`
	Message    string `json:"message"`
}

func (p PostCatMatch) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.MatchCatId, validation.Required),
		validation.Field(&p.UserCatId, validation.Required),
		validation.Field(&p.Message, validation.Required, validation.Length(5, 120)),
	)
}

type ApproveOrRejectMatch struct {
	MatchUID string `json:"matchId"`
}

func (p ApproveOrRejectMatch) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.MatchUID, validation.Required),
	)
}

type DeleteMatch struct {
	MatchUID string `schema:"id" binding:"omitempty"`
}

func (p DeleteMatch) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.MatchUID, validation.Required),
	)
}
