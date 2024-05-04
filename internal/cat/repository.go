package cat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/citadel-corp/cats-social/internal/common/db"
	"github.com/lib/pq"
)

type Repository interface {
	GetByIDAndUserID(ctx context.Context, id string, userID string) (*Cat, error)
	List(ctx context.Context, req ListCatPayload, userID string) ([]*Cat, error)
	Create(ctx context.Context, cat *Cat) (*Cat, error)
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
	err := row.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Race, &cat.Sex, &cat.Age, &cat.Description, &cat.HasMatched, pq.Array(&cat.ImageURLS), &cat.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCatNotFound
	}
	if err != nil {
		return nil, err
	}
	return cat, nil
}

// List implements Repository.
func (d *dbRepository) List(ctx context.Context, req ListCatPayload, userID string) ([]*Cat, error) {
	paramNo := 1
	listQuery := "SELECT * FROM cats WHERE "
	params := make([]interface{}, 0)
	if req.Race != "" {
		listQuery += fmt.Sprintf("race = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, req.Race)
	}
	if req.Sex != "" {
		listQuery += fmt.Sprintf("sex = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, req.Sex)
	}
	switch req.HasMatchedType {
	case HasMatched:
		listQuery += fmt.Sprintf("has_matched = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, true)
	case HasNotMatched:
		listQuery += fmt.Sprintf("has_matched = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, false)
	}
	switch req.AgeSearchType {
	case MoreThan:
		listQuery += fmt.Sprintf("age_in_month > $%d AND ", paramNo)
		paramNo += 1
		params = append(params, req.Age)
	case LessThan:
		listQuery += fmt.Sprintf("age_in_month < $%d AND ", paramNo)
		paramNo += 1
		params = append(params, req.Age)
	case EqualTo:
		listQuery += fmt.Sprintf("age_in_month = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, req.Age)
	}

	if req.Owned {
		listQuery += fmt.Sprintf("user_id = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, userID)
	}
	if req.Search != "" {
		listQuery += fmt.Sprintf("name LIKE %%$%d%% AND ", paramNo)
		params = append(params, req.Search)
	}
	if strings.HasSuffix(listQuery, "AND ") {
		listQuery, _ = strings.CutSuffix(listQuery, "AND ")
	}
	listQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d;", req.Limit, req.Offset)
	if strings.Contains(listQuery, "WHERE  ORDER") {
		listQuery = strings.Replace(listQuery, "WHERE  ORDER", "ORDER", 1)
	}
	rows, err := d.db.DB().QueryContext(ctx, listQuery, params...)
	if err != nil {
		return nil, err
	}
	res := make([]*Cat, 0)
	for rows.Next() {
		cat := &Cat{}
		err = rows.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Race, &cat.Sex, &cat.Age, &cat.Description, &cat.HasMatched, pq.Array(&cat.ImageURLS), &cat.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, cat)
	}
	return res, nil
}

// Create implements Repository.
func (d *dbRepository) Create(ctx context.Context, cat *Cat) (*Cat, error) {
	createCatQuery := `
		INSERT INTO cats (
			id, user_id, name, race, sex, age_in_month, description, has_matched, image_urls
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id, created_at;
	`
	row := d.db.DB().QueryRowContext(ctx, createCatQuery,
		cat.ID, cat.UserID, cat.Name, cat.Race, cat.Sex, cat.Age, cat.Description, cat.HasMatched, pq.Array(cat.ImageURLS))
	c := &Cat{}
	err := row.Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return c, nil
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
