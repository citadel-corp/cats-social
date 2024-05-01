package cat

import (
	"context"

	"github.com/citadel-corp/cats-social/internal/common/db"
	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, cat *Cat) error
	Delete(ctx context.Context, id string, userID string) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
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
