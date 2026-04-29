package repository

import (
	"context"
	"his/internal/dto"
	"his/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StaffRepository struct {
	db *pgxpool.Pool
}

func NewStaffRepository(db *pgxpool.Pool) *StaffRepository {
	return &StaffRepository{db: db}
}

func (r *StaffRepository) Create(ctx context.Context, staff *models.Staff) error {
	query := `
		INSERT INTO staffs (
			username, password_hash, hospital_id
		)
		VALUES ($1,$2,$3)
	`

	_, err := r.db.Exec(ctx, query,
		staff.Username,
		staff.PasswordHash,
		staff.HospitalID,
	)

	return err
}

func (r *StaffRepository) FindStaffByUsername(ctx context.Context, username string) (*dto.StaffWithHospital, error) {
	query := `
		SELECT 
			staffs.id,
			staffs.username,
			staffs.password_hash,
			staffs.hospital_id,
			hospitals.name
		FROM staffs
		JOIN hospitals ON staffs.hospital_id = hospitals.id
		WHERE staffs.username = $1
		AND staffs.deleted_at IS NULL
	`

	row := r.db.QueryRow(ctx, query, username)

	var staff dto.StaffWithHospital
	err := row.Scan(
		&staff.ID,
		&staff.Username,
		&staff.PasswordHash,
		&staff.HospitalID,
		&staff.HospitalName,
	)
	if err != nil {
		return nil, err
	}

	return &staff, nil
}

func (r *StaffRepository) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM staffs
			WHERE username = $1
			AND deleted_at IS NULL
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
