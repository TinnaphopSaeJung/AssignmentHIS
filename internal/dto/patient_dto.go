package dto

type SearchPatientRequest struct {
	NationalID  string `json:"national_id"`
	PassportID  string `json:"passport_id"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`

	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PatientResponse struct {
	ID           int64  `json:"id"`
	FirstNameTH  string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH   string `json:"last_name_th"`
	FirstNameEN  string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN   string `json:"last_name_en"`
	DateOfBirth  string `json:"date_of_birth"`
	NationalID   string `json:"national_id"`
	PassportID   string `json:"passport_id"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
	PatientHN    string `json:"patient_hn"`
}

type Pagination struct {
	Page         int  `json:"current_page"`
	Limit        int  `json:"limit"`
	Total        int  `json:"total"`
	LastPage     int  `json:"last_page"`
	PreviousPage *int `json:"previous_page"`
	NextPage     *int `json:"next_page"`
}

type SearchPatientResponse struct {
	Items      []PatientResponse `json:"items"`
	Pagination Pagination        `json:"pagination"`
}
