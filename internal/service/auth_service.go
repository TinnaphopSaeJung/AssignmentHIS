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
	repo       *repository.StaffRepository
	jwtManager *utils.JWTManager
}

func NewAuthService(repo *repository.StaffRepository, jwtManager *utils.JWTManager) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) CreateStaff(ctx context.Context, input dto.CreateStaffInput) error {
	if !utils.IsValidPassword(input.Password) {
		return errors.New("Password must be at least 8 characters and include letters, numbers, and special characters.")
	}

	existing, _ := s.repo.FindStaffByUsername(ctx, input.Username)
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

func (s *AuthService) Login(ctx context.Context, username, password string) (*dto.LoginResponse, error) {
	staff, err := s.repo.FindStaffByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("Invalid username or password.")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(staff.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, errors.New("Invalid username or password.")
	}

	token, err := s.jwtManager.GenerateJWT(staff.ID, staff.HospitalID)
	if err != nil {
		return nil, errors.New("Cannot generate token.")
	}

	return &dto.LoginResponse{
		AccessToken:  token,
		ID:           staff.ID,
		Username:     staff.Username,
		HospitalID:   staff.HospitalID,
		HospitalName: staff.HospitalName,
	}, nil
}
