package catmatch

import "time"

type MatchStatus string

const (
	Pending  MatchStatus = "pending"
	Approved MatchStatus = "approved"
	Rejected MatchStatus = "rejected"
)

type CatMatches struct {
	ID             int64       `json:"id"`
	UID            string      `json:"uid"`
	IssuerCatId    int64       `json:"issuer_cat_id"`
	IssueUserId    int64       `json:"issue_user_id"`
	MatchCatId     int64       `json:"match_cat_id"`
	MatchUserId    int64       `json:"match_user_id"`
	Message        string      `json:"message"`
	ApprovalStatus MatchStatus `json:"approval_status"`
	CreatedAt      *time.Time  `json:"created_at"`
}
