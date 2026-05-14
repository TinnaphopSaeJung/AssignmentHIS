package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, staffID int64, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (
			staff_id, token_hash, expires_at
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query, staffID, tokenHash, expiresAt)

	return err
}

func (r *RefreshTokenRepository) FindValidToken(ctx context.Context, staffID int64, tokenHash string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM refresh_tokens
			WHERE staff_id = $1
			  AND token_hash = $2
			  AND revoked_at IS NULL
			  AND expires_at > NOW()
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, staffID, tokenHash).Scan(&exists)

	return exists, err
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, staffID int64, tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE staff_id = $1
		  AND token_hash = $2
		  AND revoked_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, staffID, tokenHash)

	return err
}
