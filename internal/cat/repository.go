package cat

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/cats-social/internal/common/db"
	"github.com/lib/pq"
)

type Repository interface {
	GetByIDAndUserID(ctx context.Context, id string, userID string) (*Cat, error)
	Create(ctx context.Context, cat *Cat) error
	Update(ctx context.Context, cat *Cat) error
	Delete(ctx context.Context, id string, userID string) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// GetByIDAndUserID implements Repository.
func (d *dbRepository) GetByIDAndUserID(ctx context.Context, id string, userID string) (*Cat, error) {
	getUserQuery := `
		SELECT id, user_id, name, race, sex, age_in_month, description, has_matched, image_urls, created_at
		FROM cats
		WHERE id = $1 AND user_id = $2;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, id, userID)
	cat := &Cat{}
	err := row.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Race, &cat.Sex, &cat.Age, &cat.Description, &cat.HasMatched, pq.Array(&cat.ImageURLS), cat.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCatNotFound
	}
	if err != nil {
		return nil, err
	}
	return cat, nil
}

// Create implements Repository.
func (d *dbRepository) Create(ctx context.Context, cat *Cat) error {
	createCatQuery := `
		INSERT INTO cats (
			id, user_id, name, race, sex, age_in_month, description, has_matched, image_urls
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		);
	`
	_, err := d.db.DB().ExecContext(ctx, createCatQuery,
		cat.ID, cat.UserID, cat.Name, cat.Race, cat.Sex, cat.Age, cat.Description, cat.HasMatched, pq.Array(cat.ImageURLS))
	return err
}

// Update implements Repository.
func (d *dbRepository) Update(ctx context.Context, cat *Cat) error {
	updateQuery := `
		UPDATE cats
		SET name = $1,
		race = $2,
		sex = $3,
		age_in_month = $4,
		description = $5,
		has_matched = $6,
		image_urls = $7
		WHERE id = $8 AND user_id = $9
	`
	err := d.db.DB().QueryRowContext(ctx, updateQuery, cat.Name, cat.Race, cat.Sex, cat.Age, cat.Description, cat.HasMatched, pq.Array(cat.ImageURLS), cat.ID, cat.UserID).Err()
	if errors.Is(err, sql.ErrNoRows) {
		return ErrCatNotFound
	}
	return err
}

// Delete implements Repository.
func (d *dbRepository) Delete(ctx context.Context, id string, userID string) error {
	deleteCatQuery := `
		DELETE FROM cats
		WHERE id = $1 and user_id = $2;
	`
	row, err := d.db.DB().ExecContext(ctx, deleteCatQuery, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCatNotFound
	}
	return nil
}
