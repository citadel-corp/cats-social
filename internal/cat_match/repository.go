package catmatch

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/cats-social/internal/common/db"
	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, catMatch *CatMatches) error
	Approve(ctx context.Context, catMatch *CatMatches) error
	Reject(ctx context.Context, catMatch *CatMatches) error
	Delete(ctx context.Context, id int64, userId int64) error
	// GetByCatID(ctx context.Context, catID int64) (*CatMatches, error)
	GetByUIDAndUserID(ctx context.Context, uid string, userID int64, filter map[string]interface{}) (*CatMatches, error)
	// GetMatchingCats(ctx context.Context, matchUid string) (*CatMatchAndCats, error)
	List(ctx context.Context, userID int64, filter map[string]interface{}) ([]CatMatchList, error)
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
		WHERE uid = $1 AND issuer_user_id = $2;
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

// List implements Repository.
func (d *dbRepository) List(ctx context.Context, userID int64, filter map[string]interface{}) ([]CatMatchList, error) {
	listQuery := `
		SELECT cm.uid, cm.message, cm.created_at,
		ic.uid, ic.name, ic.race, ic.sex, ic.description, ic.age_in_month,
		ic.image_urls, ic.has_matched, ic.created_at,
		mc.uid, mc.name, mc.race, mc.sex, mc.description, mc.age_in_month,
		mc.image_urls, mc.has_matched, mc.created_at,
		u.id, u.name, u.email, u.created_at
		FROM cat_matches cm
		LEFT JOIN cats ic on cm.issuer_cat_id = ic.id
		LEFT JOIN cats mc on cm.matched_cat_id = mc.id
		LEFT JOIN users u on cm.issuer_user_id = u.id
		WHERE (cm.issuer_user_id = $1 OR cm.matched_user_id = $1)
		ORDER BY cm.created_at desc;
	`

	// var approvalStatus MatchStatus
	// if v, ok := filter["approval_status"].(MatchStatus); ok {
	// 	approvalStatus = v
	// } else {
	// 	approvalStatus = Approved
	// }

	rows, err := d.db.DB().QueryContext(ctx, listQuery, userID)
	if err != nil {
		return nil, err
	}
	res := make([]CatMatchList, 0)
	for rows.Next() {
		catMatch := CatMatchList{}
		err = rows.Scan(&catMatch.ID, &catMatch.Message, &catMatch.CreatedAt,
			&catMatch.IssuerCat.ID, &catMatch.IssuerCat.Name, &catMatch.IssuerCat.Race, &catMatch.IssuerCat.Sex, &catMatch.IssuerCat.Description, &catMatch.IssuerCat.AgeInMonth,
			pq.Array(&catMatch.IssuerCat.ImageUrls), &catMatch.IssuerCat.HasMatched, &catMatch.IssuerCat.CreatedAt,
			&catMatch.MatchCat.ID, &catMatch.MatchCat.Name, &catMatch.MatchCat.Race, &catMatch.MatchCat.Sex, &catMatch.MatchCat.Description, &catMatch.MatchCat.AgeInMonth,
			pq.Array(&catMatch.MatchCat.ImageUrls), &catMatch.MatchCat.HasMatched, &catMatch.MatchCat.CreatedAt,
			&catMatch.IssuedBy.ID, &catMatch.IssuedBy.Name, &catMatch.IssuedBy.Email, &catMatch.IssuedBy.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, catMatch)
	}
	return res, nil
}
