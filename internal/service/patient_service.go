package service

import (
	"context"
	"errors"
	"fmt"

	"his/internal/dto"
	"his/pkg/utils"
)

type PatientRepository interface {
	Search(ctx context.Context, hospitalID int64, req dto.SearchPatientRequest) ([]dto.PatientResponse, int, error)
}

type HospitalAClient interface {
	SearchPatient(ctx context.Context, id string) (*dto.HospitalAPatientResponse, int, error)
}

type PatientService struct {
	repo            PatientRepository
	hospitalAClient HospitalAClient
}

func NewPatientService(repo PatientRepository, hospitalAClient HospitalAClient) *PatientService {
	return &PatientService{
		repo:            repo,
		hospitalAClient: hospitalAClient,
	}
}

func (s *PatientService) Search(ctx context.Context, hospitalID int64, req dto.SearchPatientRequest) (*dto.SearchPatientResponse, int, error) {
	if !utils.IsValidDate(req.DateOfBirth) {
		return nil, 409, errors.New("date_of_birth must be in YYYY-MM-DD format.")
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	items, total, err := s.repo.Search(ctx, hospitalID, req)
	if err != nil {
		return nil, 500, fmt.Errorf("Internal Server Error: %w", err)
	}

	lastPage := (total + req.Limit - 1) / req.Limit

	var prevPage *int
	var nextPage *int

	if req.Page > 1 {
		p := req.Page - 1
		prevPage = &p
	}

	if req.Page < lastPage {
		n := req.Page + 1
		nextPage = &n
	}

	return &dto.SearchPatientResponse{
		Items: items,
		Pagination: dto.Pagination{
			Page:         req.Page,
			Limit:        req.Limit,
			Total:        total,
			LastPage:     lastPage,
			PreviousPage: prevPage,
			NextPage:     nextPage,
		},
	}, 200, nil
}

func (s *PatientService) SearchFromHISExternal(ctx context.Context, id string) (*dto.HospitalAPatientResponse, int, error) {
	patient, statusCode, err := s.hospitalAClient.SearchPatient(ctx, id)
	if err != nil {
		return nil, statusCode, err
	}

	return patient, 200, nil
}
