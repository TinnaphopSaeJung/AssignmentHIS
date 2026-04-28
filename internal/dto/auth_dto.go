package dto

// ใช้รับ request จาก API
type CreateStaffRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	HospitalID int64  `json:"hospital_id" binding:"required"`

	FirstNameTH  string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH   string `json:"last_name_th"`

	FirstNameEN  string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN   string `json:"last_name_en"`
}

// ใช้ใน service
type CreateStaffInput struct {
	Username   string
	Password   string
	HospitalID int64

	FirstNameTH  string
	MiddleNameTH string
	LastNameTH   string

	FirstNameEN  string
	MiddleNameEN string
	LastNameEN   string
}
