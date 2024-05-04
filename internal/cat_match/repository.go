package catmatch

import (
	"context"

	"github.com/citadel-corp/cats-social/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, catMatch *CatMatches) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, catMatch *CatMatches) error {
	createCatQuery := `
        INSERT INTO cat_matches (
            uid, issuer_cat_id, issuer_user_id, matched_cat_id, matched_user_id, message
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        );
    `

	_, err := d.db.DB().ExecContext(ctx, createCatQuery,
		catMatch.UID, catMatch.IssuerCatId, catMatch.IssueUserId,
		catMatch.MatchCatId, catMatch.MatchUserId, catMatch.Message)
	if err != nil {
		return err
	}

	return nil
}
