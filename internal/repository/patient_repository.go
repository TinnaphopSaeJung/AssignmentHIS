package repository

import (
	"context"
	"fmt"

	"his/internal/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PatientRepository struct {
	db *pgxpool.Pool
}

func NewPatientRepository(db *pgxpool.Pool) *PatientRepository {
	return &PatientRepository{
		db: db,
	}
}

func (r *PatientRepository) Search(ctx context.Context, hospitalID int64, req dto.SearchPatientRequest) ([]dto.PatientResponse, int, error) {
	baseQuery := `
		FROM patients
		JOIN patient_hospitals_mapping 
			ON patients.id = patient_hospitals_mapping.patient_id
		WHERE patient_hospitals_mapping.hospital_id = $1
			AND patients.deleted_at IS NULL
			AND patient_hospitals_mapping.deleted_at IS NULL
	`

	args := []interface{}{hospitalID}
	idx := 2

	if req.NationalID != "" {
		baseQuery += fmt.Sprintf(" AND patients.national_id = $%d", idx)
		args = append(args, req.NationalID)
		idx++
	}

	if req.PassportID != "" {
		baseQuery += fmt.Sprintf(" AND patients.passport_id = $%d", idx)
		args = append(args, req.PassportID)
		idx++
	}

	if req.FirstName != "" {
		baseQuery += fmt.Sprintf(`
			AND (
				patients.first_name_th ILIKE $%d
				OR patients.first_name_en ILIKE $%d
			)
		`, idx, idx)
		args = append(args, "%"+req.FirstName+"%")
		idx++
	}

	if req.MiddleName != "" {
		baseQuery += fmt.Sprintf(`
			AND (
				patients.middle_name_th ILIKE $%d
				OR patients.middle_name_en ILIKE $%d
			)
		`, idx, idx)
		args = append(args, "%"+req.MiddleName+"%")
		idx++
	}

	if req.LastName != "" {
		baseQuery += fmt.Sprintf(`
			AND (
				patients.last_name_th ILIKE $%d
				OR patients.last_name_en ILIKE $%d
			)
		`, idx, idx)
		args = append(args, "%"+req.LastName+"%")
		idx++
	}

	if req.DateOfBirth != "" {
		baseQuery += fmt.Sprintf(" AND patients.date_of_birth = $%d", idx)
		args = append(args, req.DateOfBirth)
		idx++
	}

	if req.PhoneNumber != "" {
		baseQuery += fmt.Sprintf(" AND patients.phone_number ILIKE $%d", idx)
		args = append(args, "%"+req.PhoneNumber+"%")
		idx++
	}

	if req.Email != "" {
		baseQuery += fmt.Sprintf(" AND patients.email ILIKE $%d", idx)
		args = append(args, "%"+req.Email+"%")
		idx++
	}

	countQuery := "SELECT COUNT(*) " + baseQuery

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT 
			patients.id,
			COALESCE(patients.first_name_th, ''),
			COALESCE(patients.middle_name_th, ''),
			COALESCE(patients.last_name_th, ''),
			COALESCE(patients.first_name_en, ''),
			COALESCE(patients.middle_name_en, ''),
			COALESCE(patients.last_name_en, ''),
			COALESCE(patients.date_of_birth::text, ''),
			COALESCE(patients.national_id, ''),
			COALESCE(patients.passport_id, ''),
			COALESCE(patients.phone_number, ''),
			COALESCE(patients.email, ''),
			COALESCE(patients.gender, ''),
			COALESCE(patient_hospitals_mapping.patient_hn, '')
	` + baseQuery

	offset := (req.Page - 1) * req.Limit

	dataQuery += fmt.Sprintf(" ORDER BY patients.id DESC LIMIT $%d OFFSET $%d", idx, idx+1)

	args = append(args, req.Limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []dto.PatientResponse

	for rows.Next() {
		var p dto.PatientResponse

		err := rows.Scan(
			&p.ID,
			&p.FirstNameTH,
			&p.MiddleNameTH,
			&p.LastNameTH,
			&p.FirstNameEN,
			&p.MiddleNameEN,
			&p.LastNameEN,
			&p.DateOfBirth,
			&p.NationalID,
			&p.PassportID,
			&p.PhoneNumber,
			&p.Email,
			&p.Gender,
			&p.PatientHN,
		)
		if err != nil {
			return nil, 0, err
		}

		results = append(results, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
