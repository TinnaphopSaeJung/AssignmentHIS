package service

import (
	"context"
	"errors"
	"fmt"
	"his/internal/dto"
	"his/internal/models"
	"his/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type StaffRepository interface {
	Create(ctx context.Context, staff *models.Staff) error
	FindStaffByUsername(ctx context.Context, username string) (*dto.StaffWithHospital, error)
	IsUsernameExists(ctx context.Context, username string) (bool, error)
	FindStaffByID(ctx context.Context, staffID int64) (*dto.StaffWithHospital, error)
}

type LoginAttempt interface {
	IsLocked(ctx context.Context, username string) (bool, time.Duration, error)
	RecordFailedAttempt(ctx context.Context, username string) (int64, error)
	Reset(ctx context.Context, username string) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, staffID int64, tokenHash string, expiresAt time.Time) error
	FindValidToken(ctx context.Context, staffID int64, tokenHash string) (bool, error)
	Revoke(ctx context.Context, staffID int64, tokenHash string) error
}

type AuthService struct {
	repo                StaffRepository
	refreshTokenRepo    RefreshTokenRepository
	jwtManager          *utils.JWTManager
	loginAttemptService LoginAttempt
}

func NewAuthService(repo StaffRepository, refreshTokenRepo RefreshTokenRepository, jwtManager *utils.JWTManager, loginAttemptService LoginAttempt) *AuthService {
	return &AuthService{
		repo:                repo,
		refreshTokenRepo:    refreshTokenRepo,
		jwtManager:          jwtManager,
		loginAttemptService: loginAttemptService,
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
		return 500, fmt.Errorf("Internal Server Error: %w", err)
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
	isLocked, ttl, err := s.loginAttemptService.IsLocked(ctx, username)
	if err != nil {
		return nil, 500, errors.New("Cannot check login attempt.")
	}

	if isLocked {
		return nil, 429, fmt.Errorf("Too many failed login attempts. Please try again in %d minute(s).", int(ttl.Minutes())+1)
	}

	staff, err := s.repo.FindStaffByUsername(ctx, username)
	if err != nil {
		statusCode, err := s.handleFailedLogin(ctx, username)
		return nil, statusCode, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(password))
	if err != nil {
		statusCode, err := s.handleFailedLogin(ctx, username)
		return nil, statusCode, err
	}

	if err := s.loginAttemptService.Reset(ctx, username); err != nil {
		return nil, 500, errors.New("Cannot reset login attempt.")
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(staff.ID, staff.HospitalID)
	if err != nil {
		return nil, 500, errors.New("Cannot generate access token.")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(staff.ID)
	if err != nil {
		return nil, 500, errors.New("Cannot generate refresh token.")
	}

	refreshTokenHash := utils.HashToken(refreshToken)
	refreshTokenExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	if err := s.refreshTokenRepo.Create(ctx, staff.ID, refreshTokenHash, refreshTokenExpiresAt); err != nil {
		return nil, 500, errors.New("Cannot save refresh token.")
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ID:           staff.ID,
		Username:     staff.Username,
		HospitalID:   staff.HospitalID,
		HospitalName: staff.HospitalName,
	}, 200, nil
}

func (s *AuthService) handleFailedLogin(ctx context.Context, username string) (int, error) {
	count, err := s.loginAttemptService.RecordFailedAttempt(ctx, username)
	if err != nil {
		return 500, errors.New("Cannot record login attempt.")
	}

	if count >= maxLoginFailedAttempts {
		return 429, errors.New("Too many failed login attempts. Please try again in 15 minute(s).")
	}

	return 400, errors.New("Invalid username or password.")
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, int, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, 401, errors.New("Invalid refresh token.")
	}

	tokenHash := utils.HashToken(refreshToken)

	exists, err := s.refreshTokenRepo.FindValidToken(ctx, claims.StaffID, tokenHash)
	if err != nil {
		return nil, 500, errors.New("Cannot validate refresh token.")
	}

	if !exists {
		return nil, 401, errors.New("Invalid refresh token.")
	}

	staff, err := s.repo.FindStaffByID(ctx, claims.StaffID)
	if err != nil {
		return nil, 401, errors.New("Invalid refresh token.")
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(staff.ID, staff.HospitalID)
	if err != nil {
		return nil, 500, errors.New("Cannot generate access token.")
	}

	return &dto.RefreshTokenResponse{
		AccessToken: accessToken,
	}, 200, nil
}
