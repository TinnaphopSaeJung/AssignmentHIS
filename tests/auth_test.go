package tests

import (
	"context"
	"errors"
	"testing"

	"his/internal/dto"
	"his/internal/models"
	"his/internal/service"
	"his/pkg/utils"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type fakeStaffRepository struct {
	usernameExists bool
	findStaff      *dto.StaffWithHospital
	findErr        error
	createErr      error
	createdStaff   *models.Staff
}

func (f *fakeStaffRepository) Create(ctx context.Context, staff *models.Staff) error {
	f.createdStaff = staff
	return f.createErr
}

func (f *fakeStaffRepository) FindStaffByUsername(ctx context.Context, username string) (*dto.StaffWithHospital, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}

	return f.findStaff, nil
}

func (f *fakeStaffRepository) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	return f.usernameExists, nil
}

func TestCreateStaffSuccess(t *testing.T) {
	repo := &fakeStaffRepository{
		usernameExists: false,
	}

	jwtManager := utils.NewJWTManager("test_secret")
	authService := service.NewAuthService(repo, jwtManager)

	statusCode, err := authService.CreateStaff(context.Background(), dto.CreateStaffInput{
		Username:   "cake",
		Password:   "Password@123",
		HospitalID: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, 201, statusCode)
	assert.NotNil(t, repo.createdStaff)
	assert.Equal(t, "cake", repo.createdStaff.Username)
	assert.Equal(t, int64(1), repo.createdStaff.HospitalID)
	assert.NotEqual(t, "Password@123", repo.createdStaff.PasswordHash)
}

func TestCreateStaffInvalidPassword(t *testing.T) {
	repo := &fakeStaffRepository{}
	jwtManager := utils.NewJWTManager("test_secret")
	authService := service.NewAuthService(repo, jwtManager)

	statusCode, err := authService.CreateStaff(context.Background(), dto.CreateStaffInput{
		Username:   "cake",
		Password:   "123",
		HospitalID: 1,
	})

	assert.Error(t, err)
	assert.Equal(t, 400, statusCode)
	assert.Contains(t, err.Error(), "Password")
}

func TestCreateStaffUsernameAlreadyExists(t *testing.T) {
	repo := &fakeStaffRepository{
		usernameExists: true,
	}

	jwtManager := utils.NewJWTManager("test_secret")
	authService := service.NewAuthService(repo, jwtManager)

	statusCode, err := authService.CreateStaff(context.Background(), dto.CreateStaffInput{
		Username:   "cake",
		Password:   "Password@123",
		HospitalID: 1,
	})

	assert.Error(t, err)
	assert.Equal(t, 409, statusCode)
	assert.Contains(t, err.Error(), "username")
}

func TestCreateStaffRepositoryError(t *testing.T) {
	repo := &fakeStaffRepository{
		usernameExists: false,
		createErr:      errors.New("db error"),
	}

	jwtManager := utils.NewJWTManager("test_secret")
	authService := service.NewAuthService(repo, jwtManager)

	statusCode, err := authService.CreateStaff(context.Background(), dto.CreateStaffInput{
		Username:   "cake",
		Password:   "Password@123",
		HospitalID: 1,
	})

	assert.Error(t, err)
	assert.Equal(t, 500, statusCode)
}

func TestLoginSuccess(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), bcrypt.DefaultCost)

	repo := &fakeStaffRepository{
		findStaff: &dto.StaffWithHospital{
			ID:           1,
			Username:     "cake",
			PasswordHash: string(hash),
			HospitalID:   10,
			HospitalName: "Hospital A",
		},
	}

	jwtManager := utils.NewJWTManager("test_secret")
	authService := service.NewAuthService(repo, jwtManager)

	res, statusCode, err := authService.Login(context.Background(), "cake", "Password@123")

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.AccessToken)
	assert.Equal(t, int64(1), res.ID)
	assert.Equal(t, "cake", res.Username)
	assert.Equal(t, int64(10), res.HospitalID)
	assert.Equal(t, "Hospital A", res.HospitalName)
}

func TestLoginInvalidUsername(t *testing.T) {
	repo := &fakeStaffRepository{
		findErr: errors.New("not found"),
	}

	jwtManager := utils.NewJWTManager("test_secret")
	authService := service.NewAuthService(repo, jwtManager)

	res, statusCode, err := authService.Login(context.Background(), "wrong_user", "Password@123")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, 400, statusCode)
	assert.Contains(t, err.Error(), "Invalid username or password")
}

func TestLoginInvalidPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), bcrypt.DefaultCost)

	repo := &fakeStaffRepository{
		findStaff: &dto.StaffWithHospital{
			ID:           1,
			Username:     "cake",
			PasswordHash: string(hash),
			HospitalID:   10,
			HospitalName: "Hospital A",
		},
	}

	jwtManager := utils.NewJWTManager("test_secret")
	authService := service.NewAuthService(repo, jwtManager)

	res, statusCode, err := authService.Login(context.Background(), "cake", "WrongPassword@123")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, 400, statusCode)
	assert.Contains(t, err.Error(), "Invalid username or password")
}
