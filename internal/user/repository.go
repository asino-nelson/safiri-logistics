package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlreadyUsed = errors.New("email already in use")
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
	)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrEmailAlreadyUsed
	}

	return err
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`

	var account User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.PasswordHash,
		&account.Role,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, mapNotFound(err)
	}

	return &account, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`

	var account User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.PasswordHash,
		&account.Role,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, mapNotFound(err)
	}

	return &account, nil
}

func mapNotFound(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}

	return err
}
