package service

import (
	"context"
	"errors"
	"fmt"
	"his/internal/dto"
	"his/internal/models"
	"his/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type StaffRepository interface {
	Create(ctx context.Context, staff *models.Staff) error
	FindStaffByUsername(ctx context.Context, username string) (*dto.StaffWithHospital, error)
	IsUsernameExists(ctx context.Context, username string) (bool, error)
}

type AuthService struct {
	repo       StaffRepository
	jwtManager *utils.JWTManager
}

func NewAuthService(repo StaffRepository, jwtManager *utils.JWTManager) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) CreateStaff(ctx context.Context, input dto.CreateStaffInput) (int, error) {
	if !utils.IsValidPassword(input.Password) {
		return 400, errors.New("Password must be at least 8 characters and include letters, numbers, and special characters.")
	}

	exists, err := s.repo.IsUsernameExists(ctx, input.Username)
	if err != nil {
		return 500, fmt.Errorf("Internal Server Error: %w", err)
	}

	if exists {
		return 409, errors.New("This username already exists.")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 500, errors.New("Internal Server Error.")
	}

	staff := &models.Staff{
		Username:     input.Username,
		PasswordHash: string(hash),
		HospitalID:   input.HospitalID,
	}

	if err := s.repo.Create(ctx, staff); err != nil {
		return 500, fmt.Errorf("Internal Server Error: %w", err)
	}

	return 201, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*dto.LoginResponse, int, error) {
	staff, err := s.repo.FindStaffByUsername(ctx, username)
	if err != nil {
		return nil, 400, errors.New("Invalid username or password.")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(staff.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, 400, errors.New("Invalid username or password.")
	}

	token, err := s.jwtManager.GenerateJWT(staff.ID, staff.HospitalID)
	if err != nil {
		return nil, 500, errors.New("Cannot generate token.")
	}

	return &dto.LoginResponse{
		AccessToken:  token,
		ID:           staff.ID,
		Username:     staff.Username,
		HospitalID:   staff.HospitalID,
		HospitalName: staff.HospitalName,
	}, 200, nil
}
