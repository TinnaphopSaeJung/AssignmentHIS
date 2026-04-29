package service

import (
	"context"
	"errors"

	"his/internal/dto"
	"his/internal/repository"
	"his/pkg/utils"
)

type PatientService struct {
	repo *repository.PatientRepository
}

func NewPatientService(repo *repository.PatientRepository) *PatientService {
	return &PatientService{
		repo: repo,
	}
}

func (s *PatientService) Search(ctx context.Context, hospitalID int64, req dto.SearchPatientRequest) (*dto.SearchPatientResponse, error) {
	if !utils.IsValidDate(req.DateOfBirth) {
		return nil, errors.New("date_of_birth must be in YYYY-MM-DD format.")
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	items, total, err := s.repo.Search(ctx, hospitalID, req)
	if err != nil {
		return nil, err
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
	}, nil
}
