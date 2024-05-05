package catmatch

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/cats-social/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, catMatch *CatMatches) error
	Approve(ctx context.Context, catMatch *CatMatches) error
	Reject(ctx context.Context, catMatch *CatMatches) error
	Delete(ctx context.Context, id int64, userId int64) error
	// GetByCatID(ctx context.Context, catID int64) (*CatMatches, error)
	GetByUIDAndUserID(ctx context.Context, uid string, userID int64, filter map[string]interface{}) (*CatMatches, error)
	// GetMatchingCats(ctx context.Context, matchUid string) (*CatMatchAndCats, error)
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

func (d *dbRepository) Approve(ctx context.Context, catMatch *CatMatches) error {
	err := d.db.StartTx(ctx, func(tx *sql.Tx) error {
		updateMatchQuery := `
			UPDATE cat_matches
			SET approval_status = $1
			WHERE id = $2;
		`

		_, err := d.db.DB().ExecContext(ctx, updateMatchQuery, Approved, catMatch.ID)
		if err != nil {
			return err
		}

		// update cats' matched status
		updateCatQuery := `
			UPDATE cats
			SET has_matched = true
			WHERE id = $1;
		`
		_, err = d.db.DB().ExecContext(ctx, updateCatQuery, catMatch.IssuerCatId)
		if err != nil {
			return err
		}

		_, err = d.db.DB().ExecContext(ctx, updateCatQuery, catMatch.MatchCatId)
		if err != nil {
			return err
		}

		// reject each cat's remaining matches
		deleteMatchQuery := `
            UPDATE cat_matches
			SET approval_status = $1
            WHERE (issuer_cat_id = $2 OR matched_cat_id = $2) AND id != $3;
        `
		_, err = d.db.DB().ExecContext(ctx, deleteMatchQuery, Rejected, catMatch.IssuerCatId, catMatch.ID)
		if err != nil {
			return err
		}

		_, err = d.db.DB().ExecContext(ctx, deleteMatchQuery, Rejected, catMatch.MatchCatId, catMatch.ID)
		if err != nil {
			return err
		}

		return nil

	})

	return err
}

func (d *dbRepository) Reject(ctx context.Context, catMatch *CatMatches) error {
	updateMatchQuery := `
		UPDATE cat_matches
		SET approval_status = $1
		WHERE id = $2;
	`

	_, err := d.db.DB().ExecContext(ctx, updateMatchQuery, Rejected, catMatch.ID)
	if err != nil {
		return err
	}

	return nil
}

func (d *dbRepository) GetByUIDAndUserID(ctx context.Context, uid string, userID int64, filter map[string]interface{}) (*CatMatches, error) {
	getMatchQuery := `
		SELECT id, uid, issuer_cat_id, issuer_user_id, matched_cat_id, matched_user_id, message, approval_status
		FROM cat_matches
		WHERE uid = $1 AND matched_user_id = $2;
	`

	row := d.db.DB().QueryRowContext(ctx, getMatchQuery, uid, userID)
	catMatch := &CatMatches{}
	err := row.Scan(&catMatch.ID, &catMatch.UID, &catMatch.IssuerCatId, &catMatch.IssueUserId,
		&catMatch.MatchCatId, &catMatch.MatchUserId, &catMatch.Message, &catMatch.ApprovalStatus)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCatMatchNotFound
	}
	if err != nil {
		return nil, err
	}

	if v, ok := filter["pending_only"].(bool); v && ok && catMatch.ApprovalStatus != Pending {
		return nil, ErrCatMatchNoLongerValid
	}

	return catMatch, nil
}

func (d *dbRepository) Delete(ctx context.Context, id int64, userId int64) error {
	deleteMatchQuery := `
		DELETE FROM cat_matches
		WHERE id = $1 AND issuer_user_id = $2;
	`

	_, err := d.db.DB().ExecContext(ctx, deleteMatchQuery, id, userId)
	if err != nil {
		return err
	}

	return nil
}
