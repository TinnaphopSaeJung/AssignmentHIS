package repository

import (
	"context"
	"encoding/json"

	"his/internal/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditLogRepository struct {
	db *pgxpool.Pool
}

func NewAuditLogRepository(db *pgxpool.Pool) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, event dto.AuditLogEvent) error {
	query := `
		INSERT INTO audit_logs (
			event_type,
			staff_id,
			hospital_id,
			username,
			ip_address,
			metadata,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	metadataBytes, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		ctx,
		query,
		event.EventType,
		event.StaffID,
		event.HospitalID,
		event.Username,
		event.IPAddress,
		metadataBytes,
		event.CreatedAt,
	)

	return err
}
