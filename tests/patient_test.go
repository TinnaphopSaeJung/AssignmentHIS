package tests

import (
	"context"
	"errors"
	"testing"

	"his/internal/dto"
	"his/internal/service"

	"github.com/stretchr/testify/assert"
)

type fakePatientRepository struct {
	items []dto.PatientResponse
	total int
	err   error
}

func (f *fakePatientRepository) Search(
	ctx context.Context,
	hospitalID int64,
	req dto.SearchPatientRequest,
) ([]dto.PatientResponse, int, error) {
	return f.items, f.total, f.err
}

type fakeHospitalAClient struct {
	patient    *dto.HospitalAPatientResponse
	statusCode int
	err        error
}

func (f *fakeHospitalAClient) SearchPatient(
	ctx context.Context,
	id string,
) (*dto.HospitalAPatientResponse, int, error) {
	return f.patient, f.statusCode, f.err
}

func TestSearchPatientSuccess(t *testing.T) {
	repo := &fakePatientRepository{
		items: []dto.PatientResponse{
			{
				ID:          1,
				FirstNameTH: "สมชาย",
				LastNameTH:  "สายลม",
				PatientHN:   "HNA001",
			},
		},
		total: 1,
	}

	client := &fakeHospitalAClient{}
	patientService := service.NewPatientService(repo, client)

	res, statusCode, err := patientService.Search(context.Background(), 1, dto.SearchPatientRequest{
		FirstName: "สมชาย",
		Page:      1,
		Limit:     10,
	})

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.NotNil(t, res)
	assert.Len(t, res.Items, 1)
	assert.Equal(t, 1, res.Pagination.Page)
	assert.Equal(t, 10, res.Pagination.Limit)
	assert.Equal(t, 1, res.Pagination.Total)
	assert.Equal(t, 1, res.Pagination.LastPage)
	assert.Nil(t, res.Pagination.PreviousPage)
	assert.Nil(t, res.Pagination.NextPage)
}

func TestSearchPatientDefaultPagination(t *testing.T) {
	repo := &fakePatientRepository{
		items: []dto.PatientResponse{},
		total: 0,
	}

	client := &fakeHospitalAClient{}
	patientService := service.NewPatientService(repo, client)

	res, statusCode, err := patientService.Search(context.Background(), 1, dto.SearchPatientRequest{})

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.NotNil(t, res)
	assert.Equal(t, 1, res.Pagination.Page)
	assert.Equal(t, 10, res.Pagination.Limit)
	assert.Equal(t, 0, res.Pagination.Total)
	assert.Equal(t, 0, res.Pagination.LastPage)
}

func TestSearchPatientInvalidDateFormat(t *testing.T) {
	repo := &fakePatientRepository{}
	client := &fakeHospitalAClient{}
	patientService := service.NewPatientService(repo, client)

	res, statusCode, err := patientService.Search(context.Background(), 1, dto.SearchPatientRequest{
		DateOfBirth: "01-01-1995",
		Page:        1,
		Limit:       10,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, 409, statusCode)
	assert.Contains(t, err.Error(), "date_of_birth")
}

func TestSearchPatientRepositoryError(t *testing.T) {
	repo := &fakePatientRepository{
		err: errors.New("db error"),
	}

	client := &fakeHospitalAClient{}
	patientService := service.NewPatientService(repo, client)

	res, statusCode, err := patientService.Search(context.Background(), 1, dto.SearchPatientRequest{
		Page:  1,
		Limit: 10,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, 500, statusCode)
}

func TestSearchPatientPaginationMultiplePages(t *testing.T) {
	repo := &fakePatientRepository{
		items: []dto.PatientResponse{
			{ID: 11, FirstNameEN: "John"},
			{ID: 12, FirstNameEN: "Jane"},
		},
		total: 25,
	}

	client := &fakeHospitalAClient{}
	patientService := service.NewPatientService(repo, client)

	res, statusCode, err := patientService.Search(context.Background(), 1, dto.SearchPatientRequest{
		Page:  2,
		Limit: 10,
	})

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.NotNil(t, res)
	assert.Equal(t, 2, res.Pagination.Page)
	assert.Equal(t, 10, res.Pagination.Limit)
	assert.Equal(t, 25, res.Pagination.Total)
	assert.Equal(t, 3, res.Pagination.LastPage)

	if assert.NotNil(t, res.Pagination.PreviousPage) {
		assert.Equal(t, 1, *res.Pagination.PreviousPage)
	}

	if assert.NotNil(t, res.Pagination.NextPage) {
		assert.Equal(t, 3, *res.Pagination.NextPage)
	}
}

func TestSearchFromHISExternalSuccess(t *testing.T) {
	repo := &fakePatientRepository{}
	client := &fakeHospitalAClient{
		patient: &dto.HospitalAPatientResponse{
			FirstNameTH: "สมชาย",
			LastNameTH:  "สายลม",
			NationalID:  "1100100123456",
			PatientHN:   "HNA001",
		},
		statusCode: 200,
	}

	patientService := service.NewPatientService(repo, client)

	res, statusCode, err := patientService.SearchFromHISExternal(context.Background(), "1100100123456")

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.NotNil(t, res)
	assert.Equal(t, "สมชาย", res.FirstNameTH)
	assert.Equal(t, "1100100123456", res.NationalID)
}

func TestSearchFromHISExternalNotFound(t *testing.T) {
	repo := &fakePatientRepository{}
	client := &fakeHospitalAClient{
		patient:    nil,
		statusCode: 404,
		err:        errors.New("patient not found from hospital A"),
	}

	patientService := service.NewPatientService(repo, client)

	res, statusCode, err := patientService.SearchFromHISExternal(context.Background(), "not-found-id")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, 404, statusCode)
	assert.Contains(t, err.Error(), "not found")
}
