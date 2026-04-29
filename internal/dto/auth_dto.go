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
