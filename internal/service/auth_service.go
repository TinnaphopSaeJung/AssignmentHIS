package service

import (
	"context"
	"errors"
	"his/internal/dto"
	"his/internal/models"
	"his/internal/repository"
	"his/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.StaffRepository
}

func NewAuthService(repo *repository.StaffRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateStaff(ctx context.Context, input dto.CreateStaffInput) error {
	if !utils.IsValidPassword(input.Password) {
		return errors.New("Password must be at least 8 characters and include letters, numbers, and special characters.")
	}

	existing, _ := s.repo.FindByUsername(ctx, input.Username)
	if existing != nil {
		return errors.New("This username already exists.")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	staff := &models.Staff{
		Username:     input.Username,
		PasswordHash: string(hash),
		HospitalID:   input.HospitalID,
	}

	return s.repo.Create(ctx, staff)
}
