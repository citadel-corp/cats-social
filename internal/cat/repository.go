package cat

import (
	"context"

	"github.com/citadel-corp/cats-social/internal/common/db"
	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, cat *Cat) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// Create implements Repository.
func (d *dbRepository) Create(ctx context.Context, cat *Cat) error {
	createUserQuery := `
		INSERT INTO cats (
			id, user_id, name, race, sex, age_in_month, description, has_matched, image_urls
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		);
	`
	_, err := d.db.DB().ExecContext(ctx, createUserQuery,
		cat.ID, cat.UserID, cat.Name, cat.Race, cat.Sex, cat.Age, cat.Description, cat.HasMatched, pq.Array(cat.ImageURLS))
	return err
}
