package catmatch

import (
	"context"

	"github.com/citadel-corp/cats-social/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, catMatch *CatMatches) error
	List(ctx context.Context, userID int64) ([]*CatMatches, error)
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

// List implements Repository.
func (d *dbRepository) List(ctx context.Context, userID int64) ([]*CatMatches, error) {
	_ = `
		SELECT cm.*, ic.*, mc.*, u.*
		FROM cat_matches cm
		LEFT JOIN cats ic on cm.issuer_user_id = ic.user_id
		LEFT JOIN cats mc on cm.matched_user_id = mc.user_id
		LEFT JOIN users u on cm.issuer_user_id = u.id
		WHERE cm.issuer_user_id = $1
	`
	// todo lanjutkan

	return nil, nil
}
