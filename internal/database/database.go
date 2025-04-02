package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-users/internal/api"
	"go-users/internal/config"
	"go-users/internal/ownErrors"
)

// ConnPool represents a connection pool abstraction for executing queries and managing database connections.
type ConnPool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}

// db represents a database connection abstraction with a connection pool for executing queries and managing transactions.
type db struct {
	pool ConnPool
}

// DB defines an interface for interacting with the database, including user management and resource cleanup.
type DB interface {
	CreateUser(ctx context.Context, user *api.UserRequest) (*api.User, error)
	GetUser(ctx context.Context, id uint) (*api.User, error)
	UpdateUser(ctx context.Context, u *api.UserRequest, id uint) (*api.User, error)
	Close()
}

// New initializes a new database connection pool using the provided context and configuration. Returns a DB instance or an error.
func New(ctx context.Context, cfg config.Database) (DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &db{pool: pool}, nil
}

// Close releases all resources associated with the database connection pool if it is initialized.
func (db *db) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

// CreateUser creates a new user in the database
func (db *db) CreateUser(ctx context.Context, u *api.UserRequest) (*api.User, error) {
	query := "INSERT INTO users (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING id, first_name, last_name, email, created_at, updated_at"

	var user api.User
	err := db.pool.QueryRow(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case isDuplicateKeyError(err):
			return nil, ownErrors.ErrUserAlreadyExists
		default:
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	return &user, nil
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	if !ok {
		return false
	}
	return pgErr.Code == "23505"
}

// GetUser retrieves a user from the database by their unique ID. Returns the user or an error if not found or on failure.
func (db *db) GetUser(ctx context.Context, id uint) (*api.User, error) {
	query := "SELECT id, first_name, last_name, email, created_at, updated_at FROM users WHERE id = $1"

	var user api.User
	err := db.pool.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ownErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser updates an existing user's details in the database and sets the updated timestamp in the User struct.
// Returns ErrNotFound if the user does not exist or an error if the operation fails.
func (db *db) UpdateUser(ctx context.Context, u *api.UserRequest, id uint) (*api.User, error) {
	query := "UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4 RETURNING id, first_name, last_name, email, created_at, updated_at"

	var user api.User
	err := db.pool.QueryRow(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		id,
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ownErrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}
