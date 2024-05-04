package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/cats-social/internal/common/db"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uint64) (*User, error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// Create implements Repository.
func (d *dbRepository) Create(ctx context.Context, user *User) (*User, error) {
	createUserQuery := `
		INSERT INTO users (
			uid, email, name, hashed_password
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id;
	`
	row := d.db.DB().QueryRowContext(ctx, createUserQuery, user.UID, user.Email, user.Name, user.HashedPassword)
	u := &User{}
	err := row.Scan(&u.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, ErrEmailAlreadyExists
			default:
				return nil, err
			}
		}
		return nil, err
	}
	return u, nil
}

func (d *dbRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	getUserQuery := `
		SELECT id, uid, email, name, hashed_password FROM users
		WHERE email = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, email)
	u := &User{}
	err := row.Scan(&u.ID, &u.UID, &u.Email, &u.Name, &u.HashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (d *dbRepository) GetByID(ctx context.Context, id uint64) (*User, error) {
	getUserQuery := `
		SELECT id, uid, email, name, product_sold_total, hashed_password FROM users
		WHERE uid = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, id)
	u := &User{}
	err := row.Scan(&u.ID, &u.UID, &u.Email, &u.Name, &u.HashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}
