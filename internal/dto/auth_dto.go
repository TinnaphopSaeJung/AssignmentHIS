package dto

type CreateStaffRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	HospitalID int64  `json:"hospital_id" binding:"required"`
}

type CreateStaffInput struct {
	Username   string
	Password   string
	HospitalID int64
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	HospitalID   int64  `json:"hospital_id"`
	HospitalName string `json:"hospital_name"`
}

type StaffWithHospital struct {
	ID           int64
	Username     string
	PasswordHash string
	HospitalID   int64
	HospitalName string
}
