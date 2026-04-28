package repository

import (
	"context"
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
			username, password_hash, hospital_id,
			first_name_th, middle_name_th, last_name_th,
			first_name_en, middle_name_en, last_name_en
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := r.db.Exec(ctx, query,
		staff.Username,
		staff.PasswordHash,
		staff.HospitalID,
		staff.FirstNameTH,
		staff.MiddleNameTH,
		staff.LastNameTH,
		staff.FirstNameEN,
		staff.MiddleNameEN,
		staff.LastNameEN,
	)

	return err
}

func (r *StaffRepository) FindByUsername(ctx context.Context, username string) (*models.Staff, error) {
	query := `
		SELECT id, username
		FROM staffs
		WHERE username = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRow(ctx, query, username)

	var staff models.Staff
	err := row.Scan(&staff.ID, &staff.Username)
	if err != nil {
		return nil, err
	}

	return &staff, nil
}
